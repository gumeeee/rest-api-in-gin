package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gumeeee/rest-api-in-gin/internal/database"
)

func (app *application) createEvent(ctx *gin.Context) {
	var event database.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}

	err := app.models.Events.Insert(&event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to create event"})
		return
	}

	ctx.JSON(http.StatusCreated, event)
}

func (app *application) getEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)
	if event == nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to retreive event"})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (app *application) getAllEvents(ctx *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to retreive events"})
	}

	ctx.JSON(http.StatusOK, events)
}

func (app *application) updateEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "Invalid event ID"})
		return
	}

	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to retreive event"})
		return
	}

	if existingEvent == nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{"error": "Event not found"})
		return
	}

	updatedEvent := &database.Event{}

	if err := ctx.ShouldBindJSON(&updatedEvent); err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = id

	if err := app.models.Events.Update(updatedEvent); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to update event"})
		return
	}

	ctx.JSON(http.StatusOK, updatedEvent)
}

func (app *application) deleteEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "Invalid event ID"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to delete event"})
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Event deleted successfully"})
}
