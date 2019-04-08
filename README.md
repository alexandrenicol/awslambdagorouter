# awslambdagorouter

[![Build Status](https://travis-ci.org/alexandrenicol/awslambdagorouter.png?branch=master)](https://travis-ci.org/alexandrenicol/awslambdagorouter)
[![codecov](https://codecov.io/gh/alexandrenicol/awslambdagorouter/branch/master/graph/badge.svg)](https://codecov.io/gh/alexandrenicol/awslambdagorouter)

## examples
```golang
package main

import (
	"github.com/alexandrenicol/awslambdagorouter"
	"github.com/aws/aws-lambda-go/lambda"
)

func test1(request awslambdagorouter.RouterRequest) awslambdagorouter.RouterResponse {
	return map[string]interface{}{
		"success": "hello world",
	}
}

func main() {
	router := awslambdagorouter.Start()
	router.Get("/helloworld", test1)

	lambda.Start(router.Serve)
}
```