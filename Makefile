GOPATH:=$(shell go env GOPATH)
OBJECTDIR:=$(shell dirname $(shell dirname $(shell pwd)))
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
# IMAGE_TAG=$(GIT_TAG)-$(GIT_COMMIT)
IMAGE_TAG = 0.0.1
NAME = sms-srv
IMAGE_NAME = micro-welfare/${NAME}
NETWORK = micro-welfare
CFG_CLUSTER = prod
PROTO = sms

all: build
proto:
	@echo execute ${PROTO} proto file generate ${OBJECTDIR}
	protoc --proto_path=. --proto_path=../../ --micro_out=. --go_out=. proto/${PROTO}/${PROTO}.proto
	mv -if github.com/zzsds/micro-sms-service/sms/proto/${PROTO}/* proto/${PROTO}
	rm -rf github.com/
vendor:
	go mod vendor

build:
	go build -o ${NAME} *.go

test:
	go test -v ./... -cover

docker:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_NAME):latest

run:
	docker run --rm --name ${NAME} -p :50051 --network ${NETWORK} -e MICRO_ADDRESS=:50051 -e MICRO_REGISTRY=mdns -e MYSQL_HOST=${MYSQL_HOST} ${IMAGE_NAME}:${IMAGE_TAG}

.PHONY: build proto clean vet test docker run