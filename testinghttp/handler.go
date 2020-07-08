package testinghttp

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
		HandleResponse(w, 500, "read post body failed")
		return
	}
	defer r.Body.Close()

	var input LoginForm
	if err = json.Unmarshal(data, &input); err != nil {
		HandleResponse(w, 400, "input invalid format")
		return
	}

	if input.Code != "a@example.com" || input.Password != "password" {
		HandleResponse(w, 400, "invalid code or password")
		return
	}

	HandleResponse(w, 200, "ok")
}

func HandleResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
