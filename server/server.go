package main

import (
	"chatapp-go/proto/service"
	"context"
	"google.golang.org/grpc"
	"net"
	"os"
	"sync"

	glog "google.golang.org/grpc/grpclog"
)

var logger glog.LoggerV2

func init() {
	logger = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
	stream service.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Server struct {
	Connections []*Connection
	service.UnimplementedBroadcastServer
}

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
		logger.Fatal("Please provide host:port.")
	}
	address := arguments[1]

	var connections []*Connection

	server := &Server{connections, service.UnimplementedBroadcastServer{}}

	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("Error creating the server %s", err)
	}

	logger.Infof("Starting server at address %s", address)

	service.RegisterBroadcastServer(grpcServer, server)

	errListen := grpcServer.Serve(listener)
	if errListen != nil {
		logger.Fatalf("Error while listening %s", errListen)
	}
}
