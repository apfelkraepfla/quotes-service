FROM golang:1.21.4-alpine3.18
WORKDIR /app
COPY . .
RUN go install -v ./...
EXPOSE 3000
CMD ["cmd"]