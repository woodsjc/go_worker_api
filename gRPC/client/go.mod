module github.com/woodsjc/worker_api/gRPC/client

go 1.15

require (
	github.com/woodsjc/worker_api/gRPC/internal/worker v0.0.0
	google.golang.org/grpc v1.38.0
)

replace github.com/woodsjc/worker_api/gRPC/internal/worker => ../internal/protos
