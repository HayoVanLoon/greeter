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
	"crypto/tls"
	"flag"
	"fmt"
	pb "github.com/HayoVanLoon/genproto/greeter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

func getConn(host, port string) (*grpc.ClientConn, error) {
	addr := host + ":" + port

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithAuthority(host))

	if host == "localhost" {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// Cloud Run requires a TLS connection
		cred := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	return grpc.Dial(addr, opts...)
}

func doCall(conn *grpc.ClientConn, r *pb.GetGreetingRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cl := pb.NewGreeterClient(conn)

	resp, err := cl.GetGreeting(ctx, r)
	if err != nil {
		return fmt.Errorf("error getting greeting: %s", err)
	}

	fmt.Print(resp)
	return nil
}

func main() {
	var host = flag.String("host", "localhost", "host")
	var port = flag.String("port", "8080", "port")
	var name = flag.String("name", "Bob", "name")
	flag.Parse()

	r := &pb.GetGreetingRequest{Name: *name}

	conn, err := getConn(*host, *port)
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %s", err)
		}
	}()
	if err != nil {
		log.Fatalf("error opening connection: %s", err)
	}

	err = doCall(conn, r)
	if err != nil {
		log.Fatalf("error calling greeter: %s", err)
	}
}
