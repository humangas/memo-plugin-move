PLUGIN_DIR=$(shell grep "pluginsdir.*" ~/.config/memo/config.toml | grep -o "\".*\"" | sed -e 's/"//g')

.DEFAULT_GOAL := help

.PHONY: all help setup deps install

all:

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "target:"
	@echo " - install:    install memo-plugin-memo"
	@echo " - deps:       dep ensure"
	@echo ""

setup:
	go get -u github.com/golang/dep/cmd/dep

deps: setup
	dep ensure

install: deps
	GOOS=darwin go build -o move *.go
	mv move $(PLUGIN_DIR)
