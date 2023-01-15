package controllers

import "net/http"

// Currently the users controller defined in `controllers/users.go`
// uses the type `views.Template` defined in `views/template.go`.
// As a result, `constrollers/users.go` imports the `views` package
// and both are said to be **coupled**. In addition to that, the users
// controller is locked in to the `views.Template` type. Our goal is
// to define an interface that can replace the use of `views.Template`
// by the users controller so that any type that implements the
// interface is valid.

type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data interface{})
}
