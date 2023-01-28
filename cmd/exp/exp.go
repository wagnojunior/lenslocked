package main

import (
	std_context "context"
	"fmt"

	"github.com/wagnojunior/lenslocked/context"
	"github.com/wagnojunior/lenslocked/models"
)

func main() {
	ctx := std_context.Background()

	user := models.User{
		Email: "test@test.com",
	}

	ctx = context.WithUser(ctx, &user)

	retrievedUser := context.User(ctx)
	fmt.Println(retrievedUser.Email)
}
