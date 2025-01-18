package repository

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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
// Returns an error when the names of the tables cannot be retrieved
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

// AddUrl adds a URL and its shortened form as an entry into the DynamoDB table
// Returns an error when the URL could not be inserted into the table
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

func (client TableClient) RetrieveUrl(ctx context.Context, shortUrl string) (string, error) {
	var err error
	var response *dynamodb.QueryOutput
	var urls []models.Url

	keyEx := expression.Key("ShortUrl").Equal(expression.Value(shortUrl))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)
	} else {
		queryPaginator := dynamodb.NewQueryPaginator(client.DynamoDbClient, &dynamodb.QueryInput{
			TableName:                 aws.String(client.TableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
			IndexName:                 aws.String(constants.TableSecondaryIndex),
		})
		for queryPaginator.HasMorePages() {
			response, err = queryPaginator.NextPage(ctx)
			if err != nil {
				log.Printf("Couldn't query for urls with shortened url %v. Here's why: %v\n", shortUrl, err)
				break
			} else {
				var urlPage []models.Url
				err = attributevalue.UnmarshalListOfMaps(response.Items, &urlPage)
				if err != nil {
					log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
					break
				} else {
					urls = append(urls, urlPage...)
				}
			}
		}
	}

	if len(urls) > 0 {
		return urls[0].LongUrl, err
	} else {
		return "", err
	}
}

// // Query gets all movies in the DynamoDB table that were released in the specified year.
// // The function uses the `expression` package to build the key condition expression
// // that is used in the query.
// func (basics TableBasics) Query(ctx context.Context, releaseYear int) ([]Movie, error) {
// 	var err error
// 	var response *dynamodb.QueryOutput
// 	var movies []Movie
// 	keyEx := expression.Key("year").Equal(expression.Value(releaseYear))
// 	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
// 	if err != nil {
// 		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)
// 	} else {
// 		queryPaginator := dynamodb.NewQueryPaginator(basics.DynamoDbClient, &dynamodb.QueryInput{
// 			TableName:                 aws.String(basics.TableName),
// 			ExpressionAttributeNames:  expr.Names(),
// 			ExpressionAttributeValues: expr.Values(),
// 			KeyConditionExpression:    expr.KeyCondition(),
// 		})
// 		for queryPaginator.HasMorePages() {
// 			response, err = queryPaginator.NextPage(ctx)
// 			if err != nil {
// 				log.Printf("Couldn't query for movies released in %v. Here's why: %v\n", releaseYear, err)
// 				break
// 			} else {
// 				var moviePage []Movie
// 				err = attributevalue.UnmarshalListOfMaps(response.Items, &moviePage)
// 				if err != nil {
// 					log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
// 					break
// 				} else {
// 					movies = append(movies, moviePage...)
// 				}
// 			}
// 		}
// 	}
// 	return movies, err
// }
