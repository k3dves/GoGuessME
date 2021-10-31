package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	conn *net.Conn
}

func ConnectToServer() (*Client, error) {
	addr := flag.String("addr", "localhost:6060", "<server ip>:<port>")
	flag.Parse()
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Printf("Error conneting to the server. Error: %s \n ", err)
		return &Client{}, err
	}
	client := &Client{}
	client.conn = &conn
	d("Connected to sever")

	return client, nil
}

func (c *Client) Close() {
	(*c.conn).Close()
}

func (c *Client) send(msg string) error {
	d("start send msg")
	_, err := fmt.Fprint(*c.conn, msg)
	d("end send msg")
	//fmt.Println(err)
	return err
}

func (c *Client) read() {
	for {
		d("start read")
		msg, err := bufio.NewReader((*c.conn)).ReadString('\n')

		d("Got msg " + msg)
		if err != nil {
			fmt.Printf("Error reading msg from server. Error %s", err)
			break
		}

		if err == io.EOF {
			fmt.Printf("Server Quit")
			c.Close()
			break
		}
		fmt.Print(msg)
	}
}

func main() {
	c, err := ConnectToServer()
	if err != nil {
		fmt.Printf("Error creating a connection. Error %s", err)

	}

	go c.read()

	reader := bufio.NewReader(os.Stdin)
	for {
		d("for first")
		msg, err := reader.ReadString('\n')
		d("main: msg " + msg)
		if err != nil || err == io.EOF {
			fmt.Printf("Connection closed by client. Error: %s", err)
			break
		}
		d("Send msg before")
		c.send(msg)
		d("Send msg after")
	}
}

func d(m string) {
	//fmt.Println(m)
}
