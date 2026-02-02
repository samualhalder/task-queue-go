package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	jsonresponse "github.com/samualhalder/task-queue-go/internal/json_response"
)

type adminPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (app *Application) GetToken(w http.ResponseWriter, r *http.Request) {
	var payload adminPayload
	if err := jsonresponse.ReadJson(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := jsonresponse.Validator.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	admin, err := app.Store.Admin.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		app.notFound(w, r, err)
		return
	}
	fmt.Print("val of admin is", *admin)
	if err := admin.Password.Check(payload.Password); err != nil {
		app.authorizationError(w, r, err)
		return
	}
	claims := jwt.MapClaims{
		"sub": admin.Id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.Config.AuthCnf.Issuer,
		"aud": app.Config.AuthCnf.Auditor,
	}
	token, err := app.Auth.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonresponse.Success(w, http.StatusOK, "Token will expire after an hour", token); err != nil {
		app.internalServerError(w, r, err)
	}
}
