FROM golang:1-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /edi

## Deploy
FROM alpine:3.14

WORKDIR /

COPY --from=build /edi /bin/edi
RUN addgroup -S edi -g 1000 && adduser -S edi -G edi -u 1000
RUN mkdir -p /home/edi/edi
RUN chown -R edi:edi /home/edi
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

EXPOSE 8080

USER edi:edi

ENTRYPOINT ["edi", "server", "-addr", "0.0.0.0:8080"]
