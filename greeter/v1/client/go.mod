module github.com/HayoVanLoon/greeter/client

go 1.15

require (
	cloud.google.com/go v0.65.0 // indirect
	github.com/HayoVanLoon/genproto v0.0.0-20200907192710-ab6f07526f87
	google.golang.org/api v0.30.0
	google.golang.org/genproto v0.0.0-20200829155447-2bf3329a0021 // indirect
	google.golang.org/grpc v1.31.1
)

// replace github.com/HayoVanLoon/genproto => ../../../genproto
