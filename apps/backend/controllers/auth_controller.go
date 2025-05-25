package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(ctx *gin.Context) {
	var input models.User
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid input"})
		return
	}

	existingUser, _ := repositories.GetUserByEmail(input.Email)

	if existingUser != nil {
		ctx.JSON(http.StatusConflict , gin.H{"error" : "User Already exists"})
		return
	}

	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(input.Password) , bcrypt.DefaultCost)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError , gin.H{"error" : "Failed to Hash the password"} )
		return
	}

	input.Password = string(hashedPassword)

	if err := repositories.CreateUser(&input); err != nil {
		ctx.JSON(http.StatusInternalServerError , gin.H{"error" : "Failed to create the user"} )
		return
		
	}

	token , err := utils.GenerateJWT(input.ID , input.Email , input.Username  )

	if err != nil {
		ctx.JSON(http.StatusInternalServerError , gin.H{"error" : "Failed to generate the JWT token"})
		return 
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"token":   token,
	})
}

	


func LoginHandler(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := repositories.GetUserByEmail(input.Email)
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email , user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Login successful",
		"token":     token,
		"expiresIn": time.Hour * 24,
	})
}
