FROM golang:1.21-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

ENV GOARCH=amd64

RUN go build \
  -ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
  -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
  -o /go/bin/app

FROM alpine:latest

COPY --from=build /go/bin/app /app

EXPOSE 8081

CMD ["/app"]