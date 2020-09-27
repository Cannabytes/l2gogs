package serverpackets

import "l2gogameserver/packets"

func NewExBrExtraUserInfo() []byte {

	buffer := new(packets.Buffer)

	buffer.WriteSingleByte(0xFE)
	buffer.WriteH(0xDA)
	buffer.WriteD(1)
	buffer.WriteD(0)
	buffer.WriteD(0)
	return buffer.Bytes()
}
