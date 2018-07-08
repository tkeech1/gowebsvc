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

curl-greeting:
	#curl -d "{\"s\":\"hello\"}" -X POST http://localhost:8080/greeting
	curl -d "{\"s\":\"\"}" -X POST http://localhost:8080/greeting
	