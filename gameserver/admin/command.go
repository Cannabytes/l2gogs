package admin

import (
	"l2gogameserver/data"
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
)

func IsAdmin() {

}

func Command(client interfaces.ReciverAndSender, commandArr []string) {
	command := commandArr[0]
	//summon - вызвать предмет
	if command == "summon" {
		itemid, count := data.StrToInt(commandArr[1]), data.StrToInt64(commandArr[2])
		itemSummon(client, itemid, count)
	}
}

func itemSummon(clientInterface interfaces.ReciverAndSender, itemid int, count int64) {
	client := clientInterface.(*models.Client)
	item, ok := models.NewItemCreate(client.CurrentChar, itemid, count)
	if !ok {
		logger.Error.Println("Предмет не создался")
		return
	}
	client.CurrentChar.Inventory.AddItem(item)
	//for _, myItem := range client.CurrentChar.Inventory.Items {
	//log.Println(myItem.Name, myItem.Count)
	//}
}
