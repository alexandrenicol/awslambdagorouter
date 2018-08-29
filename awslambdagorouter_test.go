package awslambdagorouter

import (
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func test1(request RouterRequest) map[string]interface{} {
	return map[string]interface{}{
		"success": "true",
	}
}
func test2(request RouterRequest) map[string]interface{} {
	return map[string]interface{}{
		"success": "true",
	}
}
func test3(request RouterRequest) map[string]interface{} {
	return request.Body
}

func test4(request RouterRequest) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]string{
			"category": request.QueryStringParameters["category"],
			"campaign": request.QueryStringParameters["campaign"],
		},
		"success": "true",
	}
}

func TestRouter(t *testing.T) {
	cases := []struct {
		testRequest      events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		//test that get functions are registered and can be executed properly
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/",
				Body:       "",
				HTTPMethod: "GET",
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: `{"success":"true"}`,
			},
		},
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/test",
				Body:       "",
				HTTPMethod: "GET",
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: `{"success":"true"}`,
			},
		},
		//test that it returns a 404 if the path does not exist
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/path",
				Body:       "",
				HTTPMethod: "GET",
			},
			expectedResponse: response404,
		},
		//test that it returns a 404 if the method does no exist, even for an existing path
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/",
				Body:       "",
				HTTPMethod: "DELETE",
			},
			expectedResponse: response404,
		},
		//test POST request routing
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/",
				Body:       `{"success":"true", "data":[0,1,2]}`,
				HTTPMethod: "POST",
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: `{"success":"true"}`,
			},
		},
		//test POST request parameters transit
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/test",
				Body:       `{"success":"true", "data":[0,1,2]}`,
				HTTPMethod: "POST",
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: `{"data":[0,1,2],"success":"true"}`,
			},
		},
		//test GET request query string parameters transit
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/query",
				HTTPMethod: "GET",
				QueryStringParameters: map[string]string{
					"category": "one",
					"campaign": "test",
				},
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: `{"data":{"campaign":"test","category":"one"},"success":"true"}`,
			},
		},
	}
	for index, c := range cases {
		router := start()
		router.get("/", test1)
		router.get("/test", test1)
		router.post("/", test2)
		router.post("/test", test3)
		router.get("/query", test4)
		got, err := router.serve(c.testRequest)

		var logger strings.Builder
		logger.WriteString("Doing test number: ")
		logger.WriteString(strconv.Itoa(index))
		t.Logf(logger.String())

		if err != nil {
			t.Errorf("Test resulted in an error")
			continue
		}

		println(got.Body)
		println(c.expectedResponse.Body)

		if got.Body != c.expectedResponse.Body {
			t.Errorf("Test not passed, the response's body is different")
		} else {
			t.Logf("PASSED")
		}
	}
}
