#https://basefas.github.io/2019/09/24/%E4%BD%BF%E7%94%A8%20Docker%20%E6%9E%84%E5%BB%BA%20Go%20%E5%BA%94%E7%94%A8/
FROM golang:1.14.5 as mod
LABEL stage=mod
ARG GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://goproxy.io,direct
WORKDIR /root/myapp/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

FROM mod as builder
LABEL stage=intermediate0
ARG LDFLAGS
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gonelist -ldflags "${LDFLAGS}" main.go


FROM alpine:3.11.6
WORKDIR /root
ARG VERSION=v0.4
COPY --from=builder /root/myapp/gonelist /bin/gonelist
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add curl wget && \
    cd /etc && curl -sL https://github.com/Sillywa/gonelist-web/releases/download/${VERSION}/dist.tar.gz | tar -zxvf -
EXPOSE 8000
ENTRYPOINT ["/bin/gonelist"]