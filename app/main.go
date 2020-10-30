package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	_ "github.com/cowai/sam-gin-app/statik"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
)

var ginLambda *ginadapter.GinLambda

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {

		r := gin.Default()
		t, err := loadTemplate()
		if err != nil {
			panic(err)
		}

		r.SetHTMLTemplate(t)
		r.GET("/page", func(c *gin.Context) {
			c.HTML(http.StatusOK, "/index.tmpl", gin.H{})
		})

		r.GET("/api/pet/:id", func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"id":   c.Param("id"),
					"name": randomName(),
					"time": randomDate(),
				},
			)
		})

		ginLambda = ginadapter.New(r)
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

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	// 仮想ファイルシステム上を走査する
	err = fs.Walk(statikFS, "/", func(path string, info os.FileInfo, err error) error {
		// ディレクトリはスキップ
		if info.IsDir() {
			return nil
		}
		r, err := statikFS.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()

		// データ読み出し
		h, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		// テンプレートにpath名で読み出したデータをパースして格納
		t, err = t.New(path).Parse(string(h))
		if err != nil {
			return err
		}

		return nil
	})

	return t, err
}
