FROM golang:alpine as builder
COPY . /app
RUN go env -w GOPROXY=https://goproxy.cn,direct \
   && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
   && cd /app \
   && go build  -o ./export-server .

FROM alpine:latest 
WORKDIR /app
COPY --from=builder /app/export-server .
CMD [ "./export-server" ]