FROM alpine:3.14
RUN apk update && apk add go vim git make
WORKDIR /bot
COPY . .
RUN make build
CMD ["./main"]