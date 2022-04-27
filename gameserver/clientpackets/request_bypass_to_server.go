package clientpackets

import (
	inter "l2gogameserver/data"
	"l2gogameserver/data/logger"
	buff2 "l2gogameserver/gameserver/buff"
	"l2gogameserver/gameserver/community"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/models/htm"
	"l2gogameserver/gameserver/models/multisell"
	"l2gogameserver/gameserver/serverpackets"
	"l2gogameserver/packets"
	"strconv"
	"strings"
	"time"
)

/*
	Пока заметка, направление как делать.
	Разберем на будущее парсинг bypass
	Все запросы на открытие страницы будут начинатся с _bbspage
	следующие параметры разделены двоеточием:
	[вызов страницы]:[команда]:[информация]:[информация]:[информация]...
	_bbspage:open:/page/index.htm (аналог _bbspage:open:page) - открыть файл (bypass -h _bbspage:open:buffer/buffs.htm)
	_bbspage:buffer:combo:3 - наложение комбо баффа с ID 3
	_bbspage:buffer:save - сохранить бафф персонажа
	_bbspage:buffer:get:3 - бафф персонажа (ранее сохраненным баффом) с ID 3
	// Другие аналогия
	_bbspage:gmshop:multisell:1531 - открыть мультиселл 1531
	_bbspage:teleport:id:152 - Телепорт по координатам с ID 152
	_bbspage:teleport:save	- сохранение позиции (xyz) персонажа
	_bbspage:teleport:to:5 - телепорт ранее сохраненную позицию с ID 5
	_bbspage:teleport:remove:5 - удаление сохраненной точки с ID 5
	...
	// Функция наложения баффов
	_bbsbuff:1204:2:0:page/index.htm - где ID баффа 1204 и уровень 2 и цена 0, следующим параметром отправляем данные о том какую страницу открыть
*/
func BypassToServer(data []byte, clientI interfaces.ReciverAndSender) {
	client := clientI.(*models.Client).CurrentChar

	var bypassRequest = packets.NewReader(data).ReadString()
	bypassInfo := strings.Split(bypassRequest, ":")
	for i, s := range bypassInfo {
		logger.Info.Println("#", i, "->", s)
	}
	logger.Info.Println(bypassInfo)
	if bypassInfo[0] == "_bbshome" && bypassRequest == "_bbshome" {
		//Открытие диалога по умолчанию
		SendOpenDialogBBS(clientI, "./datapack/html/community/index.htm")
	} else if bypassInfo[0] == "_bbspage" {
		commandname := bypassInfo[1]
		switch commandname {
		//Запрос открытия диалога
		case "open":
			SendOpenDialogBBS(clientI, "./datapack/html/community/"+bypassInfo[2])
		//Функции телепортации
		case "teleport":
			switch bypassInfo[2] {
			case "id":
				teleportID, err := strconv.Atoi(bypassInfo[3])
				if err != nil {
					logger.Info.Println(err)
					return
				}
				pkg := community.UserTeleport(clientI, teleportID)
				clientI.EncryptAndSend(pkg)
			case "save":
				logger.Info.Println("Сохранение позиции игрока")
			case "to":
				logger.Info.Println("Телепорт по сохраненной позиции игрока #", bypassInfo[3])
			case "remove":
				logger.Info.Println("Удаление по сохраненной позиции игрока #", bypassInfo[3])
			}
		case "gmshop":
			switch bypassInfo[2] {
			case "multisell": //Open multisell
				id, err := strconv.Atoi(bypassInfo[3])
				if err != nil {
					logger.Info.Println(err)
					return
				}
				logger.Info.Println("Открыть мультиселл с ID", id)
				multisellList, ok := multisell.Get(clientI, id)
				if !ok {
					logger.Info.Println("Не найден запрашиваемый мультисел#")
				}
				pkg := serverpackets.MultiSell(multisellList)
				clientI.EncryptAndSend(pkg)
			}
		}

	} else if bypassInfo[0] == "_bbsbuff" {
		//Функция наложения баффа на персонажа
		buffId := bypassInfo[1]
		buffLevel := bypassInfo[2]
		buffCost := bypassInfo[3]
		_ = buffCost
		client.Buff = append(client.Buff, &models.BuffUser{
			Id:     inter.StrToInt(buffId),
			Level:  inter.StrToInt(buffLevel),
			Second: 60,
		})
		buff2.ComparisonBuff(client.Conn)
		pkg17 := serverpackets.AbnormalStatusUpdate(client.Buff)
		client.EncryptAndSend(pkg17)

		if len(bypassInfo) == 5 {
			page := bypassInfo[4]
			SendOpenDialogBBS(clientI, "./datapack/html/community/"+page)
		}

	}
}

//SendOpenDialogBBS Открытие диалога и отправка клиенту диалога
func SendOpenDialogBBS(client interfaces.ReciverAndSender, filename string) {
	logger.Info.Println(filename)
	htmlDialog, err := htm.Open(filename)
	if err != nil {
		logger.Info.Println(err)
		return
	}
	htmlDialog = parseVariableBoard(client, htmlDialog)
	bufferDialog := packets.Get()
	defer packets.Put(bufferDialog)
	bufferDialog1 := packets.Get()
	defer packets.Put(bufferDialog1)
	bufferDialog2 := packets.Get()
	defer packets.Put(bufferDialog2)

	if len(*htmlDialog) < 8180 {
		bufferDialog.WriteSlice(models.ShowBoard(*htmlDialog, "101"))
		bufferDialog1.WriteSlice(models.ShowBoard("", "102"))
		bufferDialog2.WriteSlice(models.ShowBoard("", "103"))
	} else if len(*htmlDialog) < 8180*2 {
		bufferDialog.WriteSlice(models.ShowBoard((*htmlDialog)[:8180], "101"))
		bufferDialog1.WriteSlice(models.ShowBoard((*htmlDialog)[8180:], "102"))
		bufferDialog2.WriteSlice(models.ShowBoard("", "103"))
	} else if len(*htmlDialog) < 8180*3 {
		bufferDialog.WriteSlice(models.ShowBoard((*htmlDialog)[:8180], "101"))
		bufferDialog1.WriteSlice(models.ShowBoard((*htmlDialog)[8180:8180*2], "102"))
		bufferDialog2.WriteSlice(models.ShowBoard((*htmlDialog)[8180*2:], "103"))
	}
	buffer := packets.Get()
	buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(bufferDialog.Bytes()))
	buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(bufferDialog1.Bytes()))
	buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(bufferDialog2.Bytes()))
	client.Send(buffer.Bytes())

	packets.Put(buffer)
}

//parseVariableBoard Псевдопеременные из html комьюнити заменяем реальными
func parseVariableBoard(client interfaces.ReciverAndSender, html *string) *string {
	r := strings.NewReplacer(
		"<?player_name?>", client.GetCurrentChar().GetName(),
		"<?player_class?>", strconv.Itoa(int(client.GetCurrentChar().GetClassId())),
		"<?cb_time?>", time.Now().Format(time.RFC850),
	)
	result := r.Replace(*html)
	return &result
}
