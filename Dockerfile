FROM alpine:3.6

RUN apk update && apk add curl

CMD hello
