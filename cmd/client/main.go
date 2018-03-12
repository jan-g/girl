package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jan-g/girl/model"
	"google.golang.org/grpc"
)

func main() {
	server := os.Args[1]
	facet := os.Args[2]
	qty, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}

	start := time.Now()
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := model.NewLimiterProtocolClient(conn)

	result, err := client.Use(context.Background(), &model.UseRequest{Facet: facet, Quantity: int64(qty)})
	conn.Close()
	end := time.Now()
	fmt.Println("Facet:", result.Facet, "Taken:", result.Quantity, "Remaining:", result.Remaining)
	fmt.Println("Time elapsed/ns:", end.Sub(start))
}
