FROM registry.cn-hangzhou.aliyuncs.com/mw5uk4snmsc/go:1.21.6-alpine3.18 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

RUN set -ex \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk --update add tzdata \
    && apk add --update ttf-dejavu fontconfig \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && apk --no-cache add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && mkfontscale \
    && mkfontdir \
    && fc-cache

WORKDIR /opt

COPY . .

RUN go mod download && go mod tidy -v && go build -o go-devops-observability .

FROM registry.cn-hangzhou.aliyuncs.com/mw5uk4snmsc/go:1.21.6-alpine3.18

WORKDIR /opt

RUN set -ex \
   && mkdir -pv /opt/etc

COPY --from=builder /opt/go-devops-observability /opt/
COPY --from=builder /opt/etc /opt/etc

EXPOSE 8080

ENTRYPOINT ["./go-devops-observability"]
