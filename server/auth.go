package wssv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func register(username, password string) error {
	return nil
}

func (sv *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
			resp, err := json.Marshal(LoginResp{
				Error: "Login failed",
			})
			if err != nil {
				log.Println(err)
				return
			}
			w.Write(resp)
		}
	}()
	var req LoginReq
	reqB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(reqB, &req)
	if err != nil {
		return
	}
	token, err := login(req.Username, req.Password)
	if err != nil {
		return
	}
	resp, err := json.Marshal(LoginResp{
		Token: token,
	})
	if err != nil {
		return
	}
	w.Write(resp)
}

func login(username, password string) (string, error) {
	if username == "test" && password == "test" {
		return "testtoken", nil
	}
	if username == "user" && password == "user" {
		return "usertoken", nil
	}
	return "", fmt.Errorf("Login failed")
}

func validate(token string) (string, error) {
	if token == "testtoken" {
		return "test", nil
	} else if token == "usertoken" {
		return "user", nil
	}
	return "", fmt.Errorf("Failed to validate token")
}
