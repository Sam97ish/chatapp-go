# Chat application







### Complication
Use *.proto --go-grpc_out=./ --go_out=./ in the proto directory to create the service files.
docker build -t chatapp-server .
docker run -it -p 8080:8080 --rm --name  chatapp-server-go chatapp-server

	// Set up login
	wait.Add(1)
	go func() {
		defer wait.Done()
		for {
			loginView.name.OnSubmit(func(username *tui.Entry) {
				name = username.Text()
				ui.SetWidget(chatView.chat)
			})
			break
		}
		return
	}()


