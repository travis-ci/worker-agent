# worker-agent

[WIP] A small go agent that runs inside the build vm and controls job execution

## Installation

We are using the [GRPC](https://grpc.io/docs/tutorials/basic/go.html) library.

```
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

GRPC also requires the `protoc` cli tool for protocol buffers to be installed, on macOS:

```
brew install protobuf
```

## Protoc compilation

If you change the grpc `.proto` file, you'll also need to re-generate the `.pb.go` file via:

```
protoc -I agent/ agent/agent.proto --go_out=plugins=grpc:agent
```
