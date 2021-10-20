FROM alpine:latest
RUN apk update && apk add go vim git
RUN git clone https://github.com/rombintu/goxmppbot.git
WORKDIR /goxmppbot
COPY ./config.toml /goxmppbot/config.toml
RUN go build main.go
CMD ["./main"]