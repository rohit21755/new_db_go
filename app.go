package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

func main() {

	dir := "./"
	db, err := New(dir, nil)
	if err != nil {
		fmt.Printf("Error", err)
	}

	employees := []User{
		{"John", "23", "2345423", "Google", Address{"banglore", "karnataka", "india", "302028"}},
		{"Sarah", "28", "9876543", "Microsoft", Address{"seattle", "washington", "usa", "98052"}},
		{"Raj", "31", "8765432", "Amazon", Address{"mumbai", "maharashtra", "india", "400001"}},
		{"Maria", "26", "3456789", "Apple", Address{"cupertino", "california", "usa", "95014"}},
		{"Yuki", "29", "7654321", "Sony", Address{"tokyo", "kanto", "japan", "100001"}},
	}

	for _, val := range employees {
		db.Write("users", val.Name, User{
			Name:    val.Name,
			Age:     val.Age,
			Contact: val.Contact,
			Company: val.Company,
			Address: val.Address,
		})
	}

	var records interface{}
	if records, err = db.ReadAll("users"); err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(records)

}
