module mtrace

go 1.14

require (
	github.com/golang/protobuf v1.5.2
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/broker/grpc v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/registry/etcdv3 v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/trace/opencensus v0.0.0-20200119172437-4fe21aa238fd
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.1
	go.uber.org/ratelimit v0.2.0
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa
	google.golang.org/grpc v1.26.0
	google.golang.org/protobuf v1.27.1
)
