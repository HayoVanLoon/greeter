# Greeter

A minimal gRPC server that runs on Cloud Run.

```shell script
make -C v1 build
make -C v1 push-gcr PROJECT=<your project>
make -C v1 service-account PROJECT=<your project>
make -C v1 deploy PROJECT=<your project>
```
