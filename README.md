# Greeter

A minimal gRPC server that runs on Cloud Run.

## Local Server
```shell script
make -C greeter/v1 build docker-run
```

## Local Testing
```shell script
make -C greeter/v1 smoke-test-local
```

## Deploying to GCP
### Deploying the Backend
```shell script
make -C greeter/v1 all PROJECT=<your project>
```

### Deploying the Gateway
Requires backend to be deployed.
```shell script
make -C gateway all PROJECT=<your project>
```

## Testing
```shell script
make -C greeter/v1 smoke-test-cloud PROJECT=<your project>
```

## Deploying a new release to GCP
### Deploying the Backend
```shell script
make -C greeter/v1 release PROJECT=<your project>
```

### Deploying the Gateway
Required on API / endpoint changes
```shell script
make -C gateway release PROJECT=<your project>
```
