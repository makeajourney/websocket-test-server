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
	send chan *Message
}

type Message struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func newClient(conn *websocket.Conn) {
	c := &Client{
		conn: conn,
		send: make(chan *Message, messageBufferSize),
	}
	clients = append(clients, c)

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

func broadcast(m *Message) {
	for _, client := range clients {
		client.send <- m
	}
}

func (c *Client) read() (*Message, error) {
	var msg *Message
	if err := c.conn.ReadJSON(&msg); err != nil {
		return nil, err
	}
	msg.Content = "server : " + msg.Content
	msg.CreatedAt = time.Now()
	log.Println("read from websocket:", msg)

	return msg, nil
}

func (c *Client) write(m *Message) error {
	log.Println("write to websocket:", m)
	return c.conn.WriteJSON(m)
}
