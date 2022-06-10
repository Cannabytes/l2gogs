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

	if char.CurHp == clientI.Player().MaxHP() &&
		char.CurMp == clientI.Player().MaxMP() &&
		char.CurCp == clientI.Player().MaxCP() {
		return
	}
	char.CurHp += char.HpRegen
	if char.CurHp >= char.MaxHP() {
		char.CurHp = clientI.Player().MaxHP()
	}
	char.CurMp += char.MpRegen
	if char.CurMp >= char.MaxMP() {
		char.CurMp = clientI.Player().MaxMP()
	}
	char.CurCp += char.CpRegen
	if char.CurCp >= char.MaxCP() {
		char.CurCp = clientI.Player().MaxCP()
	}
	pkg := serverpackets.StatusUpdate(clientI)
	clientI.EncryptAndSend(pkg)
}
