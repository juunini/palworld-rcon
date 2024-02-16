package palworldrcon_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gorcon/rcon"
	palworldrcon "github.com/juunini/palworld-rcon"
)

const HOST = "127.0.0.1"
const PORT = 25575
const ADMIN_PASSWORD = "password" // change here admin password

const EXPECTED_SAVE_RESPONSE = "Complete Save"
const EXPECTED_BROADCAST_RESPONSE = "Broadcasted: "
const EXPECTED_SHOW_PLAYERS_RESPONSE = "name,playeruid,steamid"
const EXPECTED_BAN_PLAYER_RESPONSE = "Baned: " + STEAMID
const EXPECTED_KICKED_PLAYER_RESPONSE = "Kicked: " + STEAMID
const EXPECTED_SHUTDOWN_RESPONSE = "The server will shut down in 100 seconds. Please prepare to exit the game."
const EXPECTED_DO_EXIT_RESPONSE = "Shutdown..."
const STEAMID = "00000000000000000"

func setup() {
	connection, err := rcon.Dial(fmt.Sprintf("%s:%d", HOST, PORT), ADMIN_PASSWORD)
	if err != nil {
		fmt.Println("\033[93m" + fmt.Sprintf("address: %s:%d\npassword: %s", HOST, PORT, ADMIN_PASSWORD) + "\033[0m")
		fmt.Println("\033[91m" + "Can't connect RCON" + "\033[0m")
		os.Exit(1)
	}

	connection.Close()
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()
	os.Exit(code)
}

func TestCommands(t *testing.T) {
	client := palworldrcon.Client{
		Host:          HOST,
		Port:          PORT,
		AdminPassword: ADMIN_PASSWORD,
	}

	client.Connect()
	defer client.Disconnect()

	defer testDoExit(t, client)

	testInfo(t, client)
	testSave(t, client)
	testBroadcast(t, client)
	testShowPlayers(t, client)
	testKickPlayer(t, client)
	testBanPlayer(t, client)
	testShutdown(t, client)
}

func testInfo(t *testing.T, client palworldrcon.Client) {
	response, err := client.Info()
	if err != nil {
		t.Fatal(err)
	}

	checkUnknownCommand(t, "Info", response)

	t.Logf(`Palworld Server Info: "%s"`, strings.TrimSpace(response))
}

func testSave(t *testing.T, client palworldrcon.Client) {
	response, err := client.Save()
	if err != nil {
		t.Fatal(err)
	}

	checkUnknownCommand(t, "Save", response)

	if strings.TrimSpace(response) != EXPECTED_SAVE_RESPONSE {
		t.Fatalf(`Save response is not "%s", Actual: "%s"`, EXPECTED_SAVE_RESPONSE, response)
	}
}

func testBroadcast(t *testing.T, client palworldrcon.Client) {
	response, err := client.Broadcast("hi")
	if err != nil {
		t.Fatal(err)
	}

	checkUnknownCommand(t, "Broadcast", response)

	if !strings.Contains(response, EXPECTED_BROADCAST_RESPONSE) {
		t.Fatalf(`Broadcast response is not contains "%s", Actual: "%s"`, EXPECTED_BROADCAST_RESPONSE, response)
	}
}

func testKickPlayer(t *testing.T, client palworldrcon.Client) {
	response, err := client.KickPlayer(STEAMID)
	if err != nil {
		t.Fatal(err)
	}

	checkUnknownCommand(t, "KickPlayer", response)

	if strings.TrimSpace(response) != EXPECTED_KICKED_PLAYER_RESPONSE {
		t.Fatalf(`KickPlayer response is not "%s", Actual: "%s"`, EXPECTED_KICKED_PLAYER_RESPONSE, response)
	}
}
func testBanPlayer(t *testing.T, client palworldrcon.Client) {
	response, err := client.BanPlayer(STEAMID)
	if err != nil {
		t.Fatal(err)
	}

	checkUnknownCommand(t, "BanPlayer", response)

	if strings.TrimSpace(response) != EXPECTED_BAN_PLAYER_RESPONSE {
		t.Fatalf(`BanPlayer response is not "%s", Actual: "%s"`, EXPECTED_BAN_PLAYER_RESPONSE, response)
	}
}

func testShowPlayers(t *testing.T, client palworldrcon.Client) {
	response, err := client.ShowPlayers()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(response)

	checkUnknownCommand(t, "ShowPlayers", response)

	if !strings.Contains(response, EXPECTED_SHOW_PLAYERS_RESPONSE) {
		t.Fatalf(`ShowPlayers response is not contains "%s", Actual: "%s"`, EXPECTED_SHOW_PLAYERS_RESPONSE, response)
	}
}

func testShutdown(t *testing.T, client palworldrcon.Client) {
	response, err := client.Shutdown(100, "shutdown after 100 seconds later")
	if err != nil {
		t.Fatal(err)
	}

	checkUnknownCommand(t, "Shutdown", response)

	if strings.TrimSpace(response) != EXPECTED_SHUTDOWN_RESPONSE {
		t.Fatalf(`Save response is not "%s", Actual: %s`, EXPECTED_SHUTDOWN_RESPONSE, response)
	}
}

func testDoExit(t *testing.T, client palworldrcon.Client) {
	response, err := client.DoExit()
	if err != nil {
		t.Fatal(err)
	}

	checkUnknownCommand(t, "DoExit", response)

	if strings.TrimSpace(response) != EXPECTED_DO_EXIT_RESPONSE {
		t.Fatalf(`Save response is not "%s", Actual: %s`, EXPECTED_DO_EXIT_RESPONSE, response)
	}
}

func checkUnknownCommand(t *testing.T, insertedCommand string, response string) {
	if strings.TrimSpace(response) == "Unknown command" {
		t.Fatalf(`Inserted command "%s" is "Unknown command"`, insertedCommand)
	}
}
