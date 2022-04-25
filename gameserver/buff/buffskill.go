package buff

import (
	"context"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
	"log"
	"strconv"
	"time"
)

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
		if ch.IsOnline == false {
			log.Println("Персонаж вышел из игры, время баффа не отнимаем")
			return
		}
		if len(ch.Buff) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		for index, buff := range ch.Buff {
			if buff.Second == 1 {
				ch.Buff = append(ch.Buff[:index], ch.Buff[index+1:]...)
				logger.Warning.Println("Бафф сейчас должен сняться", buff.Id)
				pkg17 := serverpackets.AbnormalStatusUpdate(ch.Buff)
				ch.EncryptAndSend(pkg17)
				continue
			}
			buff.Second -= 1
			logger.Warning.Println("Баффу осталось", buff.Id, buff.Second)
		}
		time.Sleep(1 * time.Second)
	}

}
