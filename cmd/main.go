package main

import (
	"errors"
	"fmt"
	"math/rand"
)

type QueryError struct {
	Message string
}

func (e QueryError) Error() string { return e.Message }

func executeQuery() (int, error) {
	randomBit := rand.Intn(2)
	if randomBit == 1 {
		return 0, errors.New("Unexpected error")
	} else {
		return 0, &QueryError{Message: "socket not found"}
	}

}

func main() {
	result, err := executeQuery()
	if e, ok := err.(*QueryError); ok {
		fmt.Println("QueryError:", e.Message)
	} else if err != nil {
		fmt.Println("An unexpected error occurred:", err)
	} else {
		fmt.Println("Query executed successfully, result:", result)
	}
}