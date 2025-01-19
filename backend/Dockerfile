FROM golang:1.21-alpine

WORKDIR /app

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 編譯
RUN go build -o main .

# 暴露端口
EXPOSE 8080

# 運行
CMD ["./main"] 