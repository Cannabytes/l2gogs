package serverpackets

import (
	"l2gogameserver/data"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/packets"
)

func UserInfo(client *models.Client) []byte {
	player := client.CurrentChar
	stat := player.Stats

	buffer := packets.Get()
	defer packets.Put(buffer)

	x, y, z := client.CurrentChar.GetXYZ()

	buffer.WriteSingleByte(0x32)
	buffer.WriteD(x)
	buffer.WriteD(y)
	buffer.WriteD(z)

	buffer.WriteD(0) // Vehicle

	buffer.WriteD(player.ObjectID()) //objId

	buffer.WriteS(player.PlayerName()) //name //TODO

	buffer.WriteD(int32(player.Race)) //race ordinal //TODO
	buffer.WriteD(player.Sex)         //sex
	buffer.WriteD(player.BaseClass)   //baseClass

	buffer.WriteD(player.Level())                                                //level //TODO
	buffer.WriteQ(int64(player.Exp))                                             //exp
	buffer.WriteF(player.GetPercentFromCurrentLevel(player.Exp, player.Level())) //percent

	buffer.WriteD(int32(stat.STR)) //str
	buffer.WriteD(int32(stat.DEX)) //dex
	buffer.WriteD(int32(stat.CON)) //con
	buffer.WriteD(int32(stat.INT)) //int
	buffer.WriteD(int32(stat.WIT)) //wit
	buffer.WriteD(int32(stat.MEN)) //men

	buffer.WriteD(int32(player.MaxHP())) //Max hp //TODO

	buffer.WriteD(int32(player.CurHp))   //hp currnebt
	buffer.WriteD(int32(player.MaxMP())) //max mp
	buffer.WriteD(int32(player.CurMp))   //mp

	buffer.WriteD(player.Sp) //sp //TODO
	buffer.WriteD(0)         //currentLoad

	buffer.WriteD(109020) //maxLoad

	if player.IsActiveWeapon() {
		buffer.WriteD(40) //equiped weapon
	} else {
		buffer.WriteD(20) //no weapon
	}

	for _, slot := range models.GetPaperdollOrder() {
		buffer.WriteD(player.Paperdoll[slot].ObjId) //objId
	}
	for _, slot := range models.GetPaperdollOrder() {
		buffer.WriteD(int32(player.Paperdoll[slot].Id)) //itemId
	}
	for _, slot := range models.GetPaperdollOrder() {
		buffer.WriteD(int32(player.Paperdoll[slot].Enchant)) //enchant (страненько, на других сборках тут аргументация передается)
	}

	buffer.WriteD(0) //talisman slot
	buffer.WriteD(0) //Cloack

	buffer.WriteD(int32(player.PAtk()))         //patack //TODO
	buffer.WriteD(int32(player.AttackSpeed()))  //atackSpeed
	buffer.WriteD(int32(player.PDef()))         //pdef
	buffer.WriteD(33)                           //evasionRate
	buffer.WriteD(34)                           //accuracy //TODO
	buffer.WriteD(int32(player.CriticalRate())) //critHit
	buffer.WriteD(int32(player.MAtk()))         //Matack
	buffer.WriteD(int32(player.MAttackSpeed())) //M atackSpped

	buffer.WriteD(330) //patackSpeed again? //L22: Возможно это анимация скорости атаки

	buffer.WriteD(int32(player.MDef())) //mdef

	buffer.WriteD(player.PvpKills) //pvp
	buffer.WriteD(player.Karma)    //karma

	//_modspd := client.Player().GetMaxRunSpeed() * (1. / stat.BaseMoveSpd.Run)
	//logger.Info.Println(_modspd)
	buffer.WriteD(int32(player.MaxRunSpeed()))      //runSpeed
	buffer.WriteD(int32(stat.BaseMoveSpd.Walk))     //walkspeed
	buffer.WriteD(int32(stat.BaseMoveSpd.SlowSwim)) //swimRunSpeed
	buffer.WriteD(int32(stat.BaseMoveSpd.FastSwim)) //swimWalkSpeed
	buffer.WriteD(25)                               //flyRunSpeed
	buffer.WriteD(25)                               //flyWalkSpeed
	buffer.WriteD(25)                               //flyRunSpeed again
	buffer.WriteD(0)                                //flyWalkSpeed again
	buffer.WriteF(1)                                //moveMultipler

	buffer.WriteF(1.23) //atackSpeedMultiplier

	buffer.WriteF(8.0)  //collisionRadius
	buffer.WriteF(23.5) //collisionHeight

	buffer.WriteD(player.HairStyle) //hairStyle
	buffer.WriteD(player.HairColor) //hairColor
	buffer.WriteD(player.Face)      //face

	if player.IsAdmin() {
		buffer.WriteD(1) //IsGM?
	} else {
		buffer.WriteD(0) //IsPlayer
	}

	buffer.WriteS(player.Title()) //title

	buffer.WriteD(player.ClanId) //clanId
	buffer.WriteD(0)             //clancrestId
	buffer.WriteD(0)             //allyId
	buffer.WriteD(0)             //allyCrestId
	buffer.WriteD(0)             //RELATION CALCULATE ?

	buffer.WriteSingleByte(0)      //mountType
	buffer.WriteSingleByte(0)      //privateStoreType
	buffer.WriteSingleByte(0)      //hasDwarfCraft
	buffer.WriteD(player.PkKills)  //pk //TODO
	buffer.WriteD(player.PvpKills) //pvp //TODO

	buffer.WriteH(0) //cubic size
	//FOR cubicks

	buffer.WriteSingleByte(0) //PartyRoom

	buffer.WriteD(0) //EFFECTS

	buffer.WriteSingleByte(0) //WATER FLY EARTH

	buffer.WriteD(0) //clanBitmask

	buffer.WriteH(0) // c2 recommendations remaining
	buffer.WriteH(0) // c2 recommendations received //TODO

	buffer.WriteD(0) //npcMountId

	buffer.WriteH(player.GetInventoryLimit()) //inventoryLimit

	buffer.WriteD(player.ClassId) //	classId
	buffer.WriteD(0)              // special effects? circles around player...

	buffer.WriteD(int32(player.MaxCP())) //MaxCP
	buffer.WriteD(int32(player.CurCp))   //CurrentCp

	buffer.WriteSingleByte(0) //mounted air
	buffer.WriteSingleByte(0) //team Id

	buffer.WriteD(0) //ClanCrestLargeId

	buffer.WriteSingleByte(0) //isNoble
	buffer.WriteSingleByte(0) //isHero

	buffer.WriteSingleByte(0) //Fishing??
	buffer.WriteD(0)
	buffer.WriteD(0)
	buffer.WriteD(0)

	if player.IsAdmin() {
		buffer.WriteD(0x1a9112) //color name
	} else {
		var namecolor int32 = 0xffffff
		if player.NameColor != "" {
			namecolor = data.StrToInt32(player.NameColor)
		}
		buffer.WriteD(namecolor)
	}

	buffer.WriteSingleByte(1) //// changes the Speed display on Status Window

	buffer.WriteD(0) // changes the text above CP on Status Window
	buffer.WriteD(0) // plegue type

	if player.IsAdmin() {
		buffer.WriteD(0x6e071b) //titleColor
	} else {
		var titlecolor int32 = 0xffffff
		if player.TitleColor != "" {
			titlecolor = data.StrToInt32(player.TitleColor)
		}
		buffer.WriteD(titlecolor)
	}

	buffer.WriteD(0) // CursedWEAPON

	buffer.WriteD(0) //TransormDisplayId

	//attribute
	buffer.WriteH(-2) //attack element //TODO
	buffer.WriteH(0)  //attack elementValue
	buffer.WriteH(0)  //fire
	buffer.WriteH(0)  //water //TODO
	buffer.WriteH(0)  //wind //TODO
	buffer.WriteH(0)  //earth
	buffer.WriteH(0)  //holy
	buffer.WriteH(0)  //dark

	buffer.WriteD(0) //agationId

	buffer.WriteD(0)               //FAME //TODO
	buffer.WriteD(0)               //minimap or hellbound
	buffer.WriteD(player.Vitality) //vitaliti Point
	buffer.WriteD(0)               //abnormalEffects

	return buffer.Bytes()
}
