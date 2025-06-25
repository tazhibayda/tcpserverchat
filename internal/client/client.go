package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
	"tcpserverchat/internal/dispatcher"
)

func Handle(conn net.Conn, d *dispatcher.Dispatcher) {
	defer conn.Close()

	id := uuid.NewString()
	msgChan := make(chan string, 5)
	d.Join(dispatcher.Client{ID: id, MsgChan: msgChan})
	defer d.Leave(dispatcher.Client{ID: id, MsgChan: msgChan})

	go func() {
		for msg := range msgChan {
			fmt.Fprintln(conn, msg)
		}
	}()
	reader := bufio.NewReader(conn)
	// после bufio.NewReader
	var nickname string = id // по умолчанию
	for {
		line, err := reader.ReadString('\n')
		if err != nil { /* ... */
			return
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/nick ") {
			newNick := strings.TrimSpace(strings.TrimPrefix(line, "/nick "))
			if newNick != "" {
				old := nickname
				nickname = newNick
				d.Broadcast(fmt.Sprintf("%s switched nickname to %s", old, nickname))
			}
		} else {
			d.Broadcast(fmt.Sprintf("[%s] %s", nickname, line))
		}
	}

	//for {
	//	line, err := reader.ReadString('\n')
	//	if err != nil {
	//		if err != io.EOF {
	//			log.Println("read:", err)
	//		}
	//		return
	//	}
	//	d.Broadcast(fmt.Sprintf("%s: %s", id, line))
	//}
}
