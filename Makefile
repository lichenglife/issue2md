.PHONY: build test clean lint install web deps e2e-test

# 构建 CLI
build:
	go build -o bin/issue2md ./cmd/issue2md

# 运行所有测试
test:
	go test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 清理构建产物
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# 运行 lint
lint:
	go vet ./...

# 安装到 GOPATH/bin
install:
	go install ./cmd/issue2md

# 构建 Web 服务
web:
	go build -o bin/issue2mdweb ./cmd/issue2mdweb

# 下载依赖
deps:
	go mod tidy
	go mod download

# 端到端测试
e2e-test: build
	./scripts/e2e-test.sh
