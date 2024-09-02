# gRPC Redis 高頻搶票系統

這個項目旨在模擬高頻搶票情境，使用 gRPC 搭配 Redis 來實現。通過使用快取，我們可以有效降低服務器負載。

## 前置要求

在開始之前，請確保您的系統已安裝以下工具：

- Go 編程語言
- Protobuf 編譯器 (protoc)
- Redis 服務器

## 安裝步驟

### 1. 安裝 protoc

對於 macOS（使用 Homebrew）：

```bash
brew install protobuf
```

對於其他操作系統，請參考 [Protobuf 官方文檔](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation)。

### 2. 安裝 Go 的 protoc 插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

### 3. 更新 PATH

將以下行添加到您的 shell 配置文件（~/.bashrc 或 ~/.zshrc）中：

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 4. 重新加載 shell 配置

對於 bash 用戶：

```bash
source ~/.bashrc
```

對於 zsh 用戶：

```bash
source ~/.zshrc
```

## 項目設置

### 1. 初始化 Go 模塊

```bash
go mod init concert-ticket-system
```

### 2. 安裝依賴

```bash
go get github.com/redis/go-redis/v9@v9.6.1
go get google.golang.org/grpc
```

### 3. 生成 gRPC 代碼

```bash
protoc --go_out=. --go-grpc_out=. pb/ticket_service.proto
```

## 使用說明
```bash
go run server/server.go
go run client/client.go
```
