# 1. build
FROM golang:1.24.2-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app cmd/api/main.go

# 2. run
FROM scratch
COPY --from=build /app /app
COPY --from=build /src/configs /configs
COPY --from=build /src/static /static
EXPOSE 8080
ENTRYPOINT ["/app"]