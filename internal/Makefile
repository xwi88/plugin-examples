# Makefile to build the command lines and tests in Seele project.
# This Makefile doesn't consider Windows Environment. If you use it in Windows, please be careful.

SHELL := /bin/bash

#BASEDIR = $(shell pwd)
BASEDIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))

$(info Current DIR: '${BASEDIR}')

.PHONY: default clean abi grpc go abi-go grpc-go

default: go

all: go

clean:
	rm -rf pb/go/*

abi: abi-go

grpc: grpc-go

go: abi-go grpc-go

abi-go:
	mkdir -p pb/go && \
	protoc -I proto/ --go_out=pb/go/ --go_opt=paths=source_relative proto/*_abi.proto
	@echo "[abi-go] done"

grpc-go:
	mkdir -p pb/go && \
	protoc -I proto hello_rpc.proto --go_out=plugins=grpc:pb/go
	@echo "[grpc-go] done"
