package sse

import (
    "log"
)

type Channel struct {
    lastEventId,
    name string
    clients map[*Client]bool
}

func NewChannel(name string) *Channel {
    return &Channel{
        "",
        name,
        make(map[*Client]bool),
    }
}

// SendMessage broadcast a message to all clients in a channel.
func (c *Channel) SendMessage(message *Message) {
    c.lastEventId = message.id

    for c, open := range c.clients {
        if open {
            c.send <- message
        }
    }

    log.Printf("go-sse: message sent to channel '%s'.", c.name)
}

func (c *Channel) Close() {
    // Kick all clients of this channel.
    for client, _ := range c.clients {
        c.removeClient(client)
    }

    log.Printf("go-sse: channel '%s' closed.", c.name)
}

func (c *Channel) ClientCount() int {
    return len(c.clients)
}

func (c *Channel) LastEventId() string {
    return c.lastEventId
}

func (c *Channel) addClient(client *Client) {
    c.clients[client] = true
    log.Printf("go-sse: new client connected to channel '%s'.", c.name)
}

func (c *Channel) removeClient(client *Client) {
    c.clients[client] = false
    close(client.send)
    delete(c.clients, client)
    log.Printf("go-sse: client disconnected from channel '%s'.", c.name)
}
