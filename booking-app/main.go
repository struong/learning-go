package main

import (
	"fmt"
	"sync"
	"time"
)

// package level variables
const conferenceName = "Go Conference"
const conferenceTickets int = 50

var remainingTickets = uint(conferenceTickets)
var bookings = make([]UserData, 0)

// type keyword creates a new type
type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

// waits for the launched go routine to finish
var wg = sync.WaitGroup{}

func main() {
	// fmt.Printf("conferenceTickets is %T conferenceTickets is %T remainingTickets is %T\n", conferenceName, conferenceTickets, remainingTickets)

	greetUser()

	firstName, lastName, email, userTickets := getUserInput()
	isValidName, isValidEmail, isValidTicketNumber := validateUserInput(firstName, lastName, email, userTickets, remainingTickets)

	for {
		if isValidName && isValidEmail && isValidTicketNumber {

			bookTicket(userTickets, firstName, lastName, email)

			// Increments the number of goroutines to wait for
			wg.Add(1)
			// go ... starts a new go routine
			go sendTicket(userTickets, firstName, lastName, email)

			firstNames := getFirstNames()

			fmt.Printf("The first names of bookings are: %v\n", firstNames)

			if remainingTickets == 0 {
				fmt.Println("Our conference is booked out. Come back next year.")
				// break
			}
		} else {
			if !isValidName {
				fmt.Println("first name or last name you entered is too short")
			}

			if !isValidEmail {
				fmt.Println("email address you entered does not contain @ sign")
			}

			if !isValidTicketNumber {
				fmt.Println("number of tickets you entered is invalid")
			}
		}

		// Wait blocks until the WaitGroup counter is zero.
		wg.Wait()
	}
}

func greetUser() {
	fmt.Printf("Welcome to %v booking application\n", conferenceName)

	fmt.Printf("We have a total of %v tickets and %v are still available.\n", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend.")
}

func getFirstNames() []string {
	firstNames := []string{}
	for _, booking := range bookings {
		firstNames = append(firstNames, booking.firstName)
	}

	return firstNames
}

func getUserInput() (string, string, string, uint) {
	var firstName, lastName, email string
	var userTickets uint

	// ask the user for input
	fmt.Println("Enter your first name:")
	fmt.Scan(&firstName)

	fmt.Println("Enter your last name:")
	fmt.Scan(&lastName)

	fmt.Println("Enter your email:")
	fmt.Scan(&email)

	fmt.Println("Enter number of tickets:")
	fmt.Scan(&userTickets)

	return firstName, lastName, email, userTickets
}

func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	remainingTickets = remainingTickets - uint(userTickets)

	// create a map for a user
	userData := UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)
	fmt.Printf("List of bookings is %v\n", bookings)

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", remainingTickets, conferenceName)
}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {
	time.Sleep(10 * time.Second)

	ticket := fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)

	fmt.Println("####################")
	fmt.Printf("Sending ticket:\n %v \nto email address %v\n", ticket, email)
	fmt.Println("####################")

	// Done decrements the WaitGroup counter by one.
	wg.Done()
}
