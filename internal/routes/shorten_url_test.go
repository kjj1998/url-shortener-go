package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/kjj1998/url-shortener-go/constants"
	"github.com/kjj1998/url-shortener-go/internal/models"
	"github.com/kjj1998/url-shortener-go/internal/repository"
	"github.com/kjj1998/url-shortener-go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTableClient struct {
	mock.Mock
}

func (m *MockTableClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockTableClient) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func TestGenerateShortenedUrl(t *testing.T) {
	mockRepo := new(MockTableClient)
	repository.Client = repository.TableClient{
		DynamoDbClient: mockRepo,
		TableName:      constants.TableName,
	}

	utils.GenerateUniqueId = func() uint64 { return 2387497 }
	utils.ShortenUrl = func(id uint64) string { return "NWER425d" }

	tests := []struct {
		name           string
		payload        interface{}
		mockRepoError  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid request",
			payload:        models.LongUrl{LongUrl: "http://example.com"},
			mockRepoError:  nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":2387497,"shortUrl":"NWER425d","longUrl":"http://example.com"}`,
		},
		{
			name:           "Failed JSON binding",
			payload:        "124",
			mockRepoError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400, "message":"json: cannot unmarshal number into Go value of type models.LongUrl"}`,
		},
		{
			name:           "Validation error",
			payload:        map[string]string{"shortUrl": "fdsfsdfdsf"},
			mockRepoError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400, "message":"invalid parameter names in json body"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockRepoError == nil {
				mockRepo.On("AddUrl", mock.Anything, mock.Anything).Return(nil).Once()
				mockRepo.On("PutItem", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil).Once()
			} else {
				mockRepo.On("AddUrl", mock.Anything, mock.Anything).Return(tt.mockRepoError).Once()
			}
		})

		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)

		var requestBody []byte
		if strPayload, ok := tt.payload.(string); ok {
			requestBody = []byte(strPayload)
		} else {
			requestBody, _ = json.Marshal(tt.payload)
		}
		req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		GenerateShortenedUrl(ctx)

		assert.Equal(t, tt.expectedStatus, rec.Code)
		assert.JSONEq(t, tt.expectedBody, rec.Body.String())
	}
}

func TestRedirectShortenedUrl(t *testing.T) {
	tests := []struct {
		name             string
		param            string
		mockRepoError    error
		expectedStatus   int
		expectedLocation string
		mockRepoResult   string
		mockRepoId       uint64
	}{
		{
			name:             "Valid request",
			param:            "NWER425d",
			mockRepoError:    nil,
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "http://example.com",
			mockRepoResult:   "http://example.com",
		},
		{
			name:             "No parameter",
			param:            "",
			mockRepoError:    nil,
			expectedStatus:   http.StatusBadRequest,
			expectedLocation: "",
			mockRepoResult:   "nil",
		},
		{
			name:             "Internal error",
			param:            "KWBG425d",
			mockRepoError:    errors.New("simulated DynamoDB error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedLocation: "",
			mockRepoResult:   "nil",
		},
		{
			name:             "URL not found",
			param:            "KXBG425d",
			mockRepoError:    nil,
			expectedStatus:   http.StatusNotFound,
			expectedLocation: "",
			mockRepoResult:   "nil",
		},
	}

	for _, tt := range tests {
		mockRepo := new(MockTableClient)
		repository.Client = repository.TableClient{
			DynamoDbClient: mockRepo,
			TableName:      constants.TableName,
		}

		if tt.mockRepoError != nil {
			fmt.Println("test else")
			mockRepo.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{}, tt.mockRepoError).Once()
			mockRepo.On("RetrieveUrl", mock.Anything, tt.param).Return("", tt.mockRepoError).Once()
		} else {
			if tt.name == "Valid request" {
				mockDynamoResponse := &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"ShortUrl": &types.AttributeValueMemberS{Value: tt.param},
							"LongUrl":  &types.AttributeValueMemberS{Value: tt.mockRepoResult},
						},
					},
				}

				mockRepo.On("Query", mock.Anything, mock.Anything).Return(mockDynamoResponse, tt.mockRepoError).Once()
			} else if tt.name == "URL not found" {
				mockDynamoResponse := &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"ShortUrl": &types.AttributeValueMemberS{Value: tt.param},
							"LongUrl":  &types.AttributeValueMemberS{Value: ""},
						},
					},
				}

				mockRepo.On("Query", mock.Anything, mock.Anything).Return(mockDynamoResponse, tt.mockRepoError).Once()
			} else {
				mockRepo.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{}, tt.mockRepoError).Once()
			}

			mockRepo.On("RetrieveUrl", mock.Anything, tt.param).Return(tt.mockRepoResult, tt.mockRepoError).Once()
		}

		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/NWER425d", nil)
		ctx.Request = req

		ctx.Params = gin.Params{
			{Key: "shortUrl", Value: tt.param},
		}

		RedirectShortenedUrl(ctx)

		assert.Equal(t, tt.expectedStatus, rec.Code)
		assert.Equal(t, tt.expectedLocation, rec.Header().Get("Location"))
	}
}
