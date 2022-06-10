package clientpackets

import (
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models/trade"
	"l2gogameserver/gameserver/serverpackets"
	"l2gogameserver/packets"
	"l2gogameserver/utils"
)

//AnswerTradeRequest Если пользователь отвечает на запрос трейда
func AnswerTradeRequest(data []byte, sender interfaces.ReciverAndSender) {
	var packet = packets.NewReader(data)
	response := packet.ReadInt32()
	if response == 0 {
		logger.Info.Println("Пользователь не захотел торговать")
		return
	}

	exchange, ok := trade.Answer(sender.Player())
	if ok {
		buffer := packets.Get()
		defer packets.Put(buffer)

		ut1 := utils.GetPacketByte()
		ut1.SetData(serverpackets.TradeStart(exchange.Sender.Client))
		exchange.Sender.Client.EncryptAndSend(ut1.GetData())

		ut := utils.GetPacketByte()
		ut.SetData(serverpackets.TradeStart(exchange.Recipient.Client))
		exchange.Recipient.Client.EncryptAndSend(ut.GetData())
	}

}
