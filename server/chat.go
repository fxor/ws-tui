package wssv

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ClientConn struct {
	socket *websocket.Conn
	id     string
}

func (s *Server) upgrade(w http.ResponseWriter, r *http.Request) {
	// TODO: what does this do
	s.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	username, err := validate(r.Header.Get("Authorization"))
	if err != nil {
		log.Println(err)
		return
	}
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	clientConn := &ClientConn{
		socket: ws,
		id:     username,
	}

	// TODO: handle multiple sessions for the same user from different clients
	c, ok := s.clientSockets[username]
	if ok {
		c.socket.Close()
	}
	s.clientSockets[username] = clientConn
	log.Printf("User %s connected from %s\n", username, r.RemoteAddr)
	err = ws.WriteJSON(MessageResp{
		Message: fmt.Sprintf("Hi %s!\n", username),
		Author:  "Server",
	})
	if err != nil {
		log.Println(err)
	}
	s.broadcaster(clientConn)
}

func (s *Server) broadcaster(conn *ClientConn) {
	for {
		// messageType, msg, err := conn.ReadMessage()
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		var msgReq MessageReq
		err := conn.socket.ReadJSON(&msgReq)
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				log.Printf("Socket from client %s closed", conn.id)
				conn.socket.Close()
				delete(s.clientSockets, conn.id)
				return
			}
			log.Println("Broadcast read error:", err)
			return
		}

		author, err := validate(msgReq.Authorization)
		if err != nil {
			err = conn.socket.WriteJSON(MessageResp{
				Author:  "Server",
				Message: "Authentication error",
			})
			if err != nil {
				log.Println(err)
				return
			}
			err = conn.socket.Close()
			if err != nil {
				log.Println(err)
				return
			}
		}

		for k, c := range s.clientSockets {
			// err = c.WriteMessage(messageType, msg)
			err = c.socket.WriteJSON(MessageResp{
				Author:  author,
				Message: msgReq.Message,
			})
			if err != nil {
				log.Println("Broadcast error:", err)
				c.socket.Close()
				delete(s.clientSockets, k)
				continue
			}
		}
	}
}
