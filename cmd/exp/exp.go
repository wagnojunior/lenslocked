package main

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, favoriteColorKey, "blue")
	value := ctx.Value(favoriteColorKey)
	strValue, ok := value.(int)
	if !ok {
		fmt.Println("Cannot asset to the desired type")
		return
	}

	fmt.Println(strValue)
	// fmt.Println(strings.HasPrefix(strValue, "b"))
}
