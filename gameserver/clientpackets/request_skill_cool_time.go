package clientpackets

import (
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
	"l2gogameserver/packets"
)

func RequestSkillCoolTime(client *models.Client, data []byte) []byte {
	buffer := packets.Get()
	defer packets.Put(buffer)

	pkg := serverpackets.SkillCoolTime()
	buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))

	return buffer.Bytes()
}
