package serverpackets

import (
	"l2gogameserver/gameserver/clientpackets"
	"l2gogameserver/gameserver/models"
)

func NewMoveToLocation(location *clientpackets.Location, client *models.Client) {

	client.Buffer.WriteH(0) //reserve for lenght
	client.Buffer.WriteSingleByte(0x2f)

	client.Buffer.WriteD(1)

	client.Buffer.WriteD(location.TargetX)
	client.Buffer.WriteD(location.TargetY)
	client.Buffer.WriteD(location.TargetZ)

	client.Buffer.WriteD(location.OriginX)
	client.Buffer.WriteD(location.OriginY)
	client.Buffer.WriteD(location.OriginZ)
	//
	//buffer := new(packets.Buffer)
	//buffer.WriteSingleByte(0x2f)
	//
	//buffer.WriteD(1)
	//
	//buffer.WriteD(location.TargetX)
	//buffer.WriteD(location.TargetY)
	//buffer.WriteD(location.TargetZ)
	//
	//buffer.WriteD(location.OriginX)
	//buffer.WriteD(location.OriginY)
	//buffer.WriteD(location.OriginZ)

	//	return buffer.Bytes()
}
