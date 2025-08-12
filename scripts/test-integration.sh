#!/bin/bash
# 运行集成测试
# 根据约定，这些测试位于 internal/interfaces/http/handler/ 目录下

set -e

echo "Running integration tests..."
go test -v -cover -race ./internal/interfaces/http/handler