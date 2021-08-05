#build stage
FROM golang:1.16-stretch as builder

WORKDIR /betterog
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build  cmd/bog.go

FROM alpine:3.14

#final stage
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /betterog/bog /root/
COPY --from=builder /betterog/fonts /root/fonts/

EXPOSE 80

CMD ["./bog"]