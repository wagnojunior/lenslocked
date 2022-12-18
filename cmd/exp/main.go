package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Age  int
	Meta UserMeta
	Bio  string
	Week []int
}

type UserMeta struct {
	Visits  int
	AvgTime int
}

func main() {
	// Pases the template <hello.gohtml>
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	week := []int{1, 2, 3, 4, 5, 6, 7}

	// Creates and populates a variable of type User
	user := User{
		Name: "Wagno Lee",
		Age:  45,
		Meta: UserMeta{
			Visits:  4,
			AvgTime: 3,
		},
		Bio:  "This is my bio. Do you like it?",
		Week: week,
	}

	// Process the template
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

}
