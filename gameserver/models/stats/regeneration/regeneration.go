package regeneration

import (
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
)

// Генерация ХП МП СП
// Если необходимо обновить данные параметров, то шлем запрос
func RenerationHpMpCp(clientI interfaces.ReciverAndSender) {
	char := clientI.(*models.Client).CurrentChar
	//Если хп одинаковое - выходим, потом доделать так же с MP/CP
	if char.CurHp == char.MaxHp && char.CurMp == char.MaxMp && char.CurCp == char.MaxCp {
		return
	}
	char.CurHp += int32(char.HpRegen)
	if char.CurHp >= char.MaxHp {
		char.CurHp = char.MaxHp
	}
	char.CurMp += int32(char.MpRegen)
	if char.CurMp >= char.MaxMp {
		char.CurMp = char.MaxMp
	}
	char.CurCp += int32(char.CpRegen)
	if char.CurCp >= char.MaxCp {
		char.CurCp = char.MaxCp
	}
	pkg := serverpackets.StatusUpdate(clientI)
	clientI.EncryptAndSend(pkg)
}
