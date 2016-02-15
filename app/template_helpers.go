package app

import (
	"html/template"
	"net/http"
	"strings"
)

type Params struct {
	// Name of the current page, not always set
	Page string
	// LoggedIn is set to true when the current user
	// is logged in.
	LoggedIn bool
	User     User
}

func TmplParams(app *App, r *http.Request, page string) Params {
	user := GetUser(app.DB(), r)
	return Params{
		Page:     page,
		LoggedIn: len(user.Email) > 0,
		User:     user,
	}
}

// ----------------------

func TemplateHelpers() template.FuncMap {
	return template.FuncMap{
		"capitalize": Capitalize,
	}
}

// Capitalize capitalizes the given string.
func Capitalize(str string) string {
	if len(str) <= 0 {
		return ""
	}

	if len(str) == 1 {
		return strings.ToUpper(str)
	}

	return strings.ToUpper(str[:1]) + strings.ToLower(str[1:])
}
