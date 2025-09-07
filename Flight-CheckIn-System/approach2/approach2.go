package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
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

	bookSeat(userId)
	printSeats()
}

func bookSeat(userId int) {
	userName := fetchUserName(userId)
	if userName == "" {
		log.Println("User not found")
		return
	}

	fmt.Printf("Booking seat for user %s\n", userName)

	// begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction: ", err)
		return
	}

	//  get all seats and then select random seat from that list
	allSeats := fetchAllSeats()
	seatNumber := selectRandomSeat(allSeats)

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

func fetchUserName(userId int) string {
	var userName string
	err := db.QueryRow("SELECT user_name FROM users WHERE user_id = $1", userId).Scan(&userName)
	if err != nil {
		log.Fatal("Failed to fetch user name: ", err)
		return ""
	}

	return userName
}

func fetchAllSeats() []int {
	seats := make([]int, 0)
	rows, err := db.Query("SELECT seat_id FROM seats WHERE user_id IS NULL")
	if err != nil {
		log.Fatal("Failed to fetch seats: ", err)
		return seats
	}
	defer rows.Close()

	for rows.Next() {
		var seatID int
		if err := rows.Scan(&seatID); err != nil {
			log.Fatal("Failed to scan seat: ", err)
		}
		seats = append(seats, seatID)
	}
	return seats
}

func selectRandomSeat(allSeats []int) int {
	if len(allSeats) == 0 {
		log.Fatal("No available seats found")
	}

	// Select a random seat from the available seats
	return allSeats[rand.Intn(len(allSeats))]
}