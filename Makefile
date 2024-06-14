BIN_OUT_DIR := out
BINARY := $(BIN_OUT_DIR)/terraform-provider-cala

version = 0.0.7
os_arch = $(shell go env GOOS)_$(shell go env GOARCH)
provider_path = registry.terraform.io/galoymoney/cala/$(version)/$(os_arch)/

install: build
	mkdir -p ~/.terraform.d/plugins/${provider_path}
	mv ${BINARY} ~/.terraform.d/plugins/${provider_path}
	rm -rf example/.terraform example/.terraform.lock.hcl example/terraform.tfstate*

build:
	go build -o $(BINARY) main.go

generate:
	go run github.com/Khan/genqlient

gen-docs:
	 go generate ./...
