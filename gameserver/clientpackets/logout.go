package clientpackets

import (
	"l2gogameserver/gameserver/buff"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/serverpackets"
	"log"
)

func Logout(client interfaces.ReciverAndSender, data []byte) {
	log.Println(client.GetCurrentChar().GetObjectId())
	if client.GetCurrentChar().GetObjectId() == 0 {
		return
	}
	buff.SaveBuff(client)
	client.GetCurrentChar().SetStatusOffline()
	pkg := serverpackets.LogoutToClient(data, client)
	client.EncryptAndSend(pkg)
}
