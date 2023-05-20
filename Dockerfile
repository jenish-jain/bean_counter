FROM golang:1.19

WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./

RUN go mod download

# Copy local code to the container image.
COPY . ./

RUN go build cmd/bean_counter/main.go

EXPOSE 8080

CMD [ "./main" ]