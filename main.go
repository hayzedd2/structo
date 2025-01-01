package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/structo/generator"
	"github.com/structo/parser"
	"github.com/structo/types"
)

func main() {
	r := gin.Default()
	r.POST("/generate-mock", handleGenerateMock)
	log.Fatal(r.Run(":8080"))
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
