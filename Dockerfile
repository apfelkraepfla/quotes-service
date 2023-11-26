FROM golang:1.21.4-alpine3.18
WORKDIR /app
COPY . .
RUN go mod download
RUN go install -v ./...
CMD ["app"]