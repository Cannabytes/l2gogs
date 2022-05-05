package regeneration

import (
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
)

// Генерация ХП МП СП
// Если необходимо обновить данные параметров, то шлем запрос
func RenerationHpMpCp(clientI interfaces.ReciverAndSender) bool {
	char := clientI.(*models.Client).CurrentChar
	//Если хп одинаковое - выходим, потом доделать так же с MP/CP
	if char.CurHp == char.MaxHp {
		return true
	}
	char.CurHp += int32(char.HpRegen)
	if char.CurHp >= char.MaxHp {
		char.CurHp = char.MaxHp
	}
	pkg := serverpackets.StatusUpdate(clientI)
	clientI.EncryptAndSend(pkg)
	return false
}
