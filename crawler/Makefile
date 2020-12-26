CGO_ENABLED := 0
GOOS ?= darwin
GOARCH ?= amd64

PROJECT_DIR := $(abspath $(shell pwd)/..)

build:
	@go install .

run: clean crawl-main crawl-add crawl-signs

clean:
	@rm -rf ${PROJECT_DIR}/data

crawl-main:
	@crawler -dir ${PROJECT_DIR}/data \
		-ids 130,131,132,276,133,134,135,136,137,138,139,140,141,142,143,144,145,146,147,148,149,150,151

crawl-add:
	@crawler -dir ${PROJECT_DIR}/data \
		-ids 165,166,153,154,152,155

crawl-signs:
	@crawler -dir ${PROJECT_DIR}/data \
		-ids 156,157,158,159,160,161,162,163,164,165
