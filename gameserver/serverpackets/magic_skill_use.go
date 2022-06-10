package serverpackets

import (
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/packets"
)

func NewMagicSkillUse(clientI interfaces.ReciverAndSender, skill models.Skill, ctrlPressed, shiftPressed bool) []byte {
	client, ok := clientI.(*models.Client)
	if !ok {
		return []byte{}
	}

	buffer := packets.Get()
	defer packets.Put(buffer)

	client.CurrentChar.IsCastingNow = true
	client.CurrentChar.CurrentSkill = &models.SkillHolder{
		Skill:        skill,
		CtrlPressed:  ctrlPressed,
		ShiftPressed: shiftPressed,
	}

	var target int32
	//switch skill.TargetType {
	//case targets.AURA, targets.FRONT_AURA, targets.BEHIND_AURA, targets.GROUND, targets.SELF, targets.AURA_CORPSE_MOB, targets.COMMAND_CHANNEL, targets.AURA_FRIENDLY, targets.AURA_UNDEAD_ENEMY:
	//	target = 0
	//default:
	//	target = client.CurrentChar.Target
	//}

	// запускаем обработчик скилла
	_ = target

	/////////////////////////////////////////////////////////////////////////////////
	buffer.WriteSingleByte(0x48)
	buffer.WriteD(client.CurrentChar.ObjectID()) // activeChar id
	buffer.WriteD(client.CurrentChar.ObjectID()) // targetChar id
	buffer.WriteD(int32(skill.SkillId))          // skillId
	buffer.WriteD(int32(skill.Level))            // skillLevel
	buffer.WriteD(int32(skill.SkillHitTime))     // hitTime
	buffer.WriteD(int32(skill.ReuseDelay))       // reuseDelay

	x, y, z := client.CurrentChar.GetXYZ()
	buffer.WriteD(x)
	buffer.WriteD(y)
	buffer.WriteD(z)

	buffer.WriteH(0) //size???
	// for  by size ???

	buffer.WriteH(0) // _groundLocations.size()
	// for by _groundLocations.size()

	//location target
	buffer.WriteD(x)
	buffer.WriteD(y)
	buffer.WriteD(z)

	return buffer.Bytes()
}
