package helpers

import (
	"bufio"
	"net"
	"strings"
)

//To format the users output
// http://www.isthe.com/chongo/tech/comp/ansi_escapes.html
var (
	reset = []byte{27, 91, 13}
	home  = []byte("\x1b[H")  //{27, 91, 72}
	clear = []byte("\x1b[2J") //{27, 91, 50, 74}
	up    = []byte("\x1b[1A") //{27, 91, 1, 65}
	down  = []byte("\x1b[1B") //{27, 91, 1, 66}
	back  = []byte("\x1b[1D") //{27, 91, 1, 68}
	red   = []byte("\x1b[33m")
	green = []byte("\x1b[32m")
	blue  = []byte("\x1b[34m")
	white = []byte("\x1b[37m")
)

func SendText(conn net.Conn, text []byte) error {
	_, err := conn.Write(text)
	return err
}

func SendClear(conn net.Conn) error {
	return SendText(conn, clear)
}

func SendReset(conn net.Conn) error {
	return SendText(conn, reset)
}

func SendHome(conn net.Conn) error {
	return SendText(conn, home)
}

func SendUp(conn net.Conn) error {
	return SendText(conn, up)
}

func SendDown(conn net.Conn) error {
	return SendText(conn, down)
}

func SendBack(conn net.Conn) error {
	return SendText(conn, back)
}

func SendColor(conn net.Conn, text []byte, color []byte) error {
	if err := SendText(conn, color); err != nil {
		return err
	}
	if err := SendText(conn, text); err != nil {
		return err
	}
	return SendText(conn, white)
}

func SendRed(conn net.Conn, text []byte) error {
	return SendColor(conn, text, red)
}

func SendGreen(conn net.Conn, text []byte) error {
	return SendColor(conn, text, green)
}

func SendBlue(conn net.Conn, text []byte) error {
	return SendColor(conn, text, blue)
}

func ReadString(conn net.Conn) (string, error) {
	io := bufio.NewReader(conn)
	line, err := io.ReadString('\n')
	if err != nil {
		return "", err
	}
	if err := SendDown(conn); err != nil {
		return "", err
	}

	remove := []string{
		"\n", "\r", "\b",
	}

	for _, r := range remove {
		line = strings.Replace(line, r, "", -1)
	}
	return line, nil
}

func Processing(line string) string {
	remove := []string{
		"\n", "\r", "\b",
	}
	for _, r := range remove {
		line = strings.Replace(line, r, "", -1)
	}
	words := strings.Split(line, " ")
	line = ""
	for _, r := range words {
		if r != "" {
			line = line + strings.ToUpper(r[:1]) + strings.ToLower(r[1:]) + " "
		}
	}
	if line != "" {
		line = line[:len(line)-1]
	}
	return line
}
