// Данный пакет относится к работе схем баффов

package community

import (
	"context"
	"fmt"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"log"
)

// GetLoadCharacterScheme Список сохраненных комбинаций баффов
func GetLoadCharacterScheme(clientI interfaces.ReciverAndSender) {
	var all []*models.BuffScheme
	client, ok := clientI.(*models.Client)
	if !ok {
		return
	}
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
		return
	}
	defer dbConn.Release()
	sql := `SELECT id, char_id, name FROM "character_scheme" WHERE char_id=$1`
	we, err := dbConn.Query(context.Background(), sql, client.GetCurrentChar().GetObjectId())
	if err != nil {
		logger.Error.Println(err)
		return
	}
	for we.Next() {
		var bs models.BuffScheme
		err = we.Scan(&bs.Id, &bs.CharId, &bs.Name)
		if err != nil {
			logger.Error.Println(err)
		}
		bs.Buffs = getSchemeUserBuff(bs.Id)
		log.Println(bs)
		all = append(all, &bs)
	}
	client.CurrentChar.BuffScheme = all
}

// GetSchemeUserBuff Получение баффа схемы
func getSchemeUserBuff(schemeId int) []models.BuffSchemeSkill {
	var all []models.BuffSchemeSkill
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
		return all
	}
	defer dbConn.Release()
	sql := `SELECT skill_id, skill_level FROM "character_scheme_buffs" WHERE scheme_id=$1`
	we, err := dbConn.Query(context.Background(), sql, schemeId)
	if err != nil {
		logger.Error.Println(err)
		return all
	}
	for we.Next() {
		var bsk models.BuffSchemeSkill
		we.Scan(&bsk.SkillId, &bsk.SkillLevel)
		all = append(all, bsk)
	}
	return all
}

// SchemeSave Сохранение баффа персонажа в бд
func SchemeSave(clientI interfaces.ReciverAndSender, schemeName string) bool {
	client := clientI.(*models.Client).CurrentChar
	if len(client.Buff) == 0 {
		return false
	}
	schemeId, ok := createRegistryScheme(client.GetObjectId(), schemeName)
	if !ok {
		logger.Error.Panicln("Добавление новой записи схемы бафа не произошла")
		return false
	}
	return createSchemeListBuff(client.Buff, schemeId)
}

//Регистрируем новую схему
func createRegistryScheme(char_id int32, name string) (int, bool) {
	lastInsertId := 0
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
		return lastInsertId, false
	}
	defer dbConn.Release()
	err = dbConn.QueryRow(context.Background(), `INSERT INTO "character_buffs_save_list" ("char_id", "name") VALUES ($1, $2) RETURNING id`, char_id, name).Scan(&lastInsertId)
	if err != nil {
		logger.Error.Println(err)
		return lastInsertId, false
	}
	return lastInsertId, true
}

// Записываем все скиллы персонажа в новую схему
func createSchemeListBuff(bufflist []*models.BuffUser, schemeId int) bool {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
		return false
	}
	defer dbConn.Release()
	buffcount := len(bufflist)
	sql := `INSERT INTO "character_buffs_scheme" ("scheme_id", "skill_id", "skill_level") VALUES`
	for index, buff := range bufflist {
		sql += fmt.Sprintf("(%d, %d, %d)", schemeId, buff.Id, buff.Level)
		if buffcount != index+1 {
			sql += ","
		}
	}
	dbConn.Exec(context.Background(), sql)
	return true
}
