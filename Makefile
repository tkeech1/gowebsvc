clean-testcache:
	go clean -testcache github.com/tkeech1/gowebsvc/gokit/
	go clean -testcache github.com/tkeech1/gowebsvc/simple/

test-gokit: clean-testcache	
	go test -v -race -covermode=atomic github.com/tkeech1/gowebsvc/gokit/

test-simple: clean-testcache	
	go test -v -race -covermode=atomic github.com/tkeech1/gowebsvc/simple/

test: clean-testcache	
	go test -race -covermode=atomic ./...

test-circleci: 
	go test -race -covermode=atomic -coverprofile=coverage.txt ./...

deps: 
	go get -v -t -d ./...

run-gokit:	
	cd gokit/; go build; ./gokit

run-simple:	
	cd simple/; go build; ./simple

curl-greeting:
	#curl -d "{\"s\":\"hello\"}" -X POST http://localhost:8080/greeting
	curl -d "{\"s\":\"\"}" -X POST http://localhost:8080/greeting

compile-grpc:
	cd svc/; protoc greeting.proto --go_out=plugins=grpc:.
	
run-grpc-client:
	cd client/; go build; ./client safsdfadfs