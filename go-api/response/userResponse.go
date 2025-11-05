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
		panic(err)
	}

	/* users, exist := res[1]
	if !exist {
		panic("This user not exists")
	} */

	// setting the header to the response
	/* w.Header().Set("Content-Type", "application/json")

	// setting the response
	json.NewEncoder(w).Encode(res) */

	helper.Response(res, w)

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	id_codified, err := helper.Encrypt(id)
	if err != nil {
		panic(err)
	}

	id_decodified, err := helper.Decrypt(id_codified)
	if err != nil {
		panic(err)
	}

	helper.Response(id_codified+" Decodificado: "+id_decodified, w)
}
