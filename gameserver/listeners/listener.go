package listeners

import (
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver"
	"l2gogameserver/gameserver/broadcast"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
	"l2gogameserver/packets"
	"l2gogameserver/utils"
)

func StartClientListener(client interfaces.ReciverAndSender) {
	go channelListener(client)
	go npcListener(client)
	go moveListener(client)
}
func channelListener(client interfaces.ReciverAndSender) {
	ch, ok := client.(*models.Client)
	if !ok {
		logger.Error.Panicln("ChannelListenerlogger.Error.Panicln")
	}

	for q := range ch.CurrentChar.ChannelUpdateShadowItem {
		pkg := serverpackets.ItemUpdate(client, q.UpdateType, q.ObjId)
		client.EncryptAndSend(pkg)
		if q.UpdateType == models.UpdateTypeRemove {
			broadcast.BroadCastUserInfoInRadius(client, 2000)
		}
	}
}

func npcListener(client interfaces.ReciverAndSender) {
	ch, ok := client.(*models.Client)
	if !ok {
		logger.Error.Panicln("NpcListenerlogger.Error.Panicln")
	}
	for q := range ch.CurrentChar.NpcInfo {
		buff := packets.Get()
		for i := range q {
			pkg := serverpackets.NpcInfo(q[i])
			buff.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))
		}
		client.Send(buff.Bytes())
		packets.Put(buff)
	}
}

func moveListener(client interfaces.ReciverAndSender) {
	ch, ok := client.(*models.Client)
	if !ok {
		logger.Error.Panicln("NpcListenerlogger.Error.Panicln")
	}

	pkg := utils.GetPacketByte()
	defer pkg.Release()

	for q := range ch.CurrentChar.CharInfoTo {
		pkg.SetData(serverpackets.CharInfo(ch.CurrentChar))
		for _, v := range q {
			gameserver.OnlineCharacters.Mu.Lock()
			gameserver.OnlineCharacters.Char[v].Conn.EncryptAndSend(pkg.GetData())
			gameserver.OnlineCharacters.Mu.Unlock()
		}
	}

	pkg.Free()

	for q := range ch.CurrentChar.DeleteObjectTo {
		pkg.SetData(serverpackets.DeleteObject(ch.CurrentChar))
		for _, v := range q {
			gameserver.OnlineCharacters.Mu.Lock()
			gameserver.OnlineCharacters.Char[v].Conn.EncryptAndSend(pkg.GetData())
			gameserver.OnlineCharacters.Mu.Unlock()
		}
	}

}
