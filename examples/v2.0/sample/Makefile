mkfile_dir := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

build:	
	go install github.com/googleapis/gnostic
	make -C $(mkfile_dir)/../../../

all:	build
	gnostic sample.yaml --go-generator-out=sample

clean:  
	rm -rf sample service/service

test:	all
	killall service; true # ignore errors due to no matching processes
	cd service; go get .; go build; ./service &
	go test
	killall service
