package admin

import (
	"l2gogameserver/data"
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
	"log"
)

func IsAdmin() {

}

func Command(clientInterface interfaces.ReciverAndSender, commandArr []string) {
	log.Println(commandArr)
	command := commandArr[0]
	//summon - вызвать предмет
	switch command {
	case "summon":
		itemid, count := data.StrToInt(commandArr[1]), data.StrToInt64(commandArr[2])
		itemSummon(clientInterface, itemid, count)
	case "teleport":
		teleport(clientInterface, data.StrToInt(commandArr[1]), data.StrToInt(commandArr[2]))
	}

}

//Когда админ в ALT+G призывает предмет
func itemSummon(clientInterface interfaces.ReciverAndSender, itemid int, count int64) {
	client := clientInterface.(*models.Client)
	item, ok := models.NewItemCreate(client.CurrentChar, itemid, count)
	if !ok {
		logger.Error.Println("Предмет не создался")
		return
	}
	item, updateType := client.CurrentChar.Inventory.AddItem(item)
	client.EncryptAndSend(serverpackets.InventoryUpdate(item, updateType))
}

//Телепортация админа
//Ему достаточно зажать CTRL+ALT+SHIFT и кликнуть по карте
func teleport(clientInterface interfaces.ReciverAndSender, x, y int) {
	z, h := 0, 0
	clientInterface.EncryptAndSend(serverpackets.TeleportToLocation(clientInterface, x, y, z, h))
}
