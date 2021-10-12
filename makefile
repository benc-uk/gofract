# Common variables
VERSION := 0.0.1
BUILD_INFO := Manual build 

SRC_DIR := cmd
GO_PKG := github.com/benc-uk/gofract

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint # Remove if not using Go

.PHONY: help run lint lint-fix
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## 🔨 Run a local build without a container
	@mkdir -p bin
	go mod tidy
	GOOS=linux go build -o bin/fract ./$(SRC_DIR)/fract/...
	GOOS=windows go build -o bin/fract.exe ./$(SRC_DIR)/fract/...