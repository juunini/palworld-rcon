package palworldrcon

import (
	"fmt"
	"strings"

	"github.com/gorcon/rcon"
)

type Client struct {
	Host          string
	Port          uint
	AdminPassword string

	connection *rcon.Conn
}

func (c *Client) Connect() (err error) {
	c.connection, err = rcon.Dial(fmt.Sprintf("%s:%d", c.Host, c.Port), c.AdminPassword)
	return
}

func (c *Client) Disconnect() error {
	return c.connection.Close()
}

/*
Shutdown the server.

If <seconds> is specified, the server will shut down after the specified time has elapsed.

The server participant will be notified of what you have entered in <message>.
*/
func (c *Client) Shutdown(seconds uint, message string) (string, error) {
	return c.connection.Execute(fmt.Sprintf("Shutdown %d %s", seconds, message))
}

// Force stop the server.
func (c *Client) DoExit() (string, error) {
	return c.connection.Execute("DoExit")
}

// Send message to all player in the server.
func (c *Client) Broadcast(message string) (string, error) {
	return c.connection.Execute(fmt.Sprintf("Broadcast %s", message))
}

// Kick player by <steamID> from the server.
func (c *Client) KickPlayer(steamID string) (string, error) {
	return c.connection.Execute(fmt.Sprintf("KickPlayer %s", steamID))
}

// Ban player by <steamID> from the server.
func (c *Client) BanPlayer(steamID string) (string, error) {
	return c.connection.Execute(fmt.Sprintf("BanPlayer %s", steamID))
}

/*
Show information on all connected players.

ISSUE: when call "ShowPlayers" command, receive correct response with i/o timeout together.
*/
func (c *Client) ShowPlayers() (response string, err error) {
	response, err = c.connection.Execute("ShowPlayers")

	if strings.Contains(response, "name,playeruid,steamid") {
		err = nil
	}

	return
}

// Show server information.
func (c *Client) Info() (string, error) {
	return c.connection.Execute("Info")
}

// Save the world data.
func (c *Client) Save() (string, error) {
	return c.connection.Execute("Save")
}
