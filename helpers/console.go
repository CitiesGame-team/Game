package helpers

import (
	"bufio"
	"net"
	"strings"
)

func SendText(conn net.Conn, text []byte) error {
	_, err := conn.Write(text)
	return err
}

func SendClear(conn net.Conn) error {
	return SendText(conn, []byte{27, 91, 50, 74})
}

func SendReset(conn net.Conn) error {
	return SendText(conn, []byte{27, 91, 13})
}

func SendHome(conn net.Conn) error {
	return SendText(conn, []byte{27, 91, 72})
}

func SendDown(conn net.Conn) error {
	return SendText(conn, []byte{27, 91, 1, 66})
}

func SendUp(conn net.Conn) error {
	return SendText(conn, []byte{27, 91, 65})
}

func SendColor(conn net.Conn, text []byte, color []byte) error {
	if err := SendText(conn, color); err != nil {
		return err
	}
	if err := SendText(conn, text); err != nil {
		return err
	}
	return SendText(conn, []byte("\x1b[0m"))
}

func SendRed(conn net.Conn, text []byte) error {
	return SendColor(conn, text, []byte("\x1b[33m"))
}

func SendGreen(conn net.Conn, text []byte) error {
	return SendColor(conn, text, []byte("\x1b[32m"))
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
		"\n", "\r",
	}

	for _, r := range remove {
		line = strings.Replace(line, r, "", -1)
	}
	return line, nil
}
