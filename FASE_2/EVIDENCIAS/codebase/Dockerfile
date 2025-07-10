FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go mod init github.com/gonzalohonorato/servercorego
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/servercorego ./main.go

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 80
ENTRYPOINT ["/go/bin/servercorego", "--port", "80"]
