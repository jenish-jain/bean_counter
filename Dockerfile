FROM golang:1.21-alpine as builder

WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./

RUN go mod download

# Copy local code to the container image.
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" cmd/bean_counter/main.go

FROM scratch

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

COPY --from=builder /app/main main

EXPOSE 8080

CMD [ "./main" ]