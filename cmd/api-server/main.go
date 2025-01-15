package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sony/sonyflake"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func main() {
	// id, _ := sf.NextID()
	// byteSlice := make([]byte, 8)
	// binary.BigEndian.PutUint64(byteSlice, id)
	// encoded := base62.EncodeToString(byteSlice)
	// fmt.Println(encoded)
	// r := gin.Default()
	// docs.SwaggerInfo.BasePath = "/api/v1"
	// v1 := r.Group("/api/v1")
	// {
	// 	eg := v1.Group("/example")
	// 	{
	// 		eg.GET("/helloworld", Helloworld)
	// 	}
	// }
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// r.Run(":8080")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("admin"),
		config.WithRegion("ap-southeast-1"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Build the request with its input parameters
	resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}

	fmt.Println("Tables:")
	for _, tableName := range resp.TableNames {
		fmt.Println(tableName)
	}
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}
