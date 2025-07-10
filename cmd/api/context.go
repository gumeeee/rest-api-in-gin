package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gumeeee/rest-api-in-gin/internal/database"
)

func (app *application) GetUserFromContext(ctx *gin.Context) *database.User {
	contextUser, exists := ctx.Get("user")
	if !exists {
		return &database.User{}
	}

	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{}
	}

	return user
}
