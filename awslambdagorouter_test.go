package awslambdagorouter

import (
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func test1(request RouterRequest) map[string]interface{} {
	return request.Body
}

func TestRouter(t *testing.T) {
	cases := []struct {
		testRequest      events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		//test1
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/",
				Body:       `{"test":"test"}`,
				HTTPMethod: "GET",
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: `{"test":"test"}`,
			},
		},
		{
			testRequest: events.APIGatewayProxyRequest{
				Path:       "/path",
				Body:       `{"test":"test"}`,
				HTTPMethod: "GET",
			},
			expectedResponse: response404,
		},
	}
	for index, c := range cases {
		router := start()
		router.get("/", test1)
		got, err := router.serve(c.testRequest)

		var logger strings.Builder
		logger.WriteString("Doing test number: ")
		logger.WriteString(strconv.Itoa(index))
		t.Logf(logger.String())

		if err != nil {
			t.Errorf("Test resulted in an error")
			continue
		}

		if got.Body != c.expectedResponse.Body {
			t.Errorf("Test not passed, the response's body is different")
		} else {
			t.Logf("PASSED")
		}
	}
}
