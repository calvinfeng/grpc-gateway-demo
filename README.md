# gRPC Gateway Demo

## Protobuf

Sometimes I use Docker images to run proto compilation but it seems quite easy to do it on Mac with
brew. 

```bash
cd $GOPATH/src
brew install protobuf
go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u -v github.com/golang/protobuf/protoc-gen-go
```

I will put both protobuf files and protobuf compiled file into the same directory. The directory
name will be the package name for my Go files. Run the bash script and the outputs will be `pb.go`
and `pb.gw.go` for gRPC and gRPC reverse proxy.

```
./protogen.sh
```

## Error Handling

The mapping of gRPC status code to HTTP status code can be found here,

https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go#L15