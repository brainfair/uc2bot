FROM golang:latest
WORKDIR /src
RUN go get -d -v github.com/go-telegram-bot-api/telegram-bot-api
RUN go get -d -v github.com/go-sql-driver/mysql
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /src/app .
CMD ["./app"]  