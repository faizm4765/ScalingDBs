// package main

// import (
// 	"fmt"
// 	"database/sql"
// 	_ "github.com/lib/pq"
// 	"log"
// 	"sync"
// )

// type User struct {
// 	Name  string
// 	Id    int64
// }

// type Seat struct {
// 	Number int
// 	Row    int
// }

// func main() {
// 	//  book 200 seats for 200 users concurrently by spinning 200 goroutines
// 	//  each goroutine should book a seat for a user and print the user and seat

// 	//  we will have a db connection and update the seat as booked for the user

// 	//  these are the values used while creating posgres container locally
// 	connStr1 := "user=postgres dbname=postgres port=5432 user=postgres password=mysecretpassword sslmode=disable"
// 	db, err := sql.Open("postgres", connStr1)	
// 	if err != nil {
// 		log.Fatal(err)
// 	}else {
// 		log.Println("Database connection established successfully")
// 	}

// 	users := make([]User, 200)
// 	seats := make([]Seat, 200)
// 	for i := 0; i < 200; i++ {
// 		users[i] = User{
// 			Name:  "User" + string(i),
// 			Id:    int64(i),
// 		}

// 		seats[i] = Seat{
// 			Number: i + 1,
// 			Row:    (i / 10) + 1,
// 		}
// 	}

// 	var wg sync.WaitGroup
// 	for i := 0; i < 200; i++ {
// 		wg.Add(1)
// 		go func(i int) {
// 			defer wg.Done()
// 			user := users[i]
// 			// println("Booked seat", seat.Number, "in row", seat.Row, "for user", user.Name, "with id", user.Id)
// 			//  goroutines should book a seat for a user in the database
// 			//  update the seat as booked for the user
// 			//  find a seat which is not booked and book it for the user
// 			db.Exec("UPDATE flight_Seats SET user_id = $1 WHERE seat_id = (SELECT seat_id FROM flight_Seats WHERE user_id IS NULL LIMIT 1)", user.Id)
// 		}(i)
// 	}

// 	wg.Wait()

// 	//  print all the seats booked for users. Show output by print rows and x which is booked
// 	rows, err := db.Query("SELECT * FROM flight_Seats ORDER BY seat_id")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var seatStatus [200]string

// 	defer rows.Close()
// 	for rows.Next() {
// 		var userId sql.NullInt64
		
// 		//  fetch seat number from flight_Seats table
// 		//  print user id and seat number
// 		var seat_id int
// 		if err := rows.Scan(&seat_id, &userId); err != nil {
// 			log.Fatal(err)
// 		}
	
// 		if userId.Valid {
// 			log.Printf("User with id %d has booked seat %d", userId.Int64, seat_id)
// 			seatStatus[seat_id-1] = "X"
// 		} else {
// 			seatStatus[seat_id-1] = "O"
// 		}
// 	}

// 	seatsPerRow := 10
// 	for i := 0; i < 200; i++ {
// 		fmt.Print(seatStatus[i], " ")
// 		if (i + 1) % seatsPerRow == 0 {
// 			fmt.Println() // new row
// 		}
// 	}
// }