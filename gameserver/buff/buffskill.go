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

var communityComboBuff []Combo

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
	MyBuffList := clientI.(*models.Client).CurrentChar.Buff()
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

//Отсчет времени баффа
func BuffTimeOut(ch *models.Character) {
	for {
		isNeedComparisonBuff := false
		if ch.InGame == false {
			return
		}
		if len(ch.Buff()) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		for _, buff := range ch.Buff() {
			buff.Second -= 1
			if buff.Second == 0 {
				isNeedComparisonBuff = true
				ch.RemoveBuffSkill(buff.Id)
			}
		}
		if isNeedComparisonBuff {
			ch.EncryptAndSend(serverpackets.AbnormalStatusUpdate(ch.Buff()))
			ch.StatsRefresh()
			ch.EncryptAndSend(serverpackets.UserInfo(ch.Conn))
		}
		time.Sleep(1 * time.Second)
	}
}
