#!/usr/bin/make -f

# Copyright (c) 2020-present, The kubequery authors
#
# This source code is licensed as defined by the LICENSE file found in the
# root directory of this source tree.
#
# SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)

all: deps test build kubequery.yaml

deps:
	@go mod download

build: deps
	@go build -ldflags="-s -w" -o . ./...

test:
	@go test -race -cover ./...

docker: build
	@docker build --build-arg KUBEQUERY_VERSION=latest -t uptycs/kubequery:latest .

genschema: build
	@echo "\`\`\`sql" >  docs/schema.md
	@./genschema      >> docs/schema.md
	@echo "\`\`\`"    >> docs/schema.md

clean:
	@rm -f kubequery genschema

.PHONY: all
