package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{

		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)

		v1.GET("/events/:id/attendees", app.GetAttendeesForEvent)

		v1.GET("/attendees/:id/events", app.GetEventsByAttendee)

		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/events", app.createEvent)
		authGroup.PUT("/events/:id", app.updateEvent)
		authGroup.DELETE("/events/:id", app.deleteEvent)
		authGroup.POST("/events/:id/attendees/:userId", app.AddAttendeeToEvent)
		authGroup.DELETE("/events/:id/attendees/:userId", app.DeleteAttendeeFromEvent)
	}

	g.GET("/swagger/*any", func(ctx *gin.Context) {
		if ctx.Request.RequestURI == "/swagger/" {
			ctx.Redirect(302, "/swagger/index.html")
		}

		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(ctx)
	})

	return g
}
