package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Sam97ish/chatapp-go/proto/service"
	"github.com/google/uuid"
	"github.com/marcusolsson/tui-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

// Attempts to connect the user to the server.
// Returns a stream used to exchange messages.
func connect(user *service.User, client service.BroadcastClient) (service.Broadcast_CreateStreamClient, error) {
	var streamerror error
	conn := &service.Connect{
		User:   user,
		Active: true,
	}
	stream, err := client.CreateStream(context.Background(), conn)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}

	return stream, streamerror
}

func main() {
	// Args collection
	arguments := os.Args
	if len(arguments) == 1 {
		log.Fatal("Please run as  go run . host:port  [-s] .")
	}
	address := arguments[1]
	secure := false
	if len(arguments) == 3 {
		secure = !secure
	}

	// Get username
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Connecting to %s ...\n", address)
	fmt.Println("input username")
	fmt.Println("default is Anon (Press Enter)")
	fmt.Printf("connecting in secure mode: %v \n", secure)
	fmt.Println("---------------------")
	username, _ := reader.ReadString('\n')
	if username == "" || username == "\n" {
		username = "Anon"
	}
	username = strings.Replace(username, "\n", "", -1)

	// UI set up
	loginView := NewLoginView()
	chatView := NewChatView()

	ui, errLogin := tui.New(loginView.root)
	if errLogin != nil {
		log.Fatal(errLogin)
	}
	// Set up login
	loginView.input.OnSubmit(func(username *tui.Entry) {
		ui.SetWidget(chatView.chat)
	})

	quit := func() { ui.Quit() }
	ui.SetKeybinding("Esc", quit)
	ui.SetKeybinding("Ctrl+c", quit)

	// Set up
	done := make(chan int)
	id := uuid.New().String()

	// Set up connection options
	var connOptions grpc.DialOption
	if secure {

		cert, err := credentials.NewClientTLSFromFile(filepath.Join("cert", "ca-cert.pem"), "")
		if err != nil {
			log.Fatal("Error parsing the cert file for TLS support.")
		}
		connOptions = grpc.WithTransportCredentials(cert)
	} else {
		connOptions = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	// Attempt connection
	conn, err := grpc.Dial(address, connOptions)
	if err != nil {
		log.Fatalf("could not connect to service: %v", err)
	}

	client := service.NewBroadcastClient(conn)
	user := &service.User{
		Id:   id,
		Name: username,
	}

	// Connect user to server
	stream, errorStream := connect(user, client)
	if errorStream != nil {
		log.Fatalf("could not connect user: %v", errorStream)
	}

	// Incoming: receive msgs from server
	wait.Add(1)
	go func(str service.Broadcast_CreateStreamClient) {
		defer wait.Done()
		for {
			msg, errRec := str.Recv()
			if errRec != nil {
				_ = fmt.Errorf("error reading message: %v", errRec)
				break
			}

			ui.Update(func() { chatView.AddMessage(msg) })
		}
	}(stream)

	// Run chat UI
	wait.Add(1)
	go func() {
		defer wait.Done()
		if errUI := ui.Run(); errUI != nil {
			log.Fatal(errUI)
		}
	}()

	// Outgoing: send msgs to server
	wait.Add(1)
	go func() {
		defer wait.Done()

		for {
			chatView.input.OnSubmit(func(entry *tui.Entry) {
				timestamp := time.Now()
				msg := &service.Message{
					User:      user,
					Content:   entry.Text(),
					Timestamp: timestamp.Format(time.RFC1123Z),
				}
				chatView.input.SetText("")
				_, errBroad := client.BroadcastMessage(context.Background(), msg)
				if errBroad != nil {
					log.Fatalf("Error sending message: %v", errBroad)
					return
				}
			})
		}

	}()

	// Wait until all goroutines finish
	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
