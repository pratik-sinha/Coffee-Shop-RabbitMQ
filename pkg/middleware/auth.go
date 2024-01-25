package middleware

import (
	"coffee-shop/pkg/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (mw *MiddlewareManager) AuthenticateSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(mw.cfg.Cookie.Name)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "error": fmt.Sprintf("Error while retrieving cookie : %s", err.Error())})
			return
		}

		payload, err := mw.t.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "error": fmt.Sprintf("Error while verifying token : %s", err.Error())})
			return
		}

		if time.Since(payload.ExpiredAt) > 0 {
			//Expired token
			utils.DeleteSessionCookie(c, mw.cfg)
			mw.sessUC.DeleteByID(c, payload.Key)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "error": "Error while verifying token : token expired"})
			return
		}

		session, err := mw.sessUC.GetSessionByID(c, payload.Key)
		if err != nil {
			utils.DeleteSessionCookie(c, mw.cfg)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "error": fmt.Sprintf("Error while retrieving session : %s", err.Error())})
			return
		}
		c.Set("profile_id", session.ProfileID)
		c.Set("session_id", payload.Key)
		c.Next()
	}
}
