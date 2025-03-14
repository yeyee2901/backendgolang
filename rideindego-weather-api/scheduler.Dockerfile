FROM golang:1.22.5-alpine3.19 AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./scheduler ./cmd/scheduler/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/scheduler ./scheduler
COPY --from=builder /build/build.env ./build.env

CMD [ "./scheduler" ]
