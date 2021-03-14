FROM golang:1.15 as build-env

ADD . /go/src/todosync
WORKDIR /go/src/todosync

RUN go get -d -v ./...

RUN go build -o /go/bin/todosync ./cmd/todosync/todosync.go

# see https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/todosync /app
CMD ["/app"]