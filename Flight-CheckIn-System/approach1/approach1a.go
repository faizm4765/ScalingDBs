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

func printSeats() {
	fmt.Println("Current seats layout (5x4):")

	// query all seats ordered by seat_id
	rows, err := db.Query("SELECT seat_id, user_id FROM seats ORDER BY seat_id")
	if err != nil {
		log.Fatal("Failed to query seats: ", err)
	}
	defer rows.Close()

	seats := make([]int, 0, 20) // store user_id; 0 means empty
	for rows.Next() {
		var seatID, userID sql.NullInt32
		err := rows.Scan(&seatID, &userID)
		if err != nil {
			log.Fatal(err)
		}

		if userID.Valid {
			seats = append(seats, int(userID.Int32))
		} else {
			seats = append(seats, 0)
		}
	}

	// print as 5x5 grid
	for i := 0; i < 20; i++ {
		if seats[i] == 0 {
			fmt.Print(". ")
		} else {
			fmt.Print("x ")
		}

		if (i+1)%5 == 0 {
			fmt.Println()
		}
	}
}


func main() {
	initDB()
	// resetDB()
	fmt.Println("Enter  user id for whom you want to book a seat:")

	var userId int
	fmt.Scanln(&userId)

	fmt.Println("Enter seat number you would like to book:")
	var seatNumber int
	fmt.Scanln(&seatNumber)

	fmt.Printf("User ID: %d, Seat Number: %d\n", userId, seatNumber)
	bookSeat(userId, seatNumber)
	printSeats()
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