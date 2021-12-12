package main

import (
	"context"
	"fmt"
	helloworld "go101/rpc/grpc"
	"go101/rpc/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"reflect"
)

func main() {
	srv := &service.HelloWorldService{}
	grpcMethodDescs := make([]grpc.MethodDesc, 0)
	rType := reflect.TypeOf(srv)
	for i := 0; i < rType.NumMethod(); i++ {
		method := rType.Method(i)
		if isGRPCHandler(method) {
			in := reflect.New(method.Type.In(2).Elem()).Interface()
			grpcMethodDescs = append(grpcMethodDescs, grpc.MethodDesc{
				MethodName: method.Name,
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return invoke(srv, method, ctx, in)
					}
					info := &grpc.UnaryServerInfo{
						Server:     srv,
						FullMethod: "/helloworld.Greeter/" + method.Name, //FIXME
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return invoke(srv, method, ctx, req)
					}
					return interceptor(ctx, in, info, handler)
				},
			})
		}
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("UnaryInterceptor")
		return handler(ctx, req)
	}))
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "helloworld.Greeter",             //FIXME
		HandlerType: (*helloworld.GreeterServer)(nil), //FIXME
		Methods:     grpcMethodDescs,
		//Streams:     []grpc.StreamDesc{},
		//Metadata:    "helloworld.proto",
	}, srv)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func isGRPCHandler(method reflect.Method) bool {
	rType := method.Type
	if rType.NumIn() != 3 || rType.NumOut() != 2 {
		return false
	}
	ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
	pbMsgType := reflect.TypeOf((*proto.Message)(nil)).Elem()
	errType := reflect.TypeOf((*error)(nil)).Elem()
	if !rType.In(1).Implements(ctxType) || !rType.In(2).Implements(pbMsgType) {
		return false
	}
	if !rType.Out(0).Implements(pbMsgType) || !rType.Out(1).Implements(errType) {
		return false
	}
	return true
}

func invoke(srv interface{}, method reflect.Method, ctx context.Context, in interface{}) (interface{}, error) {
	outs := method.Func.Call([]reflect.Value{
		reflect.ValueOf(srv),
		reflect.ValueOf(ctx),
		reflect.ValueOf(in),
	})
	if len(outs) != 2 {
		return nil, fmt.Errorf("len of outs: %d", len(outs))
	}

	var resp proto.Message
	if out := outs[0].Interface(); out != nil {
		resp = out.(proto.Message)
	}
	var err error
	if out := outs[1].Interface(); out != nil {
		err = out.(error)
	}
	return resp, err
}
