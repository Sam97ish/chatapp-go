package main

import (
	"bufio"
	"chatapp-go/proto/service"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"sync"
	"time"
)

var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *service.User, client service.BroadcastClient) error {
	var streamerror error
	conn := &service.Connect{
		User:   user,
		Active: true,
	}
	stream, err := client.CreateStream(context.Background(), conn)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	wait.Add(1)
	go func(str service.Broadcast_CreateStreamClient) {
		defer wait.Done()
		for {
			msg, errRec := str.Recv()
			if errRec != nil {
				streamerror = fmt.Errorf("error reading message: %v", errRec)
				break
			}
			msgTime, errTime := time.Parse(time.RFC1123Z, msg.Timestamp)
			if errTime != nil {
				fmt.Printf("error parsing time: %v", errTime)
			}
			fmt.Printf("<%s> %v: %s\n", msgTime.Format(time.Stamp), msg.User.Name, msg.Content)
		}
	}(stream)

	return streamerror
}
func main() {
	arguments := os.Args
	if len(arguments) == 2 {
		log.Fatal("Please run as [ go run client.go host:port name ].")
	}
	address := arguments[1]
	name := arguments[2]

	timestamp := time.Now()
	done := make(chan int)

	id := uuid.New().String()

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to service: %v", err)
	}
	//fmt.Printf("conn %v", conn)
	client := service.NewBroadcastClient(conn)

	user := &service.User{
		Id:   id,
		Name: name,
	}

	errUser := connect(user, client)
	if errUser != nil {
		log.Fatalf("could not connect user: %v", errUser)
	}

	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			msg := &service.Message{
				User:      user,
				Content:   scanner.Text(),
				Timestamp: timestamp.Format(time.RFC1123Z),
			}
			_, errBroad := client.BroadcastMessage(context.Background(), msg)
			if errBroad != nil {
				fmt.Printf("Error sending message: %v", errBroad)
				break
			}
		}

	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
