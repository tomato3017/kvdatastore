# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOBUILDOUT=bin
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=tomatobot
IMAGE_NAME=$(BINARY_NAME)
IMAGE_TAG=$(shell git describe --tags --always)
BUILD=$(shell git describe --always --long)
VERSION=$(IMAGE_TAG)

all: test build
build: vendor
		$(GOBUILD) .
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
vendor:
		go mod vendor
lint:
		golint .
redisup:
		docker run -d --name redissrv --network host  --rm redis
redisdown:
		docker stop redissrv