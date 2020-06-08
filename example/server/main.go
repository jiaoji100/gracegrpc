/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"os"
	"log"
	"context"

	"google.golang.org/grpc"
	"github.com/jiaoji100/gracegrpc/gracegrpc"
	pb "github.com/jiaoji100/gracegrpc/gracegrpc/example/helloworld"
)

const (
	addr = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v,pid : %d", in.GetName(), os.Getpid())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	//lis, err := net.Listen("tcp", addr)
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}

	if err := gracegrpc.Serve(s, addr); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
