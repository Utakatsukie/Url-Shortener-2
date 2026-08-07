package auth

import (
	"Url-Shortener-2/configs"
	"Url-Shortener-2/pkg/res"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	*configs.Config
}

type AuthHandlerDeps struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(handler.Config.Auth.Secret)
		var payload LoginRequest
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			res.Json(w, err.Error(), 402)
			return
		}
		validate := validator.New()
		if err := validate.Struct(payload); err != nil {
			res.Json(w, err.Error(), 402)
			return
		}
		fmt.Println(payload)
		data := LoginResponse{
			Token: "123",
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Register")
	}
}
