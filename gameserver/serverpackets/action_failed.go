package serverpackets

import (
	"l2gogameserver/gameserver/models"
)

func ActionFailed(client *models.Client) {
	client.Buffer.WriteSingleByte(0x1f)
}
