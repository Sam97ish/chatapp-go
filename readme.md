# Chat application
A fun little chat application I created using Go and gRPC.
The client has neat TUI that allows chatting.
The chat is supposed to be anonymous only so no login.

# Usage
## Client
You can run the client by either running `go run . [host:port]` in client folder or running the bash script at the root directory.
The client will try to connect to the server.

## Server
You can either run `go run . [host:port]` in server folder or use docker with:
- docker build -t chatapp-server .
- docker run -it -p 8080:8080 --rm --name  chatapp-server-go chatapp-server
you can also use the bash script to autorun the server.


# Compilation
Use `go build` on either client or server to build them.
The server is dockerized and can be built there.
you can regenerate the proto files using 
`*.proto --go-grpc_out=./ --go_out=./`




