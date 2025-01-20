package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

// func (m *MockTableClient) AddUrl(ctx context.Context, url models.Url) error {
// 	args := m.Called(ctx, url)
// 	return args.Error(0)
// }

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
