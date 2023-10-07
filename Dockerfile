FROM golang:1.21-bullseye AS builder

RUN mkdir /first-take-bot
COPY ./ /first-take-bot/
WORKDIR /first-take-bot
RUN go build \
    && chmod 755 first-take-bot

FROM gcr.io/distroless/base

COPY --from=builder /first-take-bot/first-take-bot /app/first-take-bot
ENTRYPOINT ["/app/first-take-bot"]