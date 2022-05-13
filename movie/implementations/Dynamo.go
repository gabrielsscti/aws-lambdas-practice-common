package implementations

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gabrielsscti/aws-lambdas-practice-common"
	"github.com/gabrielsscti/aws-lambdas-practice-common/movie"
	"log"
	"net/http"
	"os"
)

type Dynamo struct {
	cfg aws.Config
}

func (d Dynamo) Delete(id string) error {
	log.Println("Deleting ID " + id + " from DynamoDB")
	svc := dynamodb.NewFromConfig(d.cfg)

	dynamoKey, err := getDynamoKey(id)
	if err != nil {
		log.Println("Error while getting key")
		return err
	}

	_, err = svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       dynamoKey,
		TableName: aws.String(os.Getenv(movie.MOVIE_TABLE_ENV)),
	})
	if err != nil {
		return common.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while deleting movie from database",
			Info:       err.Error(),
		}
	}

	return nil
}

func CreateNewDynamoRepository(cfg aws.Config) Dynamo {
	return Dynamo{
		cfg,
	}
}

func (d Dynamo) GetByID(id string) ([]byte, error) {
	svc := dynamodb.NewFromConfig(d.cfg)

	dynamoKey, err := getDynamoKey(id)
	if err != nil {
		return nil, err
	}

	res, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       dynamoKey,
		TableName: aws.String(os.Getenv(movie.MOVIE_TABLE_ENV)),
	})
	if err != nil {
		return nil, common.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Body:       "Could not fetch from database",
			Info:       err.Error(),
		}
	}

	var movieObj movie.Movie
	err = attributevalue.UnmarshalMap(res.Item, &movieObj)
	if err != nil {
		return nil, common.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Body:       "Could not get correct value from database",
			Info:       err.Error(),
		}
	}

	response, err := json.Marshal(movieObj)
	if err != nil {
		return nil, common.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding values from database",
			Info:       err.Error(),
		}
	}

	return response, nil
}

func getDynamoKey(id string) (map[string]types.AttributeValue, error) {
	dynamoKey, err := attributevalue.MarshalMap(struct {
		ID *string
	}{ID: aws.String(id)})
	if err != nil {
		return nil, common.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Body:       "Could not marshal key",
			Info:       err.Error(),
		}
	}

	return dynamoKey, err
}

func (d Dynamo) GetAll() ([]byte, error) {
	svc := dynamodb.NewFromConfig(d.cfg)
	res, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{TableName: aws.String(os.Getenv(movie.MOVIE_TABLE_ENV))})
	if err != nil {
		return nil, common.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error fetching the database",
			Info:       err.Error(),
		}
	}

	movies := make([]movie.Movie, 0)
	for _, item := range res.Items {
		var movie movie.Movie
		err := attributevalue.UnmarshalMap(item, &movie)
		if err != nil {
			return nil, common.ResponseError{
				StatusCode: http.StatusInternalServerError,
				Body:       "Could not get correct value from database",
				Info:       err.Error(),
			}
		}
		movies = append(movies, movie)
	}

	response, err := json.Marshal(movies)
	if err != nil {
		return nil, common.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding values from database",
			Info:       err.Error(),
		}
	}

	return response, nil
}

func (d Dynamo) Insert(obj []byte) error {
	var movieObj movie.Movie
	err := json.Unmarshal(obj, &movieObj)
	if err != nil {
		return common.ResponseError{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid payload",
			Info:       err.Error(),
		}
	}

	svc := dynamodb.NewFromConfig(d.cfg)

	dynamoInput, err := attributevalue.MarshalMap(movieObj)
	if err != nil {
		return common.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Body:       "Could not convert to database structure",
			Info:       err.Error(),
		}
	}

	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      dynamoInput,
		TableName: aws.String(os.Getenv(movie.MOVIE_TABLE_ENV)),
	})
	if err != nil {
		return common.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Body:       "Could not insert in database",
			Info:       err.Error(),
		}
	}

	return nil
}
