package main

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {

		router := gin.Default()
		router.GET("/page", nil)

		router.GET("/api/pet/:id", func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"id":   c.Param("id"),
					"name": randomName(),
					"time": randomDate(),
				},
			)
		})

		// router.POST("/api/pet", getPet)

		ginLambda = ginadapter.New(router)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}

func randomDate() time.Time {
	now := time.Now()
	start := now.AddDate(-15, 0, 0)
	delta := now.Unix() - start.Unix()

	sec := rand.Int63n(delta) + start.Unix()
	return time.Unix(sec, 0)
}

func randomName() string {
	names := []string{"Bailey", "Bella", "Max", "Lucy", "Charlie", "Molly", "Buddy", "Daisy", "Rocky", "Maggie", "Jake", "Sophie", "Jack", "Sadie", "Toby", "Chloe", "Cody", "Bailey", "Buster", "Lola", "Duke", "Zoe", "Cooper", "Abby", "Riley", "Ginger", "Harley", "Roxy", "Bear", "Gracie", "Tucker", "Coco", "Murphy", "Sasha", "Lucky", "Lily", "Oliver", "Angel", "Sam", "Princess", "Oscar", "Emma", "Teddy", "Annie", "Winston", "Rosie"}
	return names[random(0, len(names))]
}

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}
