gen:
	protoc --proto_path=chatpb chatpb/*.proto --go_out=chatpb --go_opt=paths=source_relative \
		--go-grpc_out=chatpb --go-grpc_opt=paths=source_relative
webserver:
	go run chat_webserver/webserver.go
gRPC:
	go run chat_server/server.go
