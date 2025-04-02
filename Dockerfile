FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum (如果存在)
COPY go.mod ./
RUN go mod tidy

# 复制源代码
COPY . .

# 编译为静态二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o emailservice .

# 使用distroless作为最终基础镜像，减小体积
FROM alpine:3.19

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/emailservice .

# 端口
EXPOSE 8080

# 运行
CMD ["./emailservice"] 