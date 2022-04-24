package buff

import (
	"context"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models/buff/buffdata"
	"strconv"
)

// TODO: Необходимо сделать элементарную проверку на дубликаты баффов, сохранения, чтения.
func removeBuffDuplicates(buffs []buffdata.BuffUser) {
}

// GetBuffSkill Получение из БД всех сохраненных баффов
func GetBuffSkill(charId int32) []buffdata.BuffUser {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	rows, err := dbConn.Query(context.Background(), "SELECT id, level, second FROM buffs WHERE char_id = $1", charId)
	if err != nil {
		logger.Error.Panicln(err)
	}
	var buffs []buffdata.BuffUser
	for rows.Next() {
		var buff buffdata.BuffUser
		err = rows.Scan(&buff.Id, &buff.Level, &buff.Second)
		if err != nil {
			logger.Error.Panicln(err)
		}
		buffs = append(buffs, buff)
	}
	if len(buffs) >= 1 {
		сlearBuffListDB(charId)
	}
	removeBuffDuplicates(buffs)
	return buffs
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
	MyBuffList := clientI.GetCurrentChar().GetBuff()
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
