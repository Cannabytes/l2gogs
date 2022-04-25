package serverpackets

import (
	"l2gogameserver/gameserver/models"
	"l2gogameserver/packets"
)

func AbnormalStatusUpdate(userBuff []*models.BuffUser) []byte {
	buffer := packets.Get()
	defer packets.Put(buffer)
	buffer.WriteSingleByte(0x85)
	buffer.WriteH(40)

	for _, buff := range userBuff {
		buffer.WriteD(int32(buff.Id))
		buffer.WriteH(int16(buff.Level))
		buffer.WriteD(int32(buff.Second))
	}
	return buffer.Bytes()
}
