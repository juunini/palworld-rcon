package palworldrcon

import (
	"fmt"
	"strings"
	"time"
)

type Client struct {
	Host          string
	Port          uint
	AdminPassword string
	// Default is 15 seconds.
	Timeout time.Duration

	connection *gameRCON
}

func (c *Client) Connect() error {
	if c.Timeout == 0 {
		c.Timeout = 15 * time.Second
	}

	c.connection = newGameRCON(c.Host, int(c.Port), c.AdminPassword, c.Timeout)
	return c.connection.connect()
}

func (c *Client) Disconnect() error {
	return c.connection.close()
}

/*
Shutdown the server.

If <seconds> is specified, the server will shut down after the specified time has elapsed.

The server participant will be notified of what you have entered in <message>.
*/
func (c *Client) Shutdown(seconds uint, message string) (string, error) {
	return c.connection.sendCommand(fmt.Sprintf("Shutdown %d %s", seconds, message))
}

// Force stop the server.
func (c *Client) DoExit() (string, error) {
	return c.connection.sendCommand("DoExit")
}

// Send message to all player in the server.
func (c *Client) Broadcast(message string) (string, error) {
	return c.connection.sendCommand(fmt.Sprintf("Broadcast %s", message))
}

// Kick player by <steamID> from the server.
func (c *Client) KickPlayer(steamID string) (string, error) {
	return c.connection.sendCommand(fmt.Sprintf("KickPlayer %s", steamID))
}

// Ban player by <steamID> from the server.
func (c *Client) BanPlayer(steamID string) (string, error) {
	return c.connection.sendCommand(fmt.Sprintf("BanPlayer %s", steamID))
}

/*
Show information on all connected players.

ISSUE: when call "ShowPlayers" command, receive correct response with i/o timeout together.
*/
func (c *Client) ShowPlayers() (response string, err error) {
	response, err = c.connection.sendCommand("ShowPlayers")

	if strings.Contains(response, "name,playeruid,steamid") {
		err = nil
	}

	return
}

// Show server information.
func (c *Client) Info() (string, error) {
	return c.connection.sendCommand("Info")
}

// Save the world data.
func (c *Client) Save() (string, error) {
	return c.connection.sendCommand("Save")
}
