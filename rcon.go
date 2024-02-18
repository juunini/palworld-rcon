// This code copy from https://github.com/dkoz/gamercon-async/blob/main/gamercon_async/gamercon_async.py
// and convert to golang by ChatGPT 3.5
// and then fix some code by @juunini
package palworldrcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	AUTH             = 3
	AUTH_RESPONSE    = 2
	EXEC_COMMAND     = 2
	COMMAND_RESPONSE = 0
	MIN_INT_32       = int32(-2147483648)
	MAX_INT_32       = int32(2147483647)
)

type clientError struct {
	message string
}

func (e *clientError) Error() string {
	return e.message
}

type invalidPassword struct {
}

func (e *invalidPassword) Error() string {
	return "Invalid password"
}

type connectionError struct {
	message string
}

func (e *connectionError) Error() string {
	return e.message
}

type commandExecutionError struct {
	message string
}

func (e *commandExecutionError) Error() string {
	return e.message
}

type emptyResponse struct {
}

func (e *emptyResponse) Error() string {
	return "Empty response"
}

type littleEndianSignedInt32 int32

func newLittleEndianSignedInt32(value int32) littleEndianSignedInt32 {
	if value < MIN_INT_32 || value > MAX_INT_32 {
		panic("Signed int32 out of bounds")
	}
	return littleEndianSignedInt32(value)
}

func (i littleEndianSignedInt32) toBytes() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func littleEndianSignedInt32FromBytes(b []byte) littleEndianSignedInt32 {
	buf := bytes.NewReader(b)
	var value int32
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		panic(err)
	}
	return littleEndianSignedInt32(value)
}

type packet struct {
	id         littleEndianSignedInt32
	packetType int32
	payload    []byte
	terminator []byte
}

func newPacket(id littleEndianSignedInt32, packetType int32, payload []byte, terminator []byte) packet {
	return packet{
		id:         id,
		packetType: packetType,
		payload:    payload,
		terminator: terminator,
	}
}

func (p *packet) toBytes() []byte {
	payload := append(append(append(p.id.toBytes(), littleEndianSignedInt32(p.packetType).toBytes()...), p.payload...), p.terminator...)
	size := newLittleEndianSignedInt32(int32(len(payload)))
	return append(size.toBytes(), payload...)
}

func makeCommandPacket(command string) packet {
	id := newLittleEndianSignedInt32(randInt(0, MAX_INT_32))
	return newPacket(id, EXEC_COMMAND, []byte(command), []byte{0x00, 0x00})
}

func makeLoginPacket(password string) packet {
	id := newLittleEndianSignedInt32(randInt(0, MAX_INT_32))
	return newPacket(id, AUTH, []byte(password), []byte{0x00, 0x00})
}

func randInt(min int32, max int32) int32 {
	return min + rand.Int31n(max-min)
}

type gameRCON struct {
	host     string
	port     int
	password string
	timeout  time.Duration
	auth     bool
	conn     net.Conn
}

func newGameRCON(host string, port int, password string, timeout time.Duration) *gameRCON {
	return &gameRCON{
		host:     host,
		port:     port,
		password: password,
		timeout:  timeout,
		auth:     false,
		conn:     nil,
	}
}

func (rcon *gameRCON) connect() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", rcon.host, rcon.port), rcon.timeout)
	if err != nil {
		return &connectionError{message: fmt.Sprintf("Error connecting to %s:%d - %s", rcon.host, rcon.port, err)}
	}
	rcon.conn = conn

	err = rcon.authenticate()
	if err != nil {
		return err
	}

	return nil
}

func (rcon *gameRCON) close() error {
	if rcon.conn != nil {
		err := rcon.conn.Close()
		if err != nil {
			return &connectionError{message: fmt.Sprintf("Error disconnecting from %s:%d - %s", rcon.host, rcon.port, err)}
		}
	}
	return nil
}

func (rcon *gameRCON) authenticate() error {
	loginPacket := makeLoginPacket(rcon.password)
	err := rcon.sendPacket(loginPacket)
	if err != nil {
		return err
	}

	responsePacket, err := rcon.readPacket()
	if err != nil {
		return err
	}

	if responsePacket.id == -1 {
		return &invalidPassword{}
	}

	rcon.auth = true
	return nil
}

func (rcon *gameRCON) sendPacket(p packet) error {
	if !rcon.auth && p.packetType != AUTH {
		return &clientError{message: "Client not authenticated."}
	}
	if rcon.conn == nil {
		return &clientError{message: "Not connected."}
	}

	_, err := rcon.conn.Write(p.toBytes())
	if err != nil {
		return &clientError{message: fmt.Sprintf("Error sending packet: %s", err)}
	}

	return nil
}

func (rcon *gameRCON) readPacket() (packet, error) {
	sizeData := make([]byte, 4)
	_, err := rcon.conn.Read(sizeData)
	if err != nil {
		return packet{}, &emptyResponse{}
	}

	size := littleEndianSignedInt32FromBytes(sizeData)
	packetData := make([]byte, size)
	_, err = rcon.conn.Read(packetData)
	if err != nil {
		return packet{}, &emptyResponse{}
	}

	id := littleEndianSignedInt32FromBytes(packetData[:4])
	packetType := littleEndianSignedInt32FromBytes(packetData[4:8])
	payload := packetData[8 : len(packetData)-2]

	return newPacket(id, int32(packetType), payload, []byte{0x00, 0x00}), nil
}

func (rcon *gameRCON) sendCommand(cmd string) (string, error) {
	if !rcon.auth {
		return "", &clientError{message: "Not authenticated with RCON server."}
	}

	commandPacket := makeCommandPacket(cmd)
	err := rcon.sendPacket(commandPacket)
	if err != nil {
		return "", err
	}

	responsePacket, err := rcon.readPacket()
	if err != nil {
		return "", err
	}

	if responsePacket.id == -1 {
		return "", &invalidPassword{}
	}
	if responsePacket.packetType != COMMAND_RESPONSE {
		return "", &commandExecutionError{message: "Unexpected response type."}
	}

	return string(responsePacket.payload), nil
}
