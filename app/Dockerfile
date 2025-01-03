# Latest golang image on apline linux
FROM golang:1.23 AS builder

# Env variables
ENV GOOS=linux
ENV CGO_ENABLED=0

# Work directory
WORKDIR /docker-go

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Building the application
RUN go build -o docker-go


# Fetching the latest nginx image
FROM alpine:3.16 AS production

# Certificates
RUN apk add --no-cache ca-certificates

# Copying built assets from builder
COPY --from=builder docker-go .

# Starting our application
CMD ./docker-go

# Exposing server port
EXPOSE 8080
