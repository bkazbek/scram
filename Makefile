server:
	go run main.go serve
client:
	go run main.go client --address localhost:9001 --username tcp --password tcp_password