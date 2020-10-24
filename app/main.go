package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

var ginLambda *ginadapter.GinLambda

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Gin cold start")
		r := gin.Default()
		// r.GET("/pets", getPets)
		r.GET("/pets/:id", getPet)
		// r.POST("/pets", createPet)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}

func getPet(c *gin.Context) {
	petID := c.Param("id")
	randomPet := getRandomPet()
	randomPet.ID = petID
	c.JSON(200, randomPet)
}

func getRandomPet() Pet {
	pet := Pet{}

	pet.ID = getUUID()
	pet.Breed = randomBreed()
	pet.Name = randomName()

	pet.DateOfBirth = randomDate()

	return pet
}

// Pet
type Pet struct {
	ID          string    `json:"id"`
	Breed       string    `json:"breed"`
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}

var breeds = []string{"Afghan Hound", "Beagle", "Bernese Mountain Dog", "Bloodhound", "Dalmatian", "Jack Russell Terrier", "Norwegian Elkhound"}
var names = []string{"Bailey", "Bella", "Max", "Lucy", "Charlie", "Molly", "Buddy", "Daisy", "Rocky", "Maggie", "Jake", "Sophie", "Jack", "Sadie", "Toby", "Chloe", "Cody", "Bailey", "Buster", "Lola", "Duke", "Zoe", "Cooper", "Abby", "Riley", "Ginger", "Harley", "Roxy", "Bear", "Gracie", "Tucker", "Coco", "Murphy", "Sasha", "Lucky", "Lily", "Oliver", "Angel", "Sam", "Princess", "Oscar", "Emma", "Teddy", "Annie", "Winston", "Rosie"}

func randomDate() time.Time {
	now := time.Now()
	start := now.AddDate(-15, 0, 0)
	delta := now.Unix() - start.Unix()

	sec := rand.Int63n(delta) + start.Unix()
	return time.Unix(sec, 0)
}

func randomBreed() string {
	return breeds[random(0, len(breeds))]
}

func randomName() string {
	return names[random(0, len(names))]
}

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func getUUID() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return uuid.String()
}
