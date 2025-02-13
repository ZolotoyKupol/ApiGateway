package main

import (
	"net/http"
	//"time"
	"context"
	"fmt"
	"log"
	//"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Guest struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoomID int `json:"room_id"`
}



var db *pgx.Conn

func initDB(connString string)  error {
	var err error
	db, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к БД: %v", err)
	}
	fmt.Println("Успешно подключено к БД")
	return nil
}

// var guests = []Guest{}
	//{ID: "1", FirstName: "Egor", LastName: "Dmitrienko", RoomID: 10},
	//{ID: "2", FirstName: "Egor", LastName: "Dmitrienko", RoomID: 10},
	//{ID: "3", FirstName: "Egor", LastName: "Dmitrienko", RoomID: 10},
//}


// "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres"


func main() {
	connString := "postgres://postgres:postgres@localhost:5432/ApiGatway_db"
	if err := initDB(connString); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	router := gin.Default()
	router.GET("/guests", getAllGuests)
	router.GET("/guests/:id", getGuestByID)
	router.PUT("/guests/:id", updateGuestbyID)
	router.POST("/guests", postGuests)
	router.DELETE("/guests/:id", deleteByID)

	router.Run("localhost:8080")
}

// func insertGuest(Guest) (*Guest, error) {
//     conn.Exec("insert into guest ...")
// }
// func updateGuest(Guest) id error
// func getGuest(id) id error
// func deleteGuest(Guest) id error
// func getAllGuest() id error

func getAllGuests(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT id, first_name, last_name, room_id FROM guests")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var guests []Guest
	for rows.Next() {
		var guest Guest
		err := rows.Scan(&guest.ID, &guest.FirstName, &guest.LastName, &guest.RoomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		guests = append(guests, guest)
	}
	c.IndentedJSON(http.StatusOK, guests)
}

func postGuests(c *gin.Context) {
	var newGuest Guest

	if err := c.BindJSON(&newGuest); err != nil {
		// возварщай ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newGuest.FirstName == "" || newGuest.LastName == "" || newGuest.RoomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "все поля должны быть заполнены"})
		return
	}

	err := db.QueryRow(context.Background(), "INSERT INTO guests (first_name, last_name, room_id) VALUES ($1, $2, $3) RETURNING id", newGuest.FirstName, newGuest.LastName, newGuest.RoomID).Scan(&newGuest.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка добавления гостя"})
		return
	}
	
	c.IndentedJSON(http.StatusCreated, newGuest)
}

func getGuestByID(c *gin.Context) {
	id := c.Param("id")

	var guest Guest
	err := db.QueryRow(context.Background(), "SELECT id, first_name, last_name, room_id FROM guests WHERE id = $1", id).Scan(&guest.ID, &guest.FirstName, &guest.LastName, &guest.RoomID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "гость не найден"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "ошибка получения гостя"})
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

	res, err := db.Exec(context.Background(), "UPDATE guests SET first_name = $1, last_name = $2, room_id = $3 WHERE id = $4", updatedGuest.FirstName, updatedGuest.LastName, updatedGuest.RoomID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления гостя"})
		return
	}

	if res.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "гость не найден"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "гость успешно обновлен"})
}

func deleteByID(c *gin.Context) {
	id := c.Param("id")

	res, err := db.Exec(context.Background(), "DELETE FROM guests WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка удаления гостя"})
		return
	}

	if res.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "гость не найден"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "гость успешно удален"})
}
