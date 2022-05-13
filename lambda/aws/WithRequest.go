package aws

import "github.com/aws/aws-lambda-go/events"

type WithRequest interface {
	handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
