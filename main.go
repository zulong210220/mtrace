package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mtrace/trace"
	"os"
	"time"

	"mtrace/message"
	"mtrace/ratelimit"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"

	"github.com/micro/go-micro/service/grpc"
	bk "github.com/micro/go-plugins/broker/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"

	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opencensus"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	rep "github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

//学生服务管理实现
type StudentManager struct {
}

//获取学生信息的服务接口实现
func (sm *StudentManager) GetStudent(ctx context.Context, request *message.StudentRequest, response *message.Student) error {

	//tom
	studentMap := map[string]message.Student{
		"davie":  message.Student{Name: "davie", Classes: "软件工程专业", Grade: 80},
		"steven": message.Student{Name: "steven", Classes: "计算机科学与技术", Grade: 90},
		"tony":   message.Student{Name: "tony", Classes: "计算机网络工程", Grade: 85},
		"jack":   message.Student{Name: "jack", Classes: "工商管理", Grade: 96},
	}

	if request.Name == "" {
		return errors.New(" 请求参数错误,请重新请求。")
	}

	//获取对应的student
	student := studentMap[request.Name]

	if student.Name != "" {
		fmt.Println(student.Name, student.Classes, student.Grade)
		*response = student
		return nil
	}

	return errors.New(" 未查询当相关学生信息 ")
}

func KaelReqWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
		fmt.Println("OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOooookkkk")
		os.Stdout.WriteString("OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOooookkkk\n")
		fmt.Println(req.ContentType())
		fmt.Println(req.Endpoint())
		fmt.Println(req.Service())
		fmt.Println(req.Method())
		fmt.Println(req.Header())
		fmt.Println(req.Body())
		return fn(ctx, req, rsp)
	}
}

func InitTracer(zipkinURL string, hostPort string, serviceName string) rep.Reporter {
	// set up a span reporter
	reporter := zipkinhttp.NewReporter(zipkinURL)
	//defer reporter.Close()

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(serviceName, hostPort)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)

	return reporter
}

const (
	ServiceName = "student_service"
	//ZinkUrl     = "http://localhost:9411/api/v2/spans"
	ZinkUrl = "http://zipkin.knowbox.cn/api/v2/spans"
)

func main() {

	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"localhost:2379",
		}
	})

	hostPort, _ := os.Hostname()

	reporter := InitTracer(ZinkUrl, hostPort, ServiceName)
	defer reporter.Close()

	//创建一个新的服务对象实例
	service := grpc.NewService(
		micro.Name(ServiceName),
		micro.Registry(reg), //默认不在consul中注册服务，所以需要指定
		micro.Broker(bk.NewBroker()),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(15*time.Second),
		micro.Version("v1.0.0"),
		micro.WrapHandler(KaelReqWrapper, ratelimit.NewHandlerWrapper(2)),
		micro.WrapHandler(trace.ServerWrapper),
		micro.WrapClient(wrapperTrace.NewClientWrapper()),
	)

	//服务初始化
	service.Init(
		micro.BeforeStart(func() error {
			log.Println("000000")
			return nil
		}),
	)

	//注册
	message.RegisterStudentServiceHandler(service.Server(), new(StudentManager))

	//运行
	err := service.Run()
	log.Println("RUN.....")
	if err != nil {
		log.Fatal(err)
	}
}
