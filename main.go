package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/chacerapp/apiserver/server"
	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/chacerapp/apiserver/store"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Opening postgres database")
	db, err := sql.Open("postgres", "postgres://root@localhost:26257/chacerapp_tests?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to open new database connection: %v", err)
	}

	log.Print("Dialing TCP address for gRPC server")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to get TCP listener: %v", err)
	}

	log.Print("Creating a new gRPC server")
	srv := server.NewGRPCServer(store.New(db, store.NewPaginator([]byte("my-super-secure-test-secret-3234"))))

	go func() {
		log.Print("Starting gRPC server")
		err := srv.Serve(listener)
		if err != nil {
			log.Fatalf("error when running gRPC server: %v", err)
		}
	}()

	dialCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Dialing API Server")
	dial, err := grpc.DialContext(dialCtx, ":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial gRPC server: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("calling ListAccounts endpoint")

	client := serverpb.NewAccountsClient(dial)
	_, err = client.ListAccounts(ctx, &serverpb.ListAccountsRequest{
		PageSize: 50,
	})
	if err != nil {
		log.Fatalf("failed to list accounts: %v", err)
	}

	locationClient := serverpb.NewLocationsClient(dial)
	_, err = locationClient.ListLocations(ctx, &serverpb.ListLocationsRequest{
		Parent:   "accounts/my-account",
		PageSize: 50,
	})
	if err != nil {
		log.Fatalf("failed to list locations: %v", err)
	}
}
