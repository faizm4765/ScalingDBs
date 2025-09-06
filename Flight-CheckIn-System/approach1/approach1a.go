package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func initDB() {
	// Initialize database connection here
	connStr := "user=postgres dbname=postgres port=5432 user=postgres password=mysecretpassword sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}else {
		// Actually try to connect
		err = db.Ping()
		if err != nil {
			log.Fatal("Database is unreachable: ", err)
		}

		log.Println("Database connection established successfully")
	}
}

func resetDB() {
	//  reset data in users and seats table.
	_, err := db.Exec("UPDATE seats SET user_id = NULL")
	if err != nil {
		log.Fatal("Failed to reset database: ", err)
	}
}

func main() {
	initDB()
	resetDB()
	fmt.Println("Enter  user id for whom you want to book a seat:")

	var userId int
	fmt.Scanln(&userId)

	fmt.Println("Enter seat number you would like to book:")
	var seatNumber int
	fmt.Scanln(&seatNumber)

	fmt.Printf("User ID: %d, Seat Number: %d\n", userId, seatNumber)
	bookSeat(userId, seatNumber)
}

func bookSeat(userId int, seatNumber int) {
	fmt.Printf("Booking seat %d for user %d\n", seatNumber, userId)
	
	// begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction: ", err)
		return
	}

	// book the seat
	_, err = tx.Exec("UPDATE seats SET user_id = $1 WHERE seat_id = $2", userId, seatNumber)
	if err != nil {
		tx.Rollback()
		log.Fatal("Failed to book seat: ", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Failed to commit transaction: ", err)
		return
	}
}