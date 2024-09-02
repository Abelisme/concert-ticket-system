package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "concert-ticket-system/ticketservice"
)

func simulateTicketPurchase(client pb.TicketServiceClient, userID int, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	availResp, err := client.CheckAvailability(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("Error checking availability: %v", err)
		return
	}

	if availResp.Available {
		reserveResp, err := client.ReserveTicket(ctx, &pb.ReserveRequest{UserId: int32(userID)})
		if err != nil {
			log.Printf("Error reserving ticket: %v", err)
			return
		}

		if reserveResp.Success {
			fmt.Printf("User %d successfully purchased a ticket.\n", userID)
		} else {
			fmt.Printf("Sorry, tickets sold out for user %d.\n", userID)
		}
	} else {
		fmt.Println("No tickets available.")
	}
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTicketServiceClient(conn)

	var wg sync.WaitGroup
	startTime := time.Now()
	purchaseCount := 0

	for time.Since(startTime) < time.Second*5 { // 執行5 seconds
		wg.Add(1)
		go simulateTicketPurchase(client, rand.Intn(10000), &wg)
		purchaseCount++
		time.Sleep(time.Millisecond * 10) // Control request rate
	}

	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	remainingResp, err := client.GetRemainingTickets(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Could not get remaining tickets: %v", err)
	}

	fmt.Printf("\nSimulation completed. Total purchase attempts: %d\n", purchaseCount)
	fmt.Printf("Remaining tickets: %d\n", remainingResp.Count)
	fmt.Printf("Tickets sold: %d\n", 1000-int(remainingResp.Count))
}