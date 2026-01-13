package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var messages = make(map[int]Message)
var nextID = 1

func GetHandler(c echo.Context) error {
	var msgSlice []Message

	for _, msg := range messages {
		msgSlice = append(msgSlice, msg)
	}
	return c.JSON(http.StatusOK, &msgSlice)
}

func PostHandler(c echo.Context) error {
	var message Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not add the message",
		})
	}

	message.ID = nextID
	nextID++

	messages[message.ID] = message
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message  was successfully added",
	})
}

func PatchHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid id",
		})
	}

	var updatedMessage Message
	if err := c.Bind(&updatedMessage); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not update the message",
		})
	}

	if _, exists := messages[id]; !exists {

		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Message was not found",
		})
	}

	updatedMessage.ID = id
	messages[id] = updatedMessage

	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message was updated",
	})
}

func DeleteHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid id",
		})
	}

	if _, exists := messages[id]; !exists {

		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Message was not found",
		})
	}

	delete(messages, id)

	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message was deleted",
	})
}

func main() {
	e := echo.New()

	e.GET("/messages", GetHandler)
	e.POST("/messages", PostHandler)
	e.PATCH("/messages/:id", PatchHandler)
	e.DELETE("/messages/:id", DeleteHandler)

	e.Start(":8080")

}
