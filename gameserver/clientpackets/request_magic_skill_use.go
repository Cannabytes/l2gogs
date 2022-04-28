package clientpackets

import (
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/models/sysmsg"
	"l2gogameserver/gameserver/serverpackets"
	"l2gogameserver/packets"
)

func RequestMagicSkillUse(data []byte, clientI interfaces.ReciverAndSender) {

	var packet = packets.NewReader(data)

	magicId := int(packet.ReadInt32())           // Identifier of the used skill
	ctrlPressed := packet.ReadInt32() != 0       // True if it's a ForceAttack : Ctrl pressed
	shiftPressed := packet.ReadSingleByte() != 0 // True if Shift pressed

	SkillCast(clientI, magicId, ctrlPressed, shiftPressed)
}

func SkillCast(clientI interfaces.ReciverAndSender, magicId int, ctrlPressed bool, shiftPressed bool) {
	client, ok := clientI.(*models.Client)
	if !ok {
		return
	}

	buffer := packets.Get()
	defer packets.Put(buffer)

	if client.CurrentChar.IsDead {
		pkg := serverpackets.ActionFailed(client)
		buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))
		client.Send(buffer.Bytes())
		return
	}

	if client.CurrentChar.IsFakeDeath {
		pkg := serverpackets.SystemMessage(sysmsg.CantMoveSitting, client)
		buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))
		pkg2 := serverpackets.ActionFailed(client)
		buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg2))
		client.Send(buffer.Bytes())
		return
	}

	var skill models.Skill
	for _, getskill := range client.CurrentChar.Skills {
		if getskill.SkillId == int(magicId) {
			skill = getskill
			break
		}
	}

	//if если скилла нет {
	// todo тут еще идут проверки, возможно это кастомный? скилл или скилл трансформы и если нет то фейл
	//	pkg := serverpackets.ActionFailed(client)
	//	buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))
	//	client.Send(buffer.Bytes())
	//	return
	//}

	_, _, _ = magicId, ctrlPressed, shiftPressed

	//if skill.OperateType.IsPassive() {
	//	pkg := serverpackets.ActionFailed(client)
	//	buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))
	//	client.Send(buffer.Bytes())
	//	return
	//}

	if client.CurrentChar.IsCastingNow {
		currSkill := client.CurrentChar.CurrentSkill
		if currSkill != nil && skill.SkillId == currSkill.Skill.SkillId {
			pkg := serverpackets.ActionFailed(client)
			buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))
			client.Send(buffer.Bytes())
			return
		} //todo тут еще есть elseif isSkillDisabled()

		client.CurrentChar.SetSkillToQueue(skill, ctrlPressed, shiftPressed)
		pkg := serverpackets.ActionFailed(client)
		buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))
		client.Send(buffer.Bytes())
		return
	}

	pkg2 := serverpackets.NewMagicSkillUse(client, skill, ctrlPressed, shiftPressed)
	buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg2))

	pkg := serverpackets.SetupGauge(client)
	buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))
	client.Send(buffer.Bytes())
}
