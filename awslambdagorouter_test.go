package awslambdagorouter

import (
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func test1(request RouterRequest) RouterResponse {
	return map[string]interface{}{
		"success": "true",
	}
}
func test3(request RouterRequest) RouterResponse {
	return request.Body
}

func test4(request RouterRequest) RouterResponse {
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
		// 0, 1, test that get functions are registered and can be executed properly
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
		//2, test that it returns a 404 if the path does not exist
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/path",
				Body:       "",
				HTTPMethod: "GET",
			},
			expectedResponse: response404,
		},
		//3, test that it returns a 404 if the method does no exist, even for an existing path
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/",
				Body:       "",
				HTTPMethod: "DELETE",
			},
			expectedResponse: response404,
		},
		//4, test POST request routing
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
		//5, test POST request parameters transit
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
		//6, test POST request, with no body
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/noparam",
				HTTPMethod: "POST",
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: "{}",
			},
		},
		//7, test POST request, with empty body
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/noparam2",
				Body:       "",
				HTTPMethod: "POST",
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: "{}",
			},
		},
		//8, test GET request query string parameters transit
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
		//9 test POST request, with empty body
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/unvalid",
				Body:       `{"unfinishedJson`,
				HTTPMethod: "POST",
			},
			expectedResponse: response400("Unable to decode JSON"),
		},
	}
	for index, c := range cases {
		router := Start()
		router.Get("/", test1)
		router.Get("/test", test1)
		router.Post("/", test1)
		router.Post("/test", test3)
		router.Post("/noparam", test3)
		router.Post("/noparam2", test3)
		router.Post("/unvalid", test3)
		router.Get("/query", test4)
		got, err := router.Serve(c.testRequest)

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
