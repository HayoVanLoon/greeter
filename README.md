# Greeter

A minimal gRPC server that runs on Cloud Run.

## Local Server
```shell script
make -C v1 build docker-run
```

## Local Testing
```shell script
make -C v1 smoke-test-local
```

## Deploying to GCP
```shell script
make -C v1 build push-gcr service-account deploy PROJECT=<your project>
```

## Testing
```shell script
make -C v1 smoke-test-cloud PROJECT=<your project>
```
