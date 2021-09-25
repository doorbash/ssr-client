FROM golang:1.17.1-alpine3.13 as builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o /ssr-client

FROM scratch
COPY --from=builder /ssr-client /ssr-client
ENTRYPOINT [ "/ssr-client" ]
