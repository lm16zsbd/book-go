.PHONY: run dev build clean test vet lint

SWAG := $(shell go env GOPATH)/bin/swag
AIR  := $(shell go env GOPATH)/bin/air

# 直接运行 (自动生成 swagger)
run: swagger
	go run ./cmd/server

# 热重载开发 (先安装: go install github.com/air-verse/air@latest)
dev: swagger
	$(AIR)

# 编译
build: swagger
	go build -o bin/server ./cmd/server

# 清理
clean:
	rm -rf bin/

# 测试
test:
	go test ./... -v

# 代码检查
vet:
	go vet ./...

lint:
	go vet ./...

# 依赖整理
tidy:
	go mod tidy

# 生成 swagger 文档
swagger:
	@if [ -f "$(SWAG)" ]; then $(SWAG) init -g cmd/server/main.go -o docs; else echo "swag not found, skipping"; fi
