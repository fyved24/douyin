FROM golang:1.20 AS build
WORKDIR /home/douyin
COPY . .

ENV CGO_ENABLED=0
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o server .

FROM alpine:3.17 AS run
COPY --from=build /home/douyin/server /server
CMD [ "/server" ]