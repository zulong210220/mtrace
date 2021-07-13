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
	"mtrace/message"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/service/grpc"
	bk "github.com/micro/go-plugins/broker/grpc"

	//"github.com/micro/go-plugins/client/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"
)

func main() {
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"localhost:2379",
		}
	})
	service := grpc.NewService(
		micro.Name("student.client"),
		micro.Registry(reg), //默认不在consul中注册服务，所以需要指定
		//micro.Client(grpc.NewClient()),
		//client.Selector(selector.NewSelector(selector.Registry(reg))),
		//selector.NewSelector(selector.Registry(reg)),
		micro.Broker(bk.NewBroker()),
		micro.RegisterTTL(30),
		micro.RegisterInterval(15),
		micro.Version("v1.0.0"),
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
