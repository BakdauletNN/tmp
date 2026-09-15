package main

import "github.com/gin-gonic/gin"

func main() {
    router := gin.Default()
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Start succes code 200",
        })
    })
    
    router.Run(":8080") // Starts the server on port 8080
}
