package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gumeeee/rest-api-in-gin/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login logs in a user
//
//	@Summary		Logs in a user
//	@Description	Logs in a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body	loginRequest	true	"User"
//	@Success		200	{object}	loginResponse
//	@Router			/api/v1/auth/login [post]
func (app *application) login(ctx *gin.Context) {
	var auth loginRequest

	if err := ctx.ShouldBindJSON(&auth); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if existingUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Somethin went wrong"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"expr":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	ctx.JSON(http.StatusOK, loginResponse{Token: tokenString})
}

// RegisterUser registers a new user
// @Summary		Registers a new user
// @Description	Registers a new user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		registerRequest	true	"User"
// @Success		201	{object}	database.User
// @Router			/api/v1/auth/register [post]
func (app *application) registerUser(ctx *gin.Context) {
	var register registerRequest

	if err := ctx.ShouldBindJSON(&register); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	register.Password = string(hashedPassword)
	user := database.User{
		Email:    register.Email,
		Password: register.Password,
		Name:     register.Name,
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}
