/*
 * Copyright 2019 Hayo van Loon
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

package main

import (
	"fmt"
	pb "github.com/HayoVanLoon/genproto/greeter"
	"github.com/HayoVanLoon/go-commons/logjson"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

type server struct {
}

func (s server) GetGreeting(ctx context.Context, r *pb.GetGreetingRequest) (*pb.Greeting, error) {
	var name string
	if r.GetName() == "" {
		name = "J. Doe"
	} else {
		name = r.GetName()
	}
	logjson.Info(fmt.Sprintf("Got greetings from %s", name))
	resp := &pb.Greeting{Text: fmt.Sprintf("Hello %s", name)}
	return resp, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, server{})
	reflection.Register(grpcServer)

	endpoint := ":" + port
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatalf("error listening: %s", err)
	}

	logjson.Notice(fmt.Sprintf("serving on %s", endpoint))
	log.Fatal(grpcServer.Serve(lis))
}
