mkfile_dir := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

build:	
	go install github.com/googleapis/gnostic
	make -C $(mkfile_dir)/../../../

all:	build
	gnostic swagger.json --go-generator-out=xkcd

clean:  
	rm -rf xkcd

test:	all
	go run main.go
