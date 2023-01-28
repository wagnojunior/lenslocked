package context

import (
	"context"

	"github.com/wagnojunior/lenslocked/models"
)

type key string

const (
	userKey key = "user"
)

// WithUser returns a context with an user `user` associated with the key `userKey`
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User returns the user stored in the context `ctx` and associated with the key `userKey`
func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)

	// Asserts if `val` is of type `*models.User`. The most likely case the assertion fails is when nothing was ever stored in the context
	user, ok := val.(*models.User)
	if !ok {
		return nil
	}

	return user
}
