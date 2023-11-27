FROM golang:1.21.4

ENV NVM_DIR="/usr/local/share/nvm"
ENV NVM_SYMLINK_CURRENT=true \
    PATH=${NVM_DIR}/current/bin:${PATH}

ARG VSCODE_SCRIPTS_VERSION="v0.245.2"
ARG NODE_VERSION="10"
# Run some common installation scripts for a nicer dev environment. In order:
# Used to create non-root user and update system packages
# https://github.com/microsoft/vscode-dev-containers/blob/main/script-library/docs/common.md
# We use this to install Go tools used by gopls
# https://github.com/microsoft/vscode-dev-containers/blob/main/script-library/docs/go.md
# We use this to install Node
# https://github.com/microsoft/vscode-dev-containers/blob/main/script-library/docs/node.md
RUN apt-get update && \
    wget "https://raw.githubusercontent.com/microsoft/vscode-dev-containers/${VSCODE_SCRIPTS_VERSION}/script-library/common-debian.sh" && \
    chmod +x ./common-debian.sh 
RUN ./common-debian.sh false vscode automatic automatic true false 
RUN  wget "https://raw.githubusercontent.com/microsoft/vscode-dev-containers/${VSCODE_SCRIPTS_VERSION}/script-library/go-debian.sh" && \
    chmod +x ./go-debian.sh && \
    ./go-debian.sh none /usr/local/go /go vscode false true 
RUN wget "https://raw.githubusercontent.com/microsoft/vscode-dev-containers/${VSCODE_SCRIPTS_VERSION}/script-library/node-debian.sh" && \
    chmod +x ./node-debian.sh && \
    ./node-debian.sh "${NVM_DIR}" "${NODE_VERSION}" vscode true && \
    rm common-debian.sh go-debian.sh 
# node-debian.share
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
    wget \
    unzip \
    openjdk-17-jre-headless \
    bzip2 \
    protobuf-compiler \
    patch && \
    apt-get clean -y && \
    rm -rf /var/lib/apt/lists/*

# Install swagger-codegen
ENV SWAGGER_CODEGEN_VERSION=2.4.8
RUN wget https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_CODEGEN_VERSION}/swagger-codegen-cli-${SWAGGER_CODEGEN_VERSION}.jar \
    -O /usr/local/bin/swagger-codegen-cli.jar && \
    echo '#!/bin/bash\njava -jar /usr/local/bin/swagger-codegen-cli.jar "$@"' > /usr/local/bin/swagger-codegen && \
    chmod +x /usr/local/bin/swagger-codegen

# Install Bazelisk as bazel to manage Bazel
RUN go install github.com/bazelbuild/bazelisk@latest && \
    mv $(which bazelisk) /usr/local/bin/bazel

# Install buildifier for bazel formatting
RUN go install github.com/bazelbuild/buildtools/buildifier@latest

# Give vscode ownership of GOPATH
RUN chown -R vscode: /go

# To format proto file with vscode-proto3 extension
RUN apt update && apt install clang-format --yes

# Install the compiler for protocol buffer
RUN curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v22.3/protoc-22.3-linux-x86_64.zip" && \
    unzip protoc-22.3-linux-x86_64.zip -d /usr/local && \
    rm protoc-22.3-linux-x86_64.zip

USER vscode

RUN go install mvdan.cc/gofumpt@v0.5.0 && \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.1 && \
    # To generate gRPC files from proto for Golang
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
    export PATH="$PATH:$(go env GOPATH)/bin"