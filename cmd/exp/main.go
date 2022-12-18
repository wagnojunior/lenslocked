package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Age  int
	Meta UserMeta
	Bio  template.HTML
}

type UserMeta struct {
	Visits int
}

func main() {
	// Pases the template <hello.gohtml>
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	// Creates and populates a variable of type User
	user := User{
		Name: "Wagno Lee",
		Age:  45,
		Meta: UserMeta{
			Visits: 4,
		},
		Bio: `<script>alert("HaHa, you have been h4x0r3d!");</script>`,
	}

	// Process the template
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

}
