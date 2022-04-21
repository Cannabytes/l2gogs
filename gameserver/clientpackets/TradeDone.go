package clientpackets

import (
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/models/trade"
	"l2gogameserver/gameserver/serverpackets"
	"l2gogameserver/packets"
	"l2gogameserver/utils"
	"log"
)

//Игрок подтвердил сделку
func TradeDone(data []byte, client interfaces.ReciverAndSender) {
	var packet = packets.NewReader(data)
	response := packet.ReadInt32() // 1 - пользователь нажал ОК, 0 пользователь отменил трейд

	player2, exchange, ok := trade.FindTrade(client.GetCurrentChar())
	if !ok {
		log.Println("Обменивающихся не найдено")
		return
	}
	if response == 1 {
		if exchange.Sender.ObjectId == client.GetCurrentChar().GetObjectId() {
			exchange.Sender.Completed = true
			log.Printf("Игрок %s подтвердил сделку\n", client.GetCurrentChar().GetName())
			serverpackets.TradeOtherDone(player2)
		}
		if exchange.Recipient.ObjectId == client.GetCurrentChar().GetObjectId() {
			exchange.Recipient.Completed = true
			log.Printf("Игрок %s подтвердил сделку\n", client.GetCurrentChar().GetName())
			serverpackets.TradeOtherDone(client.GetCurrentChar())
		}
		if exchange.Recipient.Completed == exchange.Sender.Completed {
			log.Println("Обмен завершен успешно")
			serverpackets.TradeOK(client.GetCurrentChar(), player2)
			//Теперь сделаем физическую передачу предметов от персонажа к персонажу
			cplayer, toplayer := trade.TradeAddInventory(client.GetCurrentChar(), player2, exchange)

			buffer := packets.Get()
			defer packets.Put(buffer)

			for _, item := range cplayer {
				log.Println(item)
				ut1 := utils.GetPacketByte()
				ut1.SetData(serverpackets.InventoryUpdate(item, models.UpdateTypeModify))
				client.EncryptAndSend(ut1.GetData())
			}

			for _, item := range toplayer {
				log.Println(item)

				ut1 := utils.GetPacketByte()
				ut1.SetData(serverpackets.InventoryUpdate(item, models.UpdateTypeModify))
				player2.EncryptAndSend(ut1.GetData())
			}

			if ok = trade.TradeUserClear(client.GetCurrentChar()); !ok {
				log.Println("Трейд не найден")
				return
			}
		}
	} else if response == 0 {
		serverpackets.TradeCancel(client.GetCurrentChar(), player2)
		if ok = trade.TradeUserClear(client.GetCurrentChar()); !ok {
			log.Println("Трейд не найден")
			return
		}
	}

}