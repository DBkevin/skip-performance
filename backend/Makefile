.PHONY: test coverage build run clean

# 测试
test:
	go test -v ./tests/...

# 带覆盖率测试
coverage:
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 显示覆盖率百分比
coverage-report:
	go test -cover ./tests/...
	go test -coverprofile=coverage.out ./tests/
	@echo "---"
	@go tool cover -func=coverage.out | tail -1

# 构建
build:
	go build -o bin/server .

# 运行
run:
	go run .

# 清理
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# 安装依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 运行所有检查
check: fmt test coverage-report
	@echo "所有检查通过"
