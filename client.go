package main

/*
 * Author : lijinya
 * Email : yajin160305@gmail.com
 * File : client.go
 * CreateDate : 2021-04-12 10:34:18
 * */

import (
	"context"
	"fmt"
	"log"
	"mtrace/message"
	"mtrace/trace"
	"os"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/service/grpc"
	bk "github.com/micro/go-plugins/broker/grpc"
	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opencensus"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	rep "github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"

	//"github.com/micro/go-plugins/client/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"
)

func CInitTracer(zipkinURL string, hostPort string, serviceName string) rep.Reporter {
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
	CServiceName = "student.client"
	CZinkUrl     = "http://localhost:9411/api/v2/spans"
)

func main() {
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"localhost:2379",
		}
	})

	hostPort, _ := os.Hostname()

	reporter := CInitTracer(CZinkUrl, hostPort, CServiceName)
	defer reporter.Close()

	service := grpc.NewService(
		micro.Name(CServiceName),
		micro.Registry(reg), //默认不在consul中注册服务，所以需要指定
		//micro.Client(grpc.NewClient()),
		//client.Selector(selector.NewSelector(selector.Registry(reg))),
		//selector.NewSelector(selector.Registry(reg)),
		micro.Broker(bk.NewBroker()),
		micro.RegisterTTL(30),
		micro.RegisterInterval(15),
		micro.Version("v1.0.0"),
		micro.WrapHandler(trace.ServerWrapper),
		micro.WrapClient(wrapperTrace.NewClientWrapper()),
	)

	service.Init()
	//client.Selector(selector.NewSelector(selector.Registry(reg))),

	studentService := message.NewStudentService("student_service", service.Client())

	now := time.Now()
	for i := 0; i < 29; i++ {
		res, err := studentService.GetStudent(context.TODO(), &message.StudentRequest{Name: "davie"})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res.Name)
		fmt.Println(res.Classes)
		fmt.Println(res.Grade)
	}
	fmt.Println(time.Since(now))
}

/* vim: set tabstop=4 set shiftwidth=4 */
