package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user
type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	//send is a buffered channel on which messages are sent to the user's browser via the web socket
	send chan *message
	// room is the room this client is chatting in
	room *room
	// userData holds information about the user
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
