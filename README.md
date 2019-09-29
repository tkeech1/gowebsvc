[![codecov](https://codecov.io/gh/tkeech1/gowebsvc/branch/master/graph/badge.svg)](https://codecov.io/gh/tkeech1/gowebsvc)
[![Go Report Card](https://goreportcard.com/badge/github.com/tkeech1/gowebsvc)](https://goreportcard.com/report/github.com/tkeech1/gowebsvc)
[![CircleCI](https://circleci.com/gh/tkeech1/gowebsvc.svg?style=svg)](https://circleci.com/gh/tkeech1/gowebsvc)
[![Build Status](https://dev.azure.com/tkeech1/gowebsvc/_apis/build/status/tkeech1.gowebsvc?branchName=master)](https://dev.azure.com/tkeech1/gowebsvc/_build/latest?definitionId=1&branchName=master)
[![Git Action Status](https://github.com/tkeech1/gowebsvc/workflows/Go/badge.svg)](https://github.com/tkeech1/gowebsvc/actions?workflow=Go)


A comparison between gokit and simple web services.

See the Makefile.  

### Go Kit Web Service
```
make run-gokit
```

### Simple Web Service

The simple web service runs a HTTP service on 8080 and a GRPC service on 50051. 

```
make run-simple
```

To use the GRPC client:

``` 
make run-grpc-client
```
