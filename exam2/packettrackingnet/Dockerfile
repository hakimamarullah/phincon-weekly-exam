FROM golang:1.20-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o packettrackingnet

ENTRYPOINT [ "/app/packettrackingnet" ]