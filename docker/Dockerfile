FROM alpine:3.14
RUN apk update && apk add go vim git
COPY . /bot
WORKDIR /bot
RUN make build
CMD ["./main"]