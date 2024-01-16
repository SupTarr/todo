FROM cgr.dev/chainguard/go:latest AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

COPY . .

RUN go build \
    -ldflags "-linkmode external -extldflags -static" \
    -o api

FROM cgr.dev/chainguard/static:latest

# USER nonroot

# ENV TINI_VERSION v0.19.0
# ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static ./tini
# # RUN chmod +x ./tini

COPY --from=build /app/api .

EXPOSE 8081

USER nonroot:nonroot

CMD ["/api"]ß 