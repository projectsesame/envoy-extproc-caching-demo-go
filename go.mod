module github.com/projectsesame/envoy-extproc-caching-demo-go

go 1.21

require (
	github.com/envoyproxy/go-control-plane v0.12.0
	github.com/wrossmorrow/envoy-extproc-sdk-go v0.0.21
)

require (
	github.com/cncf/xds/go v0.0.0-20231128003011-0fa0005c9caa // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240228224816-df926f6c8641 // indirect
	google.golang.org/grpc v1.62.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace github.com/wrossmorrow/envoy-extproc-sdk-go => github.com/izturn/envoy-extproc-sdk-go v0.0.3
