FROM golang:1.10
WORKDIR /go/src/github.com/glerchundi/grpc-helloworld

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go install -ldflags="-w -s" -v github.com/glerchundi/grpc-helloworld/cmd/helloworld-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/helloworld-server /bin/helloworld-server