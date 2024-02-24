# <div align="center">Palworld RCON</div>

<div align="center">
    <img src="https://github.com/juunini/palworld-rcon/assets/41536271/8414cd69-68f4-45bc-a052-9c4afa652582" alt="Palworld Icon" />
</div>

<div align="center">
    <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
    <a href="https://goreportcard.com/report/github.com/juunini/palworld-rcon" target="_blank">
        <img src="https://goreportcard.com/badge/github.com/juunini/palworld-rcon" alt="Go Report Card" />
    </a>
    <a href="https://godoc.org/github.com/juunini/palworld-rcon" target="_blank">
        <img src="https://img.shields.io/badge/godoc-reference-blue.svg" alt="Go Report Card" />
    </a>
</div>

## Install

```
go get github.com/juunini/palworld-rcon
```

## Usage

```go
package main

import (
    "fmt"
    "time"

    palworldrcon "github.com/juunini/palworld-rcon"
)

func main() {
    client, err := palworldrcon.Connect("127.0.0.1", 25575, "your admin password", 15 * time.Second)
    if err != nil {
        panic(err)
    }

    defer client.Disconnect()

    response, err := client.Info()
    if err != nil {
        panic(err)
    }

    fmt.Println(response) // "Welcome to Pal Server[v0.1.4.1] Default Palworld Server"
}
```

## Client Methods

see: https://tech.palworldgame.com/settings-and-operation/commands

| Method name | Properties | Description |
| - | - | - |
| Connect | | Connect to RCON Server |
| Disconnect | | Close Connection |
| Shutdown | seconds, message | Shutdown the server. If \<seconds\> is specified, the server will shut down after the specified time has elapsed. The server participant will be notified of what you have entered in \<message\>. |
| DoExit | | Force stop the server. |
| Broadcast | message | Send message to all player in the server. |
| KickPlayer | steamID | Kick player by \<steamID\> from the server. |
| BanPlayer | steamID | Ban player by \<steamID\> from the server. |
| ShowPlayers | | Show information on all connected players. |
| Info | | Show server information. |
| Save | | Save the world data. |

## License

[MIT License](./LICENSE)
