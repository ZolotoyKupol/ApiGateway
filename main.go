package main

import (
	"net/http"
	//"time"

	"github.com/gin-gonic/gin"
)

type Guest struct {
	ID string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	//CheckIn time.Time `json:"check_in"`
	//CheckOut time.Time `json:"check_out"`
	RoomID int `json:"room_id"`
}

var guests = []Guest{
	{ID: "1", FirstName: "Egor", LastName: "Dmitrienko", RoomID: 10},
	{ID: "2", FirstName: "Egor", LastName: "Dmitrienko", RoomID: 10},
	{ID: "3", FirstName: "Egor", LastName: "Dmitrienko", RoomID: 10},
}

func getAllGuests(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, guests)
}

func postGuests(c *gin.Context) {
	var newGuest Guest

	if err := c.BindJSON(&newGuest); err != nil {
		return
	}
	
	guests = append(guests, newGuest)
	c.IndentedJSON(http.StatusCreated, newGuest)
}

func getGuestByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range guests {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "guest not found"})
}

func updateGuestbyID(c *gin.Context) {
	id := c.Param("id")

	var updatedGuest Guest
	if err := c.BindJSON(&updatedGuest); err != nil {
		return
	}

	for i, a := range guests {
		if a.ID == id {
			updatedGuest.ID = a.ID
			guests[i] = updatedGuest
			c.IndentedJSON(http.StatusOK, updatedGuest)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "guest not found"})
}

func deleteByID(c *gin.Context) {
	id := c.Param("id")

	for i, a := range guests {
		if a.ID == id {
			guests = append(guests[:i], guests[i+1])
			break
		}
	}
	c.IndentedJSON(http.StatusOK, guests)
}


func main() {
	router := gin.Default()
	router.GET("/guests", getAllGuests)
	router.GET("/guests/:id", getGuestByID)
	router.PUT("/guests/:id", updateGuestbyID)
	router.POST("/guests", postGuests)
	router.DELETE("/guests/:id", deleteByID)

	router.Run("localhost:8080")
}