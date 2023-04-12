# development stage
FROM golang:1.20-alpine as development

WORKDIR /usr/app

RUN apk add --no-cache pkgconfig gcc imagemagick imagemagick-dev libc-dev

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

# build stage
FROM golang:1.20-alpine as builder

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app
RUN apk add --no-cache --virtual .build-deps build-base tzdata
RUN go get -v
RUN go build -o main

# final stage
FROM alpine as production
ENV TZ=Asia/Novosibirsk
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/src/app/main /usr/app/main

WORKDIR /usr/app
RUN mkdir -p /tmp
ENTRYPOINT /usr/app/main
