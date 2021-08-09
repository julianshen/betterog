#build stage
FROM golang:1.16-stretch as builder

WORKDIR /betterog
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build  cmd/bog.go

FROM chromedp/headless-shell:latest
#final stage
# RUN apk --no-cache add ca-certificates
RUN apt update
RUN apt install dumb-init
RUN apt-get update -y \
    && apt-get install -y fonts-noto \
    && apt-get install -y fonts-noto-cjk

WORKDIR /root/
COPY --from=builder /betterog/bog /root/
COPY --from=builder /betterog/fonts /root/fonts/
COPY --from=builder /betterog/static /root/static/
COPY --from=builder /betterog/templates /root/templates/

EXPOSE 80
ENTRYPOINT ["dumb-init", "--"]
CMD ["./bog"]