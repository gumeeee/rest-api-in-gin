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

	user := app.GetUserFromContext(ctx)
	event.OwnerId = user.Id

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

	user := app.GetUserFromContext(ctx)
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

	if existingEvent.OwnerId != user.Id {
		ctx.JSON(http.StatusForbidden,
			gin.H{"error": "You are not authorized to update this event"})
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

	user := app.GetUserFromContext(ctx)
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

	if existingEvent.OwnerId != user.Id {
		ctx.JSON(http.StatusForbidden,
			gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to delete event"})
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Event deleted successfully"})
}

func (app *application) AddAttendeeToEvent(ctx *gin.Context) {
	eventId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userId, err := strconv.Atoi(ctx.Query("userId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}
	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive user"})
		return
	}
	if userToAdd == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := app.GetUserFromContext(ctx)
	if event.OwnerId != user.Id {
		ctx.JSON(http.StatusForbidden,
			gin.H{"error": "You are not authorized to add an attendee to this event"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendeeId(event.Id, userToAdd.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive attendee"})
		return
	}
	if existingAttendee != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Attendee already exists"})
		return
	}

	attendee := database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	ctx.JSON(http.StatusCreated, attendee)
}

func (app *application) GetAttendeesForEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEventId(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to retreive attendees for event"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (app *application) DeleteAttendeeFromEvent(ctx *gin.Context) {
	eventId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
		return
	}

	userId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user Id"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Something went wrong"})
		return
	}

	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.GetUserFromContext(ctx)
	if event.OwnerId != user.Id {
		ctx.JSON(http.StatusForbidden,
			gin.H{"error": "You are not authorized to delete an attendee from this event"})
		return
	}

	err = app.models.Attendees.Delete(userId, eventId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to delete attendee from event"})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Attendee deleted successfully"})
}

func (app *application) GetEventsByAttendee(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee Id"})
		return
	}

	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to get events"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
