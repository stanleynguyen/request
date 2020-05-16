FROM golang:alpine as builder
WORKDIR /go/src/github.com/stanleynguyen/request
RUN apk update && apk upgrade
COPY . .
RUN GOOS=linux go build -o request.out .

FROM node:alpine
RUN apk update && apk upgrade
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/stanleynguyen/request/index.js .
COPY --from=builder /go/src/github.com/stanleynguyen/request/request.out .
CMD ["./request.out"]