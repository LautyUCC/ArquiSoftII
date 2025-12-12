package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := []string{"admin123", "usuario123"}
	
	for _, pwd := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Error hasheando %s: %v\n", pwd, err)
			continue
		}
		fmt.Printf("Password: %s\n", pwd)
		fmt.Printf("Hash: %s\n\n", string(hash))
	}
}

