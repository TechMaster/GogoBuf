package main

import (
	"context"
	"log"

	proto "github.com/TechMaster/GogoBuf/proto"
	"github.com/micro/go-micro"
)

type Huy struct{}

func (g *Huy) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func (g *Huy) GoodBye(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {

	return nil
}

func (g *Huy) VuiQua(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
	)

	service.Init()

	proto.RegisterGreeterHandler(service.Server(), new(Huy))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
