package custom_errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

func HandleSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func HandleSuccessWithCookie(c *gin.Context, data interface{}, cookieName string, cookieValue string, expiryTime int) {
	c.SetCookie(cookieName, cookieValue, expiryTime, "/", "principalityofcogito.com", true, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// HandleError : error format
func HandleHttpError(c *gin.Context, err error) {
	// Get ErrorType
	errorType := GetErrorType(err)

	// Get Status Code of the ErrorType
	status := GetHttpStatusCode(errorType)

	// Check if there's additional context to the error
	errorContext := GetErrorContext(err)
	if errorContext != nil {
		c.JSON(status, gin.H{
			"error":   err.Error(),
			"context": errorContext,
		})
		return
	}

	// if status == 500 || (status == 400 && utils.GetEnvWithKey("ENV") == "production") {
	// 	sentry.CaptureException(err)
	// }
	// No error context to the error
	c.JSON(status, gin.H{"success": false, "error": err.Error()})
}

func HandleGrpcError(err error) error {
	// Get ErrorType
	errorType := GetErrorType(err)

	// Get Status Code of the ErrorType
	code := GetGrpcStatusCode(errorType)

	return status.Error(code, err.Error())
}

// func HandleCronErrors(err error) {
// 	// Get ErrorType
// 	errorType := GetErrorType(err)

// 	// Get Status Code of the ErrorType
// 	status := GetStatusCode(errorType)
// 	// if status == 500 || (status == 400 && utils.GetEnvWithKey("ENV") == "production") {
// 	// 	sentry.CaptureException(err)
// 	// }
// 	logger.Logger.Fatal(err)
// }
