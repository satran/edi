FROM golang:1-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /dabba

## Deploy
FROM alpine:3.14

WORKDIR /

COPY --from=build /dabba /bin/dabba
RUN addgroup -S dabba -g 1000 && adduser -S dabba -G dabba -u 1000
RUN mkdir -p /home/dabba/dabba
RUN chown -R dabba:dabba /home/dabba
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

EXPOSE 8080

USER dabba:dabba

ENTRYPOINT ["dabba", "server", "-addr", "0.0.0.0:8080"]
