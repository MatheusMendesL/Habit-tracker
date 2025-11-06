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
		helper.Response(map[string]string{"Error": "Cannot get the users with the function"}, w, http.StatusBadRequest)
		return
	}

	helper.Response(res, w, http.StatusOK)

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	id_codified, err := helper.Encrypt(id)
	if err != nil {
		helper.Response(map[string]string{"error": "error to codify"}, w, http.StatusBadRequest)
		return
	}

	id_decodified, err := helper.Decrypt(id_codified)
	if err != nil {
		helper.Response(map[string]string{"error": "error to decodify"}, w, http.StatusBadRequest)
		return
	}

	ids := map[string]string{
		"Codificado:":   id_codified,
		"Decodificado:": id_decodified,
	}

	helper.Response(ids, w, http.StatusOK)
}
