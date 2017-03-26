package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Message struct {
	fromServer bool
	name       string
	connection *net.Conn
	text       string
}

func (message *Message) String() string {
	if message.fromServer {
		return message.text
	}

	if message.name != "" {
		return message.name + ":\t" + message.text
	} else if message.connection != nil {
		return (*message.connection).RemoteAddr().String() + "\t" + message.text
	} else {
		return message.text
	}
}

func receiveMessages(connection *net.Conn, messageChannel chan Message) {
	// receive messages and write them to messageChannel
	defer (*connection).Close()
	var msg Message
	msg.connection = connection
	inputChannel := startReader(bufio.NewReader(*connection))
	for {
		var gotData bool
		select {
		case msg.text, gotData = <-inputChannel:
			if !gotData {
				msg.text = "LOGGED OFF"
				msg.text = fmt.Sprintf("%v", &msg)
				msg.connection = nil
				msg.name = ""
				messageChannel <- msg
				return
			}
			namePrefix := "name:"
			if strings.HasPrefix(msg.text, namePrefix) {
				msg.name = strings.TrimSpace(strings.TrimPrefix(msg.text, namePrefix))
				continue
			}
			messageChannel <- msg
		}
	}
}

func forwardMessage(connection *net.Conn, message *Message) {
	//log.Printf("Sending '%v' to %v", message, (*connection).RemoteAddr())

	writer := bufio.NewWriter(*connection)
	writer.WriteString(strings.TrimSpace((*message).String()) + "\n")
	writer.Flush()
}

func startListener(port string, name string) chan *net.Conn {
	connectionChannel := make(chan *net.Conn)
	// start a listener which sends new connections on the returned channel
	listener, err := net.Listen("tcp", port)
	if err != nil {
		// handle error
		log.Fatalf("Failed to start server %v - %v", port, err)
	}

	go func() {
		var welcomeMessage Message
		welcomeMessage.text = "Welcome to " + name
		// handle new connections and pass them to the connection channel
		for {
			conn, err := listener.Accept()
			if err != nil {
				// handle error
			}
			forwardMessage(&conn, &welcomeMessage)
			connectionChannel <- &conn
		}
	}()
	return connectionChannel
}

func runServer(port string, name string) {
	messageChannel := make(chan Message)
	connectionChannel := startListener(port, name)
	var connections []*net.Conn

	for {
		select {
		case connection := <-connectionChannel:
			connections = append(connections, connection)
			go receiveMessages(connection, messageChannel)
		case message := <-messageChannel:
			for _, connection := range connections {
				forwardMessage(connection, &message)
			}
		}
	}
}

func startReader(reader *bufio.Reader) chan string {
	// reads from reader and puts text on inputChannel
	inputChannel := make(chan string)
	go func() {
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				// handle err
				close(inputChannel)
				return
			}
			inputChannel <- input
		}
	}()
	return inputChannel
}

func runClient(serverAddress string, name string) {
	connection, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Println("Failed to connect to %v - %v", serverAddress, err)
		return
	}

	defer connection.Close()
	writer := bufio.NewWriter(connection)
	if name != "" {
		writer.WriteString("name:" + name + "\n")
		writer.Flush()
	}
	receiveChannel := startReader(bufio.NewReader(connection))
	stdinReader := bufio.NewReader(os.Stdin)
	sendChannel := startReader(stdinReader)

	for {
		select {
		case message, gotMessage := <-receiveChannel:
			if !gotMessage {
				// server disconnected
				return
			}
			fmt.Printf(message)
		case toSend, gotData := <-sendChannel:
			if !gotData {
				// stdin closed?
				return
			}
			writer.WriteString(strings.TrimSpace(toSend) + "\n")
			writer.Flush()
		}
	}
}

func main() {
	var port = flag.String("port", "90", "Server port")
	portStr := ":" + *port
	var server = flag.String("server", "", "Server to connect to")
	var name = flag.String("name", "", "name")

	flag.Parse()
	if *server == "" {
		if *name == "" {
			*name = "ChatServer"
		}
		// spawn chat server
		runServer(portStr, *name)
	} else {
		// connect to chat server
		//log.Printf("Connecting to %v", *server)
		runClient(*server+portStr, *name)
	}

}
