# Greeter

A minimal gRPC server that runs on Cloud Run.

## Local Server
```shell script
make -C greeter/v1 build docker-run
```

## Local Testing
```shell script
make -C v1 smoke-test-local
```

## Deploying to GCP
```shell script
make -C greeter/v1 deploy PROJECT=<your project>
```

## Testing
```shell script
make -C greeter/v1 smoke-test-cloud PROJECT=<your project>
```

## Deploying a new release to GCP
```shell script
make -C greeter/v1 release PROJECT=<your project>
```
