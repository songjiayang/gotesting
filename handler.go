package gotesting

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type LoginForm struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleErr(w, 500, "read post body failed")
		return
	}
	defer r.Body.Close()

	var input LoginForm
	if err = json.Unmarshal(data, &input); err != nil {
		handleErr(w, 400, "input invalid format")
		return
	}

	if input.Code != "a@example.com" || input.Password != "password" {
		handleErr(w, 400, "invalid code or password")
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func handleErr(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
