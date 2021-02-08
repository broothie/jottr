FROM golang:1.15-alpine AS builder
WORKDIR /go/src/github.com/broothie/jottr
COPY . .
RUN apk add --update ca-certificates
RUN go build cmd/server/main.go

FROM alpine:3.7
COPY templates templates
COPY public public
COPY --from=builder /go/src/github.com/broothie/jottr/main main
CMD ./main
