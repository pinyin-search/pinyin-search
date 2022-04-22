# build stage
FROM golang:1.17 AS builder

WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && make clean build

# final stage
FROM alpine
LABEL name=pinyin-search
LABEL url=https://github.com/pinyin-search/pinyin-search

WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /app/pinyin-search /app/pinyin-search
EXPOSE 9876
ENTRYPOINT ["/app/pinyin-search"]
