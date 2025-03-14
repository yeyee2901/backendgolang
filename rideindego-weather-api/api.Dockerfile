FROM golang:1.22.5-alpine3.19 AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/api/main.go --markdownFiles swagger-markdown --parseDependency true
RUN go build -o ./api-gateway ./cmd/api/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/api-gateway ./api-gateway
COPY --from=builder /build/build.env ./build.env

CMD ["./api-gateway"]
