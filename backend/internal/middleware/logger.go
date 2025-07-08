package middleware

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// // Logger middleware для логирования запросов и ошибок
// func Logger() gin.HandlerFunc {
// 	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
// 		logLine := fmt.Sprintf("[INFO] %s | %d | %s | %s | %s | %s",
// 			param.TimeStamp.Format(time.RFC3339),
// 			param.StatusCode,
// 			param.Method,
// 			param.Path,
// 			param.ClientIP,
// 			param.Latency,
// 		)

// 		// Логируем все запросы
// 		log.Printf(logLine)

// 		if param.StatusCode >= 400 {
// 			log.Printf("[ERROR] %s | %d | %s | %s | %s | %s",
// 				param.TimeStamp.Format(time.RFC3339),
// 				param.StatusCode,
// 				param.Method,
// 				param.Path,
// 				param.ClientIP,
// 				param.ErrorMessage,
// 			)
// 		}

// 		return logLine
// 	})
// }

// func ErrorLogger() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Next()

// 		if len(c.Errors) > 0 {
// 			for _, err := range c.Errors {
// 				log.Printf("[ERROR] %s | %s | %s | %s",
// 					time.Now().Format(time.RFC3339),
// 					c.Request.Method,
// 					c.Request.URL.Path,
// 					err.Error(),
// 				)
// 			}
// 		}
// 	}
// }
