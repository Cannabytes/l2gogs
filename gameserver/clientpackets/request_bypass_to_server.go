package clientpackets

import (
	"bytes"
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
	"text/template"
	"time"
)

/*
	Пока заметка, направление как делать.
	Разберем на будущее парсинг bypass
	Все запросы на открытие страницы будут начинатся с _bbspage
	следующие параметры разделены двоеточием:
	[вызов страницы]:[команда]:[информация]:[информация]:[информация]...
	_bbspage:open:/page/index.htm (аналог _bbspage:open:page) - открыть файл (bypass -h _bbspage:open:buffer/buffs.htm)
	// Другие аналогия
	_bbspage:gmshop:multisell:1531 - открыть мультиселл 1531
	_bbspage:teleport:id:152 - Телепорт по координатам с ID 152
	_bbspage:teleport:save	- сохранение позиции (xyz) персонажа
	_bbspage:teleport:to:5 - телепорт ранее сохраненную позицию с ID 5
	_bbspage:teleport:remove:5 - удаление сохраненной точки с ID 5
	...
	// Функция наложения баффов
	_bbsbuff:cast:1204:2:0:page/index.htm - где ID баффа 1204 и уровень 2 и цена 0, следующим параметром отправляем данные о том какую страницу открыть
	_bbsbuff:combo:3 - наложение комбо баффа с ID 3
	_bbsbuff:cancel - отменяет весь бафф на пользователе
	_bbsbuff:scheme:create: $name  - создание новой схемы баффа, пробел обязательный после create: (иначе клиент не передает переменную)
	_bbsbuff:scheme:get:3 - бафф персонажа (ранее сохраненным баффом) с ID 3

*/
func BypassToServer(data []byte, clientI interfaces.ReciverAndSender) {
	client := clientI.(*models.Client)

	var bypassRequest = packets.NewReader(data).ReadString()
	bypassInfo := strings.Split(bypassRequest, ":")
	//for i, s := range bypassInfo {
	//logger.Info.Println("#", i, "->", s)
	//}
	//logger.Info.Println(bypassInfo)
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
		buffAnalysis(clientI, bypassInfo)
		clientI.EncryptAndSend(serverpackets.UserInfo(client))
	}
}

func buffAnalysis(clientI interfaces.ReciverAndSender, bypassInfo []string) {
	client := clientI.(*models.Client).CurrentChar
	buffCommand := bypassInfo[1]
	if buffCommand == "cast" {
		//Функция наложения баффа на персонажа
		buffId := bypassInfo[2]
		buffLevel := bypassInfo[3]
		buffCost := bypassInfo[4]
		_ = buffCost

		client.AddBuff(inter.StrToInt(buffId), inter.StrToInt(buffLevel), 60)

		client.EncryptAndSend(serverpackets.AbnormalStatusUpdate(client.Buff()))

		if len(bypassInfo) == 6 {
			page := bypassInfo[5]
			SendOpenDialogBBS(clientI, "./datapack/html/community/"+page)
		}
		return
	}
	if buffCommand == "combo" {
		comboId := inter.StrToInt(bypassInfo[2])
		combobuff, ok := buff2.GetCommunityComboBuff(comboId)
		if !ok {
			logger.Error.Printf("Комбо %d не найдено\n", comboId)
			return
		}
		for _, buff := range combobuff.Buffs {
			client.AddBuff(buff.SkillID, buff.Level, combobuff.Time)
		}
		client.EncryptAndSend(serverpackets.AbnormalStatusUpdate(client.Buff()))

		if len(bypassInfo) == 4 {
			page := bypassInfo[3]
			SendOpenDialogBBS(clientI, "./datapack/html/community/"+page)
		}
		return
	}
	//Отмена всего баффа
	if buffCommand == "cancel" {
		client.ClearBuff()
		pkg17 := serverpackets.AbnormalStatusUpdate(client.Buff())
		client.EncryptAndSend(pkg17)
		if len(bypassInfo) == 3 {
			page := bypassInfo[2]
			SendOpenDialogBBS(clientI, "./datapack/html/community/"+strings.Trim(page, " "))
		}
		return
	}
	//Отвечает за создание схемы баффа
	if buffCommand == "scheme" {
		action := bypassInfo[2]
		if action == "create" {
			//Создание нового профиля (схем)
			schemeName := strings.Trim(bypassInfo[3], " ")
			if community.SchemeSave(clientI, schemeName) {
				if len(bypassInfo) == 5 {
					page := bypassInfo[4]
					SendOpenDialogBBS(clientI, "./datapack/html/community/"+strings.Trim(page, " "))
				}
			}
		} else if action == "get" {
			//Наложение баффа из профиля
			id := inter.StrToInt(bypassInfo[3])
			client.ClearBuff()
			for _, scheme := range client.BuffScheme {
				if scheme.Id == id {
					for _, buff := range scheme.Buffs {
						client.AddBuff(buff.SkillId, buff.SkillLevel, 60)
					}
				}
			}
			pkg17 := serverpackets.AbnormalStatusUpdate(client.Buff())
			client.EncryptAndSend(pkg17)
			if len(bypassInfo) == 5 {
				page := bypassInfo[4]
				SendOpenDialogBBS(clientI, "./datapack/html/community/"+strings.Trim(page, " "))
			}
		}

		return
	}

}

//SendOpenDialogBBS Открытие диалога и отправка клиенту диалога
func SendOpenDialogBBS(client interfaces.ReciverAndSender, filename string) {
	//logger.Info.Println(filename)
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
func parseVariableBoard(clientI interfaces.ReciverAndSender, htmlcode *string) *string {
	client := clientI.(*models.Client).CurrentChar

	var Data = struct {
		Player_name  string
		Player_class string
		Server_time  string
		Buff_scheme  []*models.BuffScheme
	}{
		Player_name:  clientI.Player().PlayerName(),
		Player_class: strconv.Itoa(int(clientI.Player().ClassID())),
		Server_time:  time.Now().Format(time.Stamp),
		Buff_scheme:  client.BuffScheme,
	}
	var tpl bytes.Buffer
	t, err := template.New("").Parse(*htmlcode)
	if err != nil {
		logger.Error.Panicln(err)
	}
	if err = t.Execute(&tpl, Data); err != nil {
		logger.Error.Panicln(err)
	}

	resultStr := tpl.String()
	return &resultStr
}
