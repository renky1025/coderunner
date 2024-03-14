# Copyright (C) 2023-present, expertsystem.digitsvalue.com

TARGET := $(shell uname -s)
GO     := GO111MODULE=on go

ifeq ("$(TARGET)", "Darwin")
	ARCH := darwin
endif

ifeq ("$(TARGET)", "Linux")
	ARCH := linux
endif

ifeq ("$(os)", "darwin")
	TARGET := Darwin
	ARCH   := darwin
endif

ifeq ("$(os)", "linux")
	TARGET := Linux
	ARCH   := linux
endif

ifeq ("$(os)", "windows")
	TARGET := Windows
	ARCH   := windows
	EXT    := .exe
endif

GOBUILD = CGO_ENABLED=0 GOOS=$(ARCH) GOARCH=amd64 $(GO) build -ldflags "-s -w"
BIN     = ./bin
CMD     = ./
.PHONY: info all clean coderunner

default: all

all: info clean  coderunner

info:
	@echo ---Building coderunner for $(TARGET)...

coderunner: info
	@$(GOBUILD) -o $(BIN)/$@$(EXT) $(CMD)/main.go
	@echo "Build $@ successfully!"

clean:
	@rm -f $(BIN)/*
	@echo "Clean successfully!"
