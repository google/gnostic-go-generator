
build:	
	go install github.com/googleapis/gnostic-go-generator
	rm -f $(GOPATH)/bin/gnostic-go-client $(GOPATH)/bin/gnostic-go-server
	ln -s $(GOPATH)/bin/gnostic-go-generator $(GOPATH)/bin/gnostic-go-client
	ln -s $(GOPATH)/bin/gnostic-go-generator $(GOPATH)/bin/gnostic-go-server

test:
	pushd examples/v2.0/sample && make test && popd
	pushd examples/v2.0/bookstore && make test && popd
	pushd examples/v3.0/bookstore && make test && popd