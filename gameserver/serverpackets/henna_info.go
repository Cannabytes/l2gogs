package serverpackets

import "l2gogameserver/packets"

func NewHennaInfo() []byte {

	buffer := new(packets.Buffer)

	buffer.WriteSingleByte(0xE5)

	buffer.WriteSingleByte(0) // equip INT
	buffer.WriteSingleByte(0) // equip STR
	buffer.WriteSingleByte(0) //equip CON
	buffer.WriteSingleByte(0) // equip MEN
	buffer.WriteSingleByte(0) //equip DEX
	buffer.WriteSingleByte(0) //equip WIT
	buffer.WriteD(3)          //slots
	buffer.WriteD(0)          //Size

	return buffer.Bytes()
}
