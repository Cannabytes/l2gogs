package clientpackets

import (
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models/buff"
	"l2gogameserver/gameserver/serverpackets"
)

func Logout(client interfaces.ReciverAndSender, data []byte) {
	buff.SaveBuff(client)
	pkg := serverpackets.LogoutToClient(data, client)
	client.EncryptAndSend(pkg)
}
