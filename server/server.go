package main

import (
	"context"
	"crypto/tls"
	"github.com/Sam97ish/chatapp-go/proto/service"
	"google.golang.org/grpc"
	"net"
	"os"
	"path/filepath"
	"sync"

	glog "google.golang.org/grpc/grpclog"
)

var logger glog.LoggerV2

// Logger set up
func init() {
	logger = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

// Connection Represents a connection on the server
type Connection struct {
	stream service.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

// Server Represents all the connections on the server
type Server struct {
	Connections []*Connection
	service.UnimplementedBroadcastServer
}

// CreateStream Creates a stream used to exchange data
func (s *Server) CreateStream(pconn *service.Connect, stream service.Broadcast_CreateStreamServer) error {
	newconn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connections = append(s.Connections, newconn)
	logger.Infof("User Id %s connected", pconn.User.Id)
	return <-newconn.error
}

// BroadcastMessage Sends message to all connected users
func (s *Server) BroadcastMessage(ctx context.Context, msg *service.Message) (*service.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connections {
		wait.Add(1)

		go func(msg *service.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				logger.Infof("Sending message to: %s", conn.id)

				if err != nil {
					logger.Errorf("Error with Stream %s - Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}

		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &service.Close{}, nil
}

func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		logger.Fatal("Please provide host:port [-s] .")
	}
	address := arguments[1]
	secure := false
	if len(arguments) == 3 {
		secure = !secure
	}

	var connections []*Connection

	server := &Server{connections, service.UnimplementedBroadcastServer{}}

	grpcServer := grpc.NewServer()
	var listener net.Listener
	var err error
	if secure {
		cert, errCert := tls.LoadX509KeyPair(filepath.Join("cert", "server-cert.pem"), filepath.Join("cert", "server-key.pem"))
		if errCert != nil {
			logger.Fatalf("Error loading the server cert %s", err)
		}

		config := &tls.Config{Certificates: []tls.Certificate{cert}, ServerName: "test"}
		listener, err = tls.Listen("tcp", address, config)
		if err != nil {
			logger.Fatalf("Error creating the server %s", err)
		}
	} else {
		listener, err = net.Listen("tcp", address)
		if err != nil {
			logger.Fatalf("Error creating the server %s", err)
		}
	}

	logger.Infof("Starting server at address %s", address)

	service.RegisterBroadcastServer(grpcServer, server)

	errListen := grpcServer.Serve(listener)
	if errListen != nil {
		logger.Fatalf("Error while listening %s", errListen)
	}
}
