package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	pb "concert-ticket-system/ticketservice"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

type TicketSystem struct {
	mu           sync.Mutex //使用 sync.Mutex 來保護共享資源（票務系統）
	totalTickets int
	soldTickets  map[int]bool
}

func NewTicketSystem(total int) *TicketSystem {
	return &TicketSystem{
		totalTickets: total,
		soldTickets:  make(map[int]bool),
	}
}

func (ts *TicketSystem) checkAvailability() bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.totalTickets > 0
}

func (ts *TicketSystem) reserveTicket(userID int) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.totalTickets > 0 {
		ts.totalTickets--
		ts.soldTickets[userID] = true
		return true
	}
	return false
}

type server struct {
	pb.UnimplementedTicketServiceServer
	ts *TicketSystem
}

func (s *server) CheckAvailability(ctx context.Context, in *pb.Empty) (*pb.AvailabilityResponse, error) {
	val, err := rdb.Get(ctx, "ticket_available").Result()
	if err == redis.Nil {
		available := s.ts.checkAvailability()
		rdb.Set(ctx, "ticket_available", available, time.Second*5)
		return &pb.AvailabilityResponse{Available: available}, nil
	} else if err != nil {
		log.Printf("Redis error: %v", err)
		return &pb.AvailabilityResponse{Available: s.ts.checkAvailability()}, nil
	}
	return &pb.AvailabilityResponse{Available: val == "1"}, nil
}

func (s *server) ReserveTicket(ctx context.Context, in *pb.ReserveRequest) (*pb.ReserveResponse, error) {
	if s.ts.reserveTicket(int(in.UserId)) {
		rdb.Decr(ctx, "ticket_count")
		return &pb.ReserveResponse{Success: true}, nil
	}
	return &pb.ReserveResponse{Success: false}, nil
}

func (s *server) GetRemainingTickets(ctx context.Context, in *pb.Empty) (*pb.RemainingTicketsResponse, error) {
	val, err := rdb.Get(ctx, "ticket_count").Int()
	if err != nil {
		log.Printf("Redis error: %v", err)
		return &pb.RemainingTicketsResponse{Count: int32(s.ts.totalTickets)}, nil
	}
	return &pb.RemainingTicketsResponse{Count: int32(val)}, nil
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	totalTickets := 1000
	ts := NewTicketSystem(totalTickets)

	rdb.Set(ctx, "ticket_count", totalTickets, 0)
	rdb.Set(ctx, "ticket_available", "1", time.Second*5)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTicketServiceServer(s, &server{ts: ts})
	fmt.Println("Server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}