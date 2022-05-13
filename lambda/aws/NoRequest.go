package aws

import (
	"github.com/aws/aws-lambda-go/events"
)

type NoRequest interface {
	handle() (events.APIGatewayProxyResponse, error)
}
