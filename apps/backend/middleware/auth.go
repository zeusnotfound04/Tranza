package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/utils" // Assuming utils package is in this path
)

var jwtSecret = []byte(utils.GetJWTSecret()) // Function to get secret from utils

// JWTAuthMiddleware provides a JWT authentication middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Debug logs
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
			fmt.Printf("DEBUG Auth: Token method: %v\n", token.Method)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Printf("DEBUG Auth: Invalid signing method: %v\n", token.Method)
				return nil, jwt.ErrSignatureInvalid
			}

			// Debug: Print the actual JWT secret being used
			secretStr := utils.GetJWTSecret()
			fmt.Printf("DEBUG Auth: JWT secret from utils.GetJWTSecret(): '%s' (length: %d)\n", secretStr, len(secretStr))
			fmt.Printf("DEBUG Auth: jwtSecret variable: '%s' (length: %d)\n", string(jwtSecret), len(jwtSecret))

			fmt.Printf("DEBUG Auth: Using JWT secret from utils.GetJWTSecret()\n")
			return []byte(secretStr), nil
		})

		fmt.Printf("DEBUG Auth: Token parsing result - Valid: %v, Error: %v\n", token != nil && token.Valid, err)

		if err != nil {
			fmt.Printf("DEBUG Auth: Token parsing failed with error: %v\n", err)
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
			fmt.Printf("DEBUG Auth: Token is valid! Claims type assertion successful\n")
			fmt.Printf("DEBUG Auth: Token claims: %v\n", claims)
			fmt.Printf("DEBUG Auth: user_id claim: %v\n", claims["user_id"])

			// Parse user_id string to UUID before setting in context
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				fmt.Printf("DEBUG Auth: user_id claim is not a string: %T\n", claims["user_id"])
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
				c.Abort()
				return
			}

			userUUID, err := uuid.Parse(userIDStr)
			if err != nil {
				fmt.Printf("DEBUG Auth: Failed to parse user_id as UUID: %v\n", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id format in token"})
				c.Abort()
				return
			}

			// Set user claims in context with proper types
			c.Set("user_id", userUUID) // Now setting as uuid.UUID instead of string
			c.Set("email", claims["email"])
			c.Set("username", claims["username"])
			fmt.Printf("DEBUG Auth: Successfully set user context with UUID: %s, proceeding to next handler\n", userUUID)
			c.Next()
		} else {
			fmt.Printf("DEBUG Auth: Token validation failed - Claims OK: %v, Token Valid: %v\n", ok, token.Valid)
			if !ok {
				fmt.Printf("DEBUG Auth: Claims type assertion failed - actual type: %T\n", token.Claims)
			}
			if !token.Valid {
				fmt.Printf("DEBUG Auth: Token is not valid\n")
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
