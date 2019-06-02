package wssv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	wsDialer *websocket.Dialer
	host     string
	conn     *websocket.Conn
	closed   bool
	token    string
}

func NewClient(url, user, password string) (*Client, error) {
	token, err := clientLogin(url, user, password)
	if err != nil {
		return nil, err
	}
	return &Client{
		wsDialer: websocket.DefaultDialer,
		host:     url,
		token:    token,
	}, nil
}

func (wsc *Client) Connect() error {
	if wsc.token == "" {
		return fmt.Errorf("Token field is empty (you must login first)")
	}
	con, _, err := wsc.wsDialer.Dial(fmt.Sprintf("ws://%s", wsc.host), http.Header{"Authorization": []string{wsc.token}})
	if err != nil {
		return err
	}
	wsc.conn = con
	return nil
}

func (wsc *Client) Send(msg string) error {
	err := wsc.conn.WriteJSON(MessageReq{
		Authorization: wsc.token,
		Message:       msg,
	})
	// err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (wsc *Client) Receive() (string, error) {
	var msgResp MessageResp
	err := wsc.conn.ReadJSON(&msgResp)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:\t%s", msgResp.Author, msgResp.Message), nil
}

func (wsc *Client) Close() error {
	wsc.closed = true
	return wsc.conn.Close()
}

func (wsc *Client) IsClosed() bool {
	return wsc.closed
}

func clientLogin(url, user, pw string) (string, error) {
	body, err := json.Marshal(LoginReq{
		Username: user,
		Password: pw,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/login", url), bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	var response LoginResp
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", fmt.Errorf(response.Error)
	}
	return response.Token, nil
}
