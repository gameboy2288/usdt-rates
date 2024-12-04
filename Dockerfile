FROM golang:1.22-alpine as builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download && go mod verify

RUN go build  -o ./app ./cmd/service

FROM alpine:3.14

RUN mkdir /app
WORKDIR /app

COPY --from=builder --chown=app:app /app/app .
COPY --from=builder --chown=app:app /app/db/ db/

CMD ["./app" ]
