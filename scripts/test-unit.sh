#!/bin/bash
# 运行单元测试，排除集成和端到端测试目录

set -e

echo "Running unit tests..."
go test -v -cover -race $(go list ./... | grep -v /tests/e2e | grep -v /internal/interfaces/http/handler)