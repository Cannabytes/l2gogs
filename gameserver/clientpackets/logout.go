package clientpackets

import (
	"l2gogameserver/gameserver"
	"l2gogameserver/gameserver/buff"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
	"log"
)

func Logout(clientI interfaces.ReciverAndSender, data []byte) {
	client := clientI.(*models.Client)

	log.Println(clientI.Player().ObjectID())
	if clientI.Player().ObjectID() == 0 {
		return
	}

	client.CurrentChar.InGame = false
	buff.SaveBuff(clientI)
	client.CurrentChar.Inventory.Save(int(clientI.Player().ObjectID()))

	clientI.Player().SetStatusOffline()
	pkg := serverpackets.LogoutToClient(data, clientI)
	clientI.EncryptAndSend(pkg)
	gameserver.CharOffline(clientI)

	client.Socket.Close()
}
