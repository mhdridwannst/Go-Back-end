FROM golang:1.23.4-alpine
WORKDIR /app

# 1. Install Git karena Alpine Linux polosan tidak punya Git
RUN apk update && apk add --no-cache git

# 2. Install Air versi stabil
RUN go install github.com/cosmtrek/air@v1.49.0

# 3. Jalankan Air
CMD ["air"]