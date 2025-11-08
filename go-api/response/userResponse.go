package response

import (
	"go-api/core"
	"go-api/helper"
	"net/http"

	"github.com/go-chi/chi"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	res, err := core.GetAllUsers()
	if err != nil {
		helper.Response(helper.Response_struct{Error: "Cannot get the users with the function"}, w, http.StatusBadRequest)
		return
	}

	// this need to be switched on the future, probably don´t use this implementation on a real db
	for _, user := range res {
		helper.Response(helper.Response_struct{Data: "Seu nome é: " + user.Name + " e o seu email é: " + user.Email}, w, http.StatusOK)
	}

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	id_codified, err := helper.Encrypt(id)
	if err != nil {
		helper.Response(helper.Response_struct{Error: "error to codify"}, w, http.StatusBadRequest)
		return
	}

	id_decodified, err := helper.Decrypt(id_codified)
	if err != nil {
		helper.Response(helper.Response_struct{Error: "error to decodify"}, w, http.StatusBadRequest)
		return
	}

	ids := map[string]string{
		"Codificado:":   id_codified,
		"Decodificado:": id_decodified,
	}

	helper.Response(helper.Response_struct{Data: ids}, w, http.StatusOK)
}
