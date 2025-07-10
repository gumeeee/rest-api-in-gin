package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gumeeee/rest-api-in-gin/internal/database"
)

// CreateEvent creates a new event
//
//	@Summary		Creates a new event
//	@Description	Creates a new event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		database.Event	true	"Event"
//	@Success		201		{object}	database.Event
//	@Router			/api/v1/events [post]
//	@Security		BearerAuth
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

// GetEvent returns a single event
//
//	@Summary		Returns a single event
//	@Description	Returns a single event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	database.Event
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

// getEvents return all events
//
// @Summary Get all events
// @Description Returns all events
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {object} database.Event
// @Router /api/v1/events [get]
func (app *application) getAllEvents(ctx *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to retreive events"})
	}

	ctx.JSON(http.StatusOK, events)
}

// UpdateEvent updates an existing event
//
//	@Summary		Updates an existing event
//	@Description	Updates an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			event	body		database.Event	true	"Event"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [put]
//	@Security		BearerAuth
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

// DeleteEvent deletes an existing event
//
//	@Summary		Deletes an existing event
//	@Description	Deletes an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		204
//	@Router			/api/v1/events/{id} [delete]
//	@Security		BearerAuth
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

// AddAttendeeToEvent adds an attendee to an event
// @Summary		Adds an attendee to an event
// @Description	Adds an attendee to an event
// @Tags			attendees
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success		201		{object}	database.Attendee
// @Router			/api/v1/events/{id}/attendees/{userId} [post]
// @Security		BearerAuth
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

// GetAttendeesForEvent returns all attendees for a given event
//
//	@Summary		Returns all attendees for a given event
//	@Description	Returns all attendees for a given event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	[]database.User
//	@Router			/api/v1/events/{id}/attendees [get]
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

// DeleteAttendeeFromEvent deletes an attendee from an event
// @Summary		Deletes an attendee from an event
// @Description	Deletes an attendee from an event
// @Tags			attendees
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success		204
// @Router			/api/v1/events/{id}/attendees/{userId} [delete]
// @Security		BearerAuth
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

// GetEventsByAttendee returns all events for a given attendee
//
//	@Summary		Returns all events for a given attendee
//	@Description	Returns all events for a given attendee
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Attendee ID"
//	@Success		200	{object}	[]database.Event
//	@Router			/api/v1/attendees/{id}/events [get]
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
