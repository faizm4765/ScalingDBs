package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"sync"
)

type User struct {
	ID   int
	Name string
}

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

		// log.Println("Database connection established successfully")
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
	resetDB()

	users := fetchAllUsers()
	var wg sync.WaitGroup
	wg.Add(len(users))

	// simulating concurrent seat booking
	for _, user := range users {
		go func(u User) {
			defer wg.Done()
			err := bookSeat(u)
			if err != nil {
				log.Println("%s could not book seat", u.Name)
			}
		}(user)
	}

	wg.Wait()
	printSeats()
}

func bookSeat(user User) error {
	userName := user.Name
	userId := user.ID

	// begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction: ", err)
		return err
	}

	//  get all seats and then select random seat from that list
	allSeats := fetchAllSeats()
	seatNumber := selectRandomSeat(allSeats)

	// book the seat
	_, err = tx.Exec("UPDATE seats SET user_id = $1 WHERE seat_id = $2", userId, seatNumber)
	if err != nil {
		tx.Rollback()
		log.Fatal("Failed to book seat: ", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Failed to commit transaction: ", err)
		return err
	}

	//  print this log in green
	fmt.Printf("\033[32mSeat %d booked for user %s\n\033[0m", seatNumber, userName)
	return nil
}

func fetchAllUsers() []User {
	users := make([]User, 0)
	rows, err := db.Query("SELECT user_id, user_name FROM users")
	if err != nil {
		log.Fatal("Failed to fetch users: ", err)
		return users
	}

	defer rows.Close()
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			log.Fatal("Failed to scan user: ", err)
		}

		users = append(users, user)
	}

	return users
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