package clientpackets

import (
	"l2gogameserver/gameserver/community"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/models/htm"
	"l2gogameserver/packets"
	"log"
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
	_bbspage:open:/page/index.htm (аналог _bbspage:open:page) - открыть файл
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
*/
func BypassToServer(data []byte, client *models.Client) {
	var bypassRequest = packets.NewReader(data).ReadString()
	bypassInfo := strings.Split(bypassRequest, ":")
	for i, s := range bypassInfo {
		log.Println("#", i, "->", s)
	}
	log.Println(bypassInfo)
	if bypassInfo[0] == "_bbshome" && bypassRequest == "_bbshome" {
		//Открытие диалога по умолчанию
		SendOpenDialogBBS(data, client, "./server/html/community/index.htm")
	} else if bypassInfo[0] == "_bbspage" {
		commandname := bypassInfo[1]
		switch commandname {
		//Запрос открытия диалога
		case "open":
			SendOpenDialogBBS(data, client, "./server/html/community/"+bypassInfo[2])

		//Функции телепортации
		case "teleport":
			switch bypassInfo[2] {
			case "id":
				teleportID, err := strconv.Atoi(bypassInfo[3])
				if err != nil {
					log.Println(err)
					return
				}
				community.UserTeleport(client, teleportID)
			case "save":
				log.Println("Сохранение позиции игрока")
			case "to":
				log.Println("Телепорт по сохраненной позиции игрока #", bypassInfo[3])
			case "remove":
				log.Println("Удаление по сохраненной позиции игрока #", bypassInfo[3])
			}

		case "gmshop":
			switch bypassInfo[2] {
			case "multisell":
				//Open multisell
				log.Println("Открыть мультиселл с ID", bypassInfo[3])
			}

		}

	}
}

//Открытие диалога и отправка клиенту диалога
func SendOpenDialogBBS(data []byte, client *models.Client, filename string) {
	log.Println(filename)
	htmlDialog, err := htm.Open(filename)
	if err != nil {
		log.Println(err)
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
	client.SSend(buffer.Bytes())
}

//Псевдопеременные из html комьюнити заменяем реальными
func parseVariableBoard(client *models.Client, htmlcode *string) *string {
	r := strings.NewReplacer(
		"<?player_name?>", client.CurrentChar.CharName,
		"<?player_class?>", strconv.Itoa(int(client.CurrentChar.ClassId)),
		"<?cb_time?>", time.Now().Format(time.RFC850),
	)
	result := r.Replace(*htmlcode)
	return &result
}