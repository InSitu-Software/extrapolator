BUILD_PREFIX?=./build
GO:=$(shell which go)
ADDITIONAL_BUILD_FLAGS=-a

os=linux freebsd darwin

all: clean $(os)

$(os):
	@mkdir -p $(BUILD_PREFIX)/$@
	env GO111MODULE=on GOOS=$@ GOARH=amd64 $(GO) build $(ADDITIONAL_BUILD_FLAGS) -o $(BUILD_PREFIX)/$@/extrapolator bin/extrapolator.go

clean:
	@rm -rf $(BUILD_PREFIX)
