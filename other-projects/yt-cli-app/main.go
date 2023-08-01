package main

// Following along with this Youtube tutorial
// https://www.youtube.com/watch?v=yyUHQIec83I&t=4831s

import (
	"fmt"
	"strings"
)

func main() {
	const conferenceTickets int = 50
	var remainingTickets uint = 50
	conferenceName := "Go Conference"
	bookings := []string{}

	fmt.Printf("Welcome to the %v booking application.\nWe have a total of %v tickets and %v are still available.\n", conferenceName, conferenceTickets, remainingTickets)

	for {
		var firstName string
		var lastName string
		var email string
		var userTickets uint

		fmt.Println("Enter your first name: ")
		fmt.Scanln(&firstName)

		fmt.Println("Enter your last name: ")
		fmt.Scanln(&lastName)

		fmt.Println("Enter your e-mail: ")
		fmt.Scanln(&email)

		fmt.Println("Enter number of tickets: ")
		fmt.Scanln(&userTickets)

		remainingTickets -= userTickets
		bookings = append(bookings, firstName+" "+lastName)

		fmt.Printf("Thank you %v %v for booking %v tickets.  You'll receive a confirmation e-mail at %v\n", firstName, lastName, userTickets, email)
		fmt.Printf("%v tickets remaining for %v\n", remainingTickets, conferenceName)

		firstNames := []string{}
		for _, booking := range bookings {
			var names = strings.Fields(booking)
			firstNames = append(firstNames, names[0])
		}

		fmt.Printf("The first names %v\n", firstNames)

	}
}
