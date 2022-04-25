package clientpackets

import (
	"l2gogameserver/gameserver/buff"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/serverpackets"
)

func Logout(client interfaces.ReciverAndSender, data []byte) {
	buff.SaveBuff(client)
	client.GetCurrentChar().SetStatusOffline()
	pkg := serverpackets.LogoutToClient(data, client)
	client.EncryptAndSend(pkg)
}
