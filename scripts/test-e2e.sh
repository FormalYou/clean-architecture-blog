#!/bin/bash
# 运行端到端（E2E）测试
# 这些测试位于 tests/e2e/ 目录下

set -e

echo "Running E2E tests..."
go test -v -cover -race ./tests/e2e