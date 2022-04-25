package command

import (
	"l2gogameserver/gameserver/clientpackets"
	"l2gogameserver/gameserver/interfaces"
)

//Проверка на реализацию комманды
//По аналогии с настройками, в будущем можно будет так же вызывать .buff, .shop, .teleport, .7rb, .rb ...
//В идеале, чтоб открывалось комьюнити с информацией
func ExistCommand(commandTxt string, me interfaces.ReciverAndSender) bool {
	if commandTxt == ".cfg" || commandTxt == ".menu" {
		openMenu(me)
		return true
	}
	return false
}

func openMenu(client interfaces.ReciverAndSender) {
	clientpackets.SendOpenDialogBBS(client, "./datapack/html/community/setting/setting.htm")
}
