/*
 * Copyright 2020 Hayo van Loon
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
	pb "github.com/HayoVanLoon/genproto/hayovanloon/greeter"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
	"log"
	"strings"
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

func getClient(host, port string) (pb.GreeterClient, func(), error) {
	conn, err := getConn(host, port)
	closeFn := func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %s", err)
		}
	}
	if err != nil {
		return nil, nil, fmt.Errorf("error opening connection: %s", err)
	}
	cl := pb.NewGreeterClient(conn)
	return cl, closeFn, nil
}

func addIdToken(ctx context.Context, aud string) (context.Context, error) {
	tokenSource, err := idtoken.NewTokenSource(ctx, aud)
	if err != nil {
		return nil, fmt.Errorf("error creating token source: %v", err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("error creating token: %v", err)
	}
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token.AccessToken)
	return ctx, nil
}

func createContext(aud string, skipAuth bool) (context.Context, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if !skipAuth && !strings.HasPrefix(aud, "localhost") {
		// TODO(hvl): check audience
		var err error
		ctx, err = addIdToken(ctx, aud)
		if err != nil {
			return nil, nil, err
		}
	}
	return ctx, cancel, nil
}

func authed(cl pb.GreeterClient, host, name string, skipAuth bool) {
	ctx, cancel, err := createContext(host, skipAuth)
	if err != nil {
		fmt.Printf("error creating context: %s\n", err)
		return
	}
	defer cancel()

	resp, err := cl.CreateGreeting(ctx, &pb.CreateGreetingRequest{Name: "Alice"})
	if err != nil {
		fmt.Printf("error creating greeting: %s\n", err)
	} else {
		fmt.Println(resp)
	}

	resp, err = cl.CreateGreeting(ctx, &pb.CreateGreetingRequest{Name: ""})
	if err != nil {
		fmt.Printf("error creating greeting: %s\n", err)
	} else {
		fmt.Println(resp)
	}
	resp, err = cl.GetGreeting(ctx, &pb.GetGreetingRequest{Name: name})
	if err != nil {
		fmt.Printf("error getting greeting: %s\n", err)
	} else {
		fmt.Println(resp)
	}

	resp, err = cl.GetGreeting(ctx, &pb.GetGreetingRequest{Name: "Alice"})
	if err != nil {
		fmt.Printf("error getting greeting: %s\n", err)
	} else {
		fmt.Println(resp)
	}

	resp3, err := cl.ListGreetings(ctx, &pb.ListGreetingsRequest{})
	if err != nil {
		fmt.Printf("error listing greetings: %s\n", err)
	} else {
		fmt.Println(resp3)
	}
}

func unauthed(cl pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp4, err := cl.ListHugs(ctx, &pb.ListHugsRequest{})
	if err != nil {
		fmt.Printf("error getting hugs: %s\n", err)
	} else {
		fmt.Println(resp4)
	}
}

func main() {
	var host = flag.String("host", "localhost", "gRPC host")
	var port = flag.String("port", "8080", "gRPC host port")
	var skipAuth = flag.Bool("skip-auth", false, "do not add ID token")
	var name = flag.String("name", "Bob", "name")
	flag.Parse()

	cl, closeFn, err := getClient(*host, *port)
	if err != nil {
		log.Fatalf("error creating client: %s\n", err)
	}
	defer closeFn()

	authed(cl, *host, *name, *skipAuth)
	unauthed(cl)
}
