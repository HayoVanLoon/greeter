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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	pb "github.com/HayoVanLoon/genproto/hayovanloon/greeter/v1"
	"github.com/HayoVanLoon/go-commons/logjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"time"
)

const ForwardedUserInfoHeader = "X-Endpoint-API-UserInfo"

type server struct {
	pb.UnimplementedGreeterServer
	cache map[string]*pb.Greeting
}

func NewServer() *server {
	return &server{cache: make(map[string]*pb.Greeting)}
}

func LogPanic() {
	if r := recover(); r != nil {
		logjson.Error(r)
		time.Sleep(5 * time.Second)
	}
}

func (s *server) CreateGreeting(ctx context.Context, r *pb.CreateGreetingRequest) (*pb.Greeting, error) {
	logjson.Debug("CreateGreeting")
	defer LogPanic()

	if r.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "no name given")
	}

	g := &pb.Greeting{Text: fmt.Sprintf("Hello %s.", r.Name)}
	s.cache[r.Name] = g
	return g, nil
}

func (s *server) GetGreeting(ctx context.Context, r *pb.GetGreetingRequest) (*pb.Greeting, error) {
	logjson.Debug("GetGreeting")
	defer LogPanic()

	u := &UserInfo{}
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		if b64 := headers.Get(ForwardedUserInfoHeader); len(b64) > 0 {
			bs, err := base64.RawStdEncoding.DecodeString(b64[0])
			if err != nil {
				logjson.Warn("could not decode " + b64[0])
			} else {
				logjson.Debug(string(bs))
				err = json.Unmarshal(bs, u)
				if err != nil {
					logjson.Warn(fmt.Sprintf("unexpected user info format %s", err))
				}
				logjson.Debug(u)
			}
		}
	}

	t := ""
	if g, ok := s.cache[r.Name]; ok {
		t = g.Text
	} else {
		t = fmt.Sprintf("Hello %s, you came unexpected.", r.Name)
	}

	resp := &pb.Greeting{Text: fmt.Sprintf("%s (%s / %s / %s)", t, u.Issuer, u.Id, u.Email)}
	return resp, nil
}

func (s *server) ListGreetings(ctx context.Context, r *pb.ListGreetingsRequest) (*pb.ListGreetingsResponse, error) {
	logjson.Debug("ListGreetings")
	defer LogPanic()

	resp := &pb.ListGreetingsResponse{}
	for _, g := range s.cache {
		resp.Greetings = append(resp.Greetings, g)
	}
	return resp, nil
}

func (s *server) ListHugs(ctx context.Context, r *pb.ListHugsRequest) (*pb.ListHugsResponse, error) {
	logjson.Debug("ListHugs")
	defer LogPanic()

	return &pb.ListHugsResponse{Hugs: []*pb.Hug{{}, {}, {}}}, nil
}

func main() {
	defer LogPanic()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	grpcServer := grpc.NewServer()
	s := NewServer()
	pb.RegisterGreeterServer(grpcServer, s)

	addr := ":" + port
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error listening: %s", err)
	}

	logjson.Notice(fmt.Sprintf("serving on %s", addr))
	if err = grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
