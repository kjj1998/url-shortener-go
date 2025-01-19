package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/kjj1998/url-shortener-go/constants"
	"github.com/kjj1998/url-shortener-go/internal/models"
)

func enterTest() (context.Context, *testtools.AwsmStubber, *TableClient) {
	ctx := context.Background()
	stubber := testtools.NewStubber()
	client := &TableClient{TableName: constants.TableName, DynamoDbClient: dynamodb.NewFromConfig(*stubber.SdkConfig)}
	return ctx, stubber, client
}

func TestTableClient_AddUrl(t *testing.T) {
	t.Run("NoErrors", func(t *testing.T) { AddUrl(nil, t) })
	t.Run("TestError", func(t *testing.T) { AddUrl(&testtools.StubError{Err: errors.New("TestError")}, t) })
}

func AddUrl(raiseErr *testtools.StubError, t *testing.T) {
	ctx, stubber, client := enterTest()

	url := models.Url{Id: 12345, LongUrl: "https://www.youtube.com", ShortUrl: "NEDF34qw"}
	item, marshErr := attributevalue.MarshalMap(url)
	if marshErr != nil {
		panic(marshErr)
	}

	stubber.Add(StubAddUrl(client.TableName, item, raiseErr))

	err := client.AddUrl(ctx, url)

	testtools.VerifyError(err, raiseErr, t)
	testtools.ExitTest(stubber, t)
}

func StubAddUrl(tableName string, item map[string]types.AttributeValue, raiseErr *testtools.StubError) testtools.Stub {
	return testtools.Stub{
		OperationName: "PutItem",
		Input:         &dynamodb.PutItemInput{TableName: aws.String(tableName), Item: item},
		Output:        &dynamodb.PutItemOutput{},
		Error:         raiseErr,
	}
}

func TestTableClient_RetrieveUrl(t *testing.T) {
	t.Run("NoErrors", func(t *testing.T) { RetrieveUrl(nil, t) })
	t.Run("TestError", func(t *testing.T) { RetrieveUrl(&testtools.StubError{Err: errors.New("TestError")}, t) })
	t.Run("NoUrlsRetrieved", func(t *testing.T) { RetrieveNoUrl(nil, t) })
}

func RetrieveUrl(raiseErr *testtools.StubError, t *testing.T) {
	ctx, stubber, client := enterTest()

	shortUrl := "NEWDSa31"
	longUrl := "https://www.youtube.com"
	var Id uint64 = 5438989247290

	stubber.Add(StubRetrieveUrl(client.TableName, shortUrl, longUrl, Id, raiseErr))

	urls, err := client.RetrieveUrl(ctx, shortUrl)

	testtools.VerifyError(err, raiseErr, t)
	if err == nil {
		if len(urls) == 0 {
			t.Errorf("Expected at least 1 url")
		}
	}

	testtools.ExitTest(stubber, t)
}

func RetrieveNoUrl(raiseErr *testtools.StubError, t *testing.T) {
	ctx, stubber, client := enterTest()

	shortUrl := "NEWDSa31"
	longUrl := "https://www.youtube.com"
	var Id uint64 = 5438989247290

	stubber.Add(StubRetrieveNoUrl(client.TableName, shortUrl, longUrl, Id, raiseErr))

	urls, err := client.RetrieveUrl(ctx, shortUrl)

	testtools.VerifyError(err, raiseErr, t)

	if err == nil {
		if len(urls) == 0 {
			fmt.Println("No URLs retrieved")
		}
	}

	testtools.ExitTest(stubber, t)
}

func StubRetrieveUrl(tableName string, shortUrl string, longUrl string, id uint64, raiseErr *testtools.StubError) testtools.Stub {
	keyEx := expression.Key("ShortUrl").Equal(expression.Value(shortUrl))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	return testtools.Stub{
		OperationName: "Query",
		Input: &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
			IndexName:                 aws.String(constants.TableSecondaryIndex),
		},
		Output: &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{{
			"Id":       &types.AttributeValueMemberN{Value: strconv.FormatUint(id, 10)},
			"ShortUrl": &types.AttributeValueMemberS{Value: shortUrl},
			"LongUrl":  &types.AttributeValueMemberS{Value: longUrl},
		}}},
		Error: raiseErr,
	}
}

func StubRetrieveNoUrl(tableName string, shortUrl string, longUrl string, id uint64, raiseErr *testtools.StubError) testtools.Stub {
	keyEx := expression.Key("ShortUrl").Equal(expression.Value(shortUrl))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	return testtools.Stub{
		OperationName: "Query",
		Input: &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
			IndexName:                 aws.String(constants.TableSecondaryIndex),
		},
		Output: &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}},
		Error:  raiseErr,
	}
}
