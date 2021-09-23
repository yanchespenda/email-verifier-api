package handler

import (
	"fmt"
	"io"
	"net/http"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/Jeffail/gabs/v2"
	"github.com/gorilla/mux"
)

var (
	verifier = emailverifier.NewVerifier()
)

// GetRouter returns the router for the API
func GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", Handler).Methods(http.MethodGet)
	return r
}

func Handler(w http.ResponseWriter, r *http.Request) {

	jsonObj := gabs.New()

	email := r.URL.Query().Get("email")
	if len(email) == 0 {
		fmt.Println("Missing email parameter")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		jsonObj.Set("Email parameter missing", "message")
		io.WriteString(w, jsonObj.String())
		return
	}

	ret, err := verifier.Verify(email)
	if err != nil {
		fmt.Println("verify email address failed, error is: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		jsonObj.Set("Email verify failed", "message")
		jsonObj.Set(err, "err")
		io.WriteString(w, jsonObj.String())
		return
	}
	if !ret.Syntax.Valid {
		fmt.Println("email address syntax is invalid")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		jsonObj.Set("Email address invalid", "message")
		io.WriteString(w, jsonObj.String())
		return
	}

	// fmt.Fprintf(w, "<h1>Hello from Go!</h1>")

	fmt.Println("Verify email address successfull, result: ", ret)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonObj.Set("Email valid", "message")
	jsonObj.Set(ret, "result")
	io.WriteString(w, jsonObj.String())

}
