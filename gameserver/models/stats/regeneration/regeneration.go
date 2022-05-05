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

	if char.CurHp == clientI.GetCurrentChar().GetMaxHP() &&
		char.CurMp == clientI.GetCurrentChar().GetMaxMP() &&
		char.CurCp == clientI.GetCurrentChar().GetMaxCP() {
		return
	}
	char.CurHp += char.HpRegen
	if char.CurHp >= char.MaxHp {
		char.CurHp = clientI.GetCurrentChar().GetMaxHP()
	}
	char.CurMp += char.MpRegen
	if char.CurMp >= char.MaxMp {
		char.CurMp = clientI.GetCurrentChar().GetMaxMP()
	}
	char.CurCp += char.CpRegen
	if char.CurCp >= char.MaxCp {
		char.CurCp = clientI.GetCurrentChar().GetMaxCP()
	}
	pkg := serverpackets.StatusUpdate(clientI)
	clientI.EncryptAndSend(pkg)
}
