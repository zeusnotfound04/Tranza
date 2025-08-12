package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zeusnotfound04/Tranza/utils" // Assuming utils package is in this path
)

var jwtSecret = []byte(utils.GetJWTSecret()) // Function to get secret from utils

// JWTAuthMiddleware provides a JWT authentication middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Debug logs
		fmt.Printf("DEBUG Auth: Request headers: %v\n", c.Request.Header)
		fmt.Printf("DEBUG Auth: Request cookies: %v\n", c.Request.Cookies())

		// First, try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("DEBUG Auth: Authorization header: '%s'\n", authHeader)
		
		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader { // No "Bearer " prefix
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not find bearer token in Authorization header"})
				c.Abort()
				return
			}
			fmt.Printf("DEBUG Auth: Using token from Authorization header\n")
		} else {
			// If no Authorization header, try to get token from HttpOnly cookie
			var err error
			tokenString, err = c.Cookie("access_token")
			fmt.Printf("DEBUG Auth: Cookie 'access_token' value: '%s', error: %v\n", tokenString, err)
			
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required: no token in header or cookie"})
				c.Abort()
				return
			}
			fmt.Printf("DEBUG Auth: Using token from cookie\n")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure that the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Debug logs
			fmt.Printf("DEBUG Auth: Token claims: %v\n", claims)
			fmt.Printf("DEBUG Auth: user_id claim: %v\n", claims["user_id"])
			
			// Set user claims in context
			c.Set("userID", claims["user_id"])
			c.Set("email", claims["email"])
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
