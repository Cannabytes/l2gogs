package models

import (
	"errors"
	"l2gogameserver/gameserver/crypt"
	"l2gogameserver/packets"
	"log"
	"net"
)

type Client struct {
	Socket          net.Conn
	ScrambleModulus []byte
	Buffer          *packets.Buffer
	Reader          *packets.Rreader
}

func NewClient() *Client {
	buff := new(packets.Buffer)
	read := packets.NewRReader()
	return &Client{Buffer: buff, Reader: read}
}

func (c *Client) Send(data []byte, need bool) error {
	if need {
		data = crypt.Encrypt(data)
	}
	// Calculate the packet length
	length := int16(len(data) + 2)
	// Put everything together
	buffer := packets.NewBuffer()
	buffer.WriteH(length)
	_, err := buffer.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.Socket.Write(buffer.Bytes())
	if err != nil {
		return errors.New("The packet couldn't be sent.")
	}
	c.Reader.B.Reset()
	return nil
}

func (c *Client) SimpleSend(data []byte) error {
	length := int16(len(data))
	data[0], data[1] = uint8(length&0xff), uint8(length>>8)

	data = crypt.SimpleEncrypt(data)

	_, err := c.Socket.Write(data)
	if err != nil {
		return errors.New("The packet couldn't be sent.")
	}
	c.Reader.B.Reset()
	return nil
}

func (c *Client) Receive() (opcode byte, data []byte, e error) {
	// Read the first two bytes to define the packet size
	header := make([]byte, 2)
	n, err := c.Socket.Read(header)
	//fmt.Println(n)
	if n != 2 || err != nil {
		return 0x00, nil, errors.New("12An error occured while reading the packet header.")
	}

	// Calculate the packet size
	size := 0
	size += int(header[0])
	size += int(header[1]) * 256

	// Allocate the appropriate size for our data (size - 2 bytes used for the length
	data = make([]byte, size-2)

	// Read the encrypted part of the packet
	n, err = c.Socket.Read(data)
	if n != size-2 || err != nil {
		return 0x00, nil, errors.New("An error occured while reading the packet data.")
	}

	// Print the raw packet
	//fmt.Printf("header packet : %X\n  Raw: %X\n", header, data)
	data = crypt.Decrypt(data)
	// Extract the op code
	opcode = data[0]
	c.Reader.AddB(data[1:])
	data = data[1:]
	e = nil
	return
}
