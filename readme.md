# Chat application







### Complication
Use *.proto --go-grpc_out=./ --go_out=./ in the proto directory to create the service files.
docker build -t chatapp-server .
docker run -it -p 8080:8080 --rm --name  chatapp-server-go chatapp-server

