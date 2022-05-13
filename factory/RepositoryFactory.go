package factory

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gabrielsscti/aws-lambdas-practice-common/movie"
	"github.com/gabrielsscti/aws-lambdas-practice-common/movie/implementations"
	"log"
	"os"
)

const DYNAMO_DATASOURCE = "DYNAMO"

func getAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opt *config.LoadOptions) error {
		opt.Region = "sa-east-1"
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func CreateNewRepository() movie.Repository {
	switch os.Getenv("DATASOURCE") {
	case DYNAMO_DATASOURCE:
		return implementations.CreateNewDynamoRepository(getAWSConfig())
	default:
		return implementations.CreateNewMockMovies()
	}
}
