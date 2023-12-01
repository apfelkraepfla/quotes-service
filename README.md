# Quotes-Service

## Description
The intent of this project is to become more familiar with Golang and GitHub.

The quotes service is a very simple gRPC service in Golang. For development I experiment with VS Codes devcontainers, and Tilt (with Minikube). It started as a part of the project https://github.com/apfelkraepfla/tilt-my-dev that grew over its initial purpose, and became too convoluted. 

## Deployment via Tilt in Minikube
The Quotes-Service can be deployed (e.g. to a local Minikube Kubernetes) by using Tilt. There is already a setup in the project https://github.com/apfelkraepfla/tilt-my-dev, which assumes that you have both projects cloned locally, and that they are located next to each other on your filesystem. 
When the deployment worked, you can port-forward its port, and then communicate with the server:
```
$ kubectl port-forward <POD_NAME> 3000:3000
```

In another console, you can communicate with the server:
```
$ grpcurl -plaintext localhost:3000 list
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
quotes.QuoteService


$ grpcurl -plaintext -d '{}' localhost:3000 quotes.QuoteService/GetQuote
```

## Added features
* devcontainer
* gRPC server, incl. vscode task to generate compiled files
* checking out grpcurl
* adding reflection to be able to list functions with grpcurl
* adding REST API besides gRPC:
  * adding annotations in `.proto`, add necessary `google.api.*` files, and extend the protoc command in the vscode task file to generate the `*.pb.gw.go` file 