package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var clients []*Client

const messageBufferSize = 256

type Client struct {
	conn *websocket.Conn
	send chan string
}

func newClient(conn *websocket.Conn) {
	c := &Client{
		conn: conn,
		send: make(chan string, messageBufferSize),
	}
	clients = append(clients, c)

	c.send <- string("connected at: " + time.Now().Format(time.UnixDate))
	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) Close() {
	for i, client := range clients {
		if client == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	close(c.send)
	c.conn.Close()
	log.Printf("close connection, addr: %s", c.conn.RemoteAddr())
}

func (c *Client) readLoop() {
	for {
		m, err := c.read()
		if err != nil {
			log.Println("read message error: ", err)
			break
		}
		broadcast(m)
	}
	c.Close()
}

func (c *Client) writeLoop() {
	for msg := range c.send {
		c.write(msg)
	}
}

func broadcast(m string) {
	for _, client := range clients {
		client.send <- m
	}
}

func (c *Client) read() (string, error) {
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return "", err
	}
	log.Printf("read from websocket: %s", msg)

	return string(msg), nil
}

func (c *Client) write(m string) error {
	log.Printf("write to websocket: %s", m)
	return c.conn.WriteMessage(websocket.TextMessage, []byte(m))
}
