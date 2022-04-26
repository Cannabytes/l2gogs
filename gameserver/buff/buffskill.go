package buff

import (
	"context"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
	"strconv"
	"time"
)

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
				unique = append(unique[:index], buff)
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
	_, err = dbConn.Exec(context.Background(), `DELETE FROM "buffs" WHERE "char_id" = $1`, charId)
	if err != nil {
		logger.Error.Panicln(err)
	}
}

// SaveBuff Сохранение баффа в БД, который на игроке
func SaveBuff(clientI interfaces.ReciverAndSender) {
	//сlearBuffListDB(clientI.(*models.Client).CurrentChar.ObjectId)
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
	_, err = dbConn.Exec(context.Background(), `INSERT INTO "buffs" ("char_id", "id", "level", "second") VALUES `+values)
	if err != nil {
		logger.Error.Panicln(err)
	}
}

func BuffTimeOut(ch *models.Character) {
	for {
		isNeedComparisonBuff := false
		if ch.IsOnline == false {
			return
		}
		if len(ch.Buff) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		for index, buff := range ch.Buff {
			buff.Second -= 1
			if buff.Second == 0 {
				isNeedComparisonBuff = true
				ch.Buff = append(ch.Buff[:index], ch.Buff[index+1:]...)
				ComparisonBuff(ch.Conn)
				continue
			}
		}
		if isNeedComparisonBuff {
			pkg17 := serverpackets.AbnormalStatusUpdate(ch.Buff)
			ch.EncryptAndSend(pkg17)
		}
		time.Sleep(1 * time.Second)
	}

}
