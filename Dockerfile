FROM golang:1.8 as builder

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN CGO_ENABLED=0 go-wrapper install

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
CMD ["/app"]
