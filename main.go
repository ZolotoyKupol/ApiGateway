package main

import (
	"net/http"
	"context"
	"fmt"
	"log"


	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Guest struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoomID    int    `json:"room_id"`
}

var db *pgx.Conn

func initDB(connString string) error {
	var err error
	db, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("couldn't connect to db: %v", err)
	}
	fmt.Println("successfully connected")
	return nil
}

func getRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/guests", getAllGuests)
	router.GET("/guests/:id", getGuestByID)
	router.PUT("/guests/:id", updateGuestbyID)
	router.POST("/guests", postGuests)
	router.DELETE("/guests/:id", deleteByID)

	return router
}



func main() {
	connString := "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
	if err := initDB(connString); err != nil {
		log.Fatalf("initialization error db: %v", err)
	}

	router := getRouter()
	router.Run(":8080")
}


func getAllGuests(c *gin.Context) {
	guests, err := fetchAllGuests(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.IndentedJSON(http.StatusOK, guests)
}

func postGuests(c *gin.Context) {
	var newGuest Guest

	if err := c.BindJSON(&newGuest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newGuest.FirstName == "" || newGuest.LastName == "" || newGuest.RoomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields must be filled in"})
		return
	}

	id, err := createGuest(db, newGuest)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error adding a guest"})
		return
	}

	newGuest.ID = id
	c.IndentedJSON(http.StatusCreated, newGuest)
}

func getGuestByID(c *gin.Context) {
	id := c.Param("id")

	guest, err := getGuest(db, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	
	if guest == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "guest not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, guest)
}

func updateGuestbyID(c *gin.Context) {
	id := c.Param("id")
	var updatedGuest Guest
	if err := c.BindJSON(&updatedGuest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	succsess, err := updateGuest(db, id, updatedGuest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update error"})
		return
	}

	if !succsess {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "guest not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "update successfully"})
}

func deleteByID(c *gin.Context) {
	id := c.Param("id")

	deleted, err := deleteGuest(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete error"})
		return
	}

	if !deleted {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "guest not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "delete successfuly"})
}


func createGuest (conn *pgx.Conn, guest Guest) (string, error) {
	var id string
	err := conn.QueryRow(context.Background(), "INSERT INTO guests (first_name, last_name, room_id) VALUES ($1, $2, $3) RETURNING id", guest.FirstName, guest.LastName, guest.RoomID).Scan(&id)
	return id, err
}

func deleteGuest (conn *pgx.Conn, id string) (bool, error) {
	res, err := conn.Exec(context.Background(), "DELETE FROM guests WHERE id = $1", id)
	if err != nil {
		return false, fmt.Errorf("error deleting guest: %v", err)
	}
	
	if res.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

func updateGuest (conn *pgx.Conn, id string, guest Guest) (bool, error) {
	res, err := conn.Exec(context.Background(), "UPDATE guests SET first_name = $1, last_name = $2, room_id = $3 WHERE id = $4", guest.FirstName, guest.LastName, guest.RoomID, id)
	if err != nil {
		return false, fmt.Errorf("error updating guest: %v", err)
	}

	if res.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

func getGuest (conn *pgx.Conn, id string) (*Guest, error) {
	var guest Guest
	err := conn.QueryRow(context.Background(), "SELECT id, first_name, last_name, room_id FROM guests WHERE id = $1", id) .Scan(&guest.ID, &guest.FirstName, &guest.LastName, &guest.RoomID)
	if err != nil{
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching guest: %v", err)
	}
	return &guest, nil

	}


func fetchAllGuests (conn *pgx.Conn) ([]Guest, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, first_name, last_name, room_id FROM guests")
	if err != nil {
		return nil, fmt.Errorf("error fetching guests: %v", err)
	}

	defer rows.Close()

	var guests []Guest
	for rows.Next() {
		var guest Guest
		err := rows.Scan(&guest.ID, &guest.FirstName, &guest.LastName, &guest.RoomID)
		if err != nil {
			return nil, fmt.Errorf("error scanning guest: %v", err)
		}
		guests = append(guests, guest)

	}

	return guests, nil
}