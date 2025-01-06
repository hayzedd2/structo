package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/structo/generator"
	"github.com/structo/parser"
	"github.com/structo/types"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://structo.vercel.app",
		},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))
	r.POST("/generate-mock", handleGenerateMock)
	r.GET("/ping", pingServer)

	go startWarmup()

	log.Fatal(r.Run(":8080"))
}

func pingServer(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Server warm up.",
	})
}

func startWarmup() {
	for {
		time.Sleep(14 * time.Minute)
		_, err := http.Get("https://structo.onrender.com/ping")
        if err != nil {
            log.Println("Error while pinging server:", err)
        } 
		log.Println("Sent warmup message to keep server alive.")
	}
}
func handleGenerateMock(c *gin.Context) {
	var req types.MockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request data"})
		return
	}
	fields, err := parser.ParseTypeOrInterface(req.TypescriptInterface, req.Language)
	fmt.Println(fields,err)
	fmt.Println(req.TypescriptInterface, req.Count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mockData := generator.GenerateMockObjects(fields, req.Count)
	c.JSON(http.StatusOK, mockData)
}
