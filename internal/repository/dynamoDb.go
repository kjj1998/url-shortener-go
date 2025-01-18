package repository

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kjj1998/url-shortener-go/constants"
	"github.com/kjj1998/url-shortener-go/internal/models"
)

type TableClient struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

var Client TableClient

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(constants.AwsProfile),
		config.WithRegion(constants.AwsRegion),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	Client = TableClient{
		DynamoDbClient: dynamodb.NewFromConfig(cfg),
		TableName:      constants.TableName,
	}
}

// ListTables returns an array of Go strings containing the names of the DynamoDB tables
// Returns an error the names of the tables cannot be retrieved
func (client TableClient) ListTables(ctx context.Context) ([]string, error) {
	var tableNames []string
	var output *dynamodb.ListTablesOutput
	var err error

	tablePaginator := dynamodb.NewListTablesPaginator(client.DynamoDbClient, &dynamodb.ListTablesInput{})
	for tablePaginator.HasMorePages() {
		output, err = tablePaginator.NextPage(ctx)
		if err != nil {
			log.Printf("Couldn't list tables. Here's why: %v\n", err)
			break
		} else {
			tableNames = append(tableNames, output.TableNames...)
		}
	}

	return tableNames, err
}

func (client TableClient) AddUrl(ctx context.Context, url models.Url) error {
	item, err := attributevalue.MarshalMap(url)
	if err != nil {
		panic(err)
	}

	_, err = client.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(client.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}
