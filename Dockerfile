FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o webservice cmd/webservice/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/webservice .

FROM alpine:3.10

EXPOSE 12345

ARG UID=1001
ARG USER=service
ARG GID=1001
ARG GROUP=service

RUN addgroup -g $GID -S $GROUP && \
    adduser -u $UID -S $USER -G $GROUP && \
    mkdir -p /service/config /service/secrets && \
    chown -R $USER:$GROUP /service

WORKDIR /service

USER $USER
COPY --from=builder /dist/webservice .

CMD ./webservice "/service/config/config.json" "/service/secrets"