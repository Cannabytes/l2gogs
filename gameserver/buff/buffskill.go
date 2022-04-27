package buff

import (
	"context"
	"encoding/json"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
	"os"
	"strconv"
	"time"
)

type Combo struct {
	ID    int `json:"id"`
	Buffs []struct {
		SkillID int `json:"skill_id"`
		Level   int `json:"level"`
	} `json:"buffs"`
	Time    int    `json:"time"`
	CostID  int    `json:"cost_id"`
	Cost    int    `json:"cost"`
	Comment string `json:"comment"`
}

var communityComboBuff = []Combo{}

// LoadCommunityComboBuff Загрузка комбо баффа комьюнити
func LoadCommunityComboBuff() {
	file, err := os.Open("./config/community/combo_buff.json")
	if err != nil {
		logger.Error.Panicln("Failed to load config file " + err.Error())
	}
	err = json.NewDecoder(file).Decode(&communityComboBuff)
	if err != nil {
		logger.Error.Panicln("Failed to decode config file " + file.Name() + " " + err.Error())
	}
}

// GetCommunityComboBuff Возвращает информацию о комбо баффах
func GetCommunityComboBuff(id int) (Combo, bool) {
	for _, combo := range communityComboBuff {
		if combo.ID == id {
			return combo, true
		}
	}
	return Combo{}, false
}

// ComparisonBuff Функция сравнения баффов
// Убирает дубликаты скиллов, и оставляет бафф (если одинаковый лвл) который больше по времени будет действовать
// Если время баффов одинаковое, тогда применяется бафф больше по уровню
func ComparisonBuff(clientI interfaces.ReciverAndSender) {
	client := clientI.(*models.Client)
	buffList := client.CurrentChar.GetBuff()
	var unique []*models.BuffUser
	buffGet := func(unique []*models.BuffUser, id int) (*models.BuffUser, int, bool) {
		for index, buff := range unique {
			if buff.Id == id {
				return buff, index, true
			}
		}
		return nil, 0, false
	}
	for _, buff := range buffList {
		duplicateBuff, index, ok := buffGet(unique, buff.Id)
		if ok {
			if duplicateBuff.Second < buff.Second || duplicateBuff.Second == buff.Second && duplicateBuff.Level < buff.Level {
				unique = append(unique[:index], unique[index+1:]...)
				unique = append(unique, buff)
			}
		} else {
			unique = append(unique, buff)
			continue
		}
	}
	client.CurrentChar.Buff = unique
}

// Очищает записи баффа в БД
func сlearBuffListDB(charId int32) {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()
	_, err = dbConn.Exec(context.Background(), `DELETE FROM "character_buffs" WHERE "char_id" = $1`, charId)
	if err != nil {
		logger.Error.Panicln(err)
	}
}

// SaveBuff Сохранение баффа в БД, который на игроке
func SaveBuff(clientI interfaces.ReciverAndSender) {
	сlearBuffListDB(clientI.(*models.Client).CurrentChar.ObjectId)
	MyBuffList := clientI.(*models.Client).CurrentChar.GetBuff()
	buffCount := len(MyBuffList)
	if buffCount == 0 {
		return
	}
	playerId := strconv.Itoa(int(clientI.GetCurrentChar().GetObjectId()))
	var values string
	for index, buff := range MyBuffList {
		id := strconv.Itoa(buff.Id)
		level := strconv.Itoa(buff.Level)
		second := strconv.Itoa(buff.Second)
		values += "(" + playerId + ", " + id + ", " + level + ", " + second + ")"
		if index+1 != buffCount {
			values += ", "
		}
	}

	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()
	_, err = dbConn.Exec(context.Background(), `INSERT INTO "character_buffs" ("char_id", "id", "level", "second") VALUES `+values)
	if err != nil {
		logger.Error.Panicln(err)
	}
}

func BuffTimeOut(ch *models.Character) {
	var buffRemove = []*models.BuffUser{}
	for {
		isNeedComparisonBuff := false
		if ch.IsOnline == false {
			return
		}
		if len(ch.Buff) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		for _, buff := range ch.Buff {
			buff.Second -= 1
			if buff.Second == 0 {
				isNeedComparisonBuff = true
				buffRemove = append(buffRemove, buff)
				continue
			}
		}
		if isNeedComparisonBuff {
			for _, bf := range buffRemove {
				RemoveBuffId(ch, bf.Id)
			}
			ComparisonBuff(ch.Conn)
			pkg17 := serverpackets.AbnormalStatusUpdate(ch.Buff)
			ch.EncryptAndSend(pkg17)
		}
		time.Sleep(1 * time.Second)
	}
}

// Удаление баффа у игрока по ID баффа
func RemoveBuffId(ch *models.Character, Id int) {
	for index, buff := range ch.Buff {
		if buff.Id == Id {
			ch.Buff = append(ch.Buff[:index], ch.Buff[index+1:]...)
		}
	}
}
