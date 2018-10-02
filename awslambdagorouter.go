package awslambdagorouter

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// RouterRequest ...
type RouterRequest struct {
	Resource              string            `json:"resource"` // The resource path defined in API Gateway
	Path                  string            `json:"path"`     // The url path for the caller
	HTTPMethod            string            `json:"httpMethod"`
	Headers               map[string]string `json:"headers"`
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	PathParameters        map[string]string `json:"pathParameters"`
	Body                  RouterResponse    `json:"body"`
	IsBase64Encoded       bool              `json:"isBase64Encoded,omitempty"`
}

// RouterResponse ...
type RouterResponse = map[string]interface{}

func convertBodyToJSON(request events.APIGatewayProxyRequest) RouterRequest {

	body := make(RouterResponse)
	if request.Body != "" {
		if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
			panic("Unable to decode JSON")
		}
	}

	return RouterRequest{
		Resource:              request.Resource,
		Path:                  request.Path,
		HTTPMethod:            request.HTTPMethod,
		Headers:               request.Headers,
		QueryStringParameters: request.QueryStringParameters,
		PathParameters:        request.PathParameters,
		Body:                  body,
		IsBase64Encoded:       request.IsBase64Encoded,
	}
}

func createResponse(jsonData RouterResponse) events.APIGatewayProxyResponse {

	body, _ := json.Marshal(jsonData)
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{},
		Body:            string(body),
		IsBase64Encoded: false,
	}
}

// CallbackFunction ...
type CallbackFunction func(RouterRequest) RouterResponse

// RouterFunctions ...
type RouterFunctions struct {
	get  map[string]CallbackFunction
	post map[string]CallbackFunction
}

// Router ...
type Router struct {
	functions RouterFunctions
}

var response404 = events.APIGatewayProxyResponse{
	StatusCode: 404,
	Body:       "Route not found for this method",
}

func response400(response string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       response,
	}
}

// Start creates a Router
func Start() Router {
	routerFunctions := RouterFunctions{
		get:  map[string]CallbackFunction{},
		post: map[string]CallbackFunction{},
	}
	router := Router{functions: routerFunctions}
	return router
}

// Get ...
func (r Router) Get(path string, callback CallbackFunction) {
	r.functions.get[path] = callback
}

// Post ...
func (r Router) Post(path string, callback CallbackFunction) {
	r.functions.post[path] = callback
}

//Serve ...
func (r Router) Serve(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {

	defer func() {
		if r := recover(); r != nil {
			if r == "Unable to decode JSON" {
				response = response400(r.(string))
			}
		}
	}()

	path := request.Path
	method := strings.ToLower(request.HTTPMethod)

	var callbackFunction CallbackFunction

	if method == "get" {
		callbackFunction = r.functions.get[path]
	} else if method == "post" {
		callbackFunction = r.functions.post[path]
	}

	if callbackFunction != nil {
		convertedRequest := convertBodyToJSON(request)
		data := callbackFunction(convertedRequest)
		response = createResponse(data)
	} else {
		response = response404
	}

	println("responding")
	return response, err

}
