clean:
	rm -f protos/go/**/*.pb.go
proto:
	protoc --go_out=. \
		--go-grpc_out=. \
		protos/banking.proto