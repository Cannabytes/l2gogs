package serverpackets

import (
	"l2gogameserver/data"
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/packets"
)

func UserInfo(client *models.Client) []byte {

	stat := client.CurrentChar.Stats

	buffer := packets.Get()
	defer packets.Put(buffer)

	x, y, z := client.CurrentChar.GetXYZ()

	buffer.WriteSingleByte(0x32)
	buffer.WriteD(x)
	buffer.WriteD(y)
	buffer.WriteD(z)

	buffer.WriteD(0) // Vehicle

	buffer.WriteD(client.CurrentChar.ObjectId) //objId

	buffer.WriteS(client.CurrentChar.CharName) //name //TODO

	buffer.WriteD(int32(client.CurrentChar.Race)) //race ordinal //TODO
	buffer.WriteD(client.CurrentChar.Sex)         //sex
	buffer.WriteD(client.CurrentChar.BaseClass)   //baseClass

	buffer.WriteD(client.CurrentChar.Level)                                                                        //level //TODO
	buffer.WriteQ(int64(client.CurrentChar.Exp))                                                                   //exp
	buffer.WriteF(client.CurrentChar.GetPercentFromCurrentLevel(client.CurrentChar.Exp, client.CurrentChar.Level)) //percent

	buffer.WriteD(int32(stat.STR)) //str
	buffer.WriteD(int32(stat.DEX)) //dex
	buffer.WriteD(int32(stat.CON)) //con
	buffer.WriteD(int32(stat.INT)) //int
	buffer.WriteD(int32(stat.WIT)) //wit
	buffer.WriteD(int32(stat.MEN)) //men

	buffer.WriteD(int32(client.CurrentChar.GetMaxHP())) //Max hp //TODO

	buffer.WriteD(int32(client.CurrentChar.CurHp))      //hp currnebt
	buffer.WriteD(int32(client.CurrentChar.GetMaxMP())) //max mp
	buffer.WriteD(int32(client.CurrentChar.CurMp))      //mp

	buffer.WriteD(client.CurrentChar.Sp) //sp //TODO
	buffer.WriteD(0)                     //currentLoad

	buffer.WriteD(109020) //maxLoad

	if client.CurrentChar.IsActiveWeapon() {
		buffer.WriteD(40) //equiped weapon
	} else {
		buffer.WriteD(20) //no weapon
	}

	for _, slot := range models.GetPaperdollOrder() {
		buffer.WriteD(client.CurrentChar.Paperdoll[slot].ObjId) //objId
	}
	for _, slot := range models.GetPaperdollOrder() {
		buffer.WriteD(int32(client.CurrentChar.Paperdoll[slot].Id)) //itemId
	}
	for _, slot := range models.GetPaperdollOrder() {
		buffer.WriteD(int32(client.CurrentChar.Paperdoll[slot].Enchant)) //enchant (страненько, на других сборках тут аргументация передается)
	}

	buffer.WriteD(0) //talisman slot
	buffer.WriteD(0) //Cloack

	buffer.WriteD(int32(stat.BasePAtk))                                                                                                                                                    //patack //TODO
	buffer.WriteD(int32(stat.BasePAtkSpd))                                                                                                                                                 //atackSpeed
	buffer.WriteD(int32(data.CalcInt(stat.BasePDef.PDef, stat.BasePDef.Gloves, stat.BasePDef.Legs, stat.BasePDef.Underwear, stat.BasePDef.Feet, stat.BasePDef.Chest, stat.BasePDef.Head))) //pdef
	buffer.WriteD(33)                                                                                                                                                                      //evasionRate
	buffer.WriteD(34)                                                                                                                                                                      //accuracy //TODO
	buffer.WriteD(int32(stat.BaseCritRate))                                                                                                                                                //critHit
	buffer.WriteD(int32(stat.BaseMAtk))                                                                                                                                                    //Matack
	buffer.WriteD(333)                                                                                                                                                                     //M atackSpped

	buffer.WriteD(330) //patackSpeed again?

	buffer.WriteD(int32(data.CalcInt(stat.BaseMDef.MDef, stat.BaseMDef.Lear, stat.BaseMDef.Neck, stat.BaseMDef.Lfinger, stat.BaseMDef.Rear, stat.BaseMDef.Rfinger))) //mdef

	buffer.WriteD(client.CurrentChar.PvpKills) //pvp
	buffer.WriteD(client.CurrentChar.Karma)    //karma

	_modspd := client.GetCurrentChar().GetMaxRunSpeed() * (1. / stat.BaseMoveSpd.Run)
	logger.Info.Println(_modspd)
	buffer.WriteD(int32(client.GetCurrentChar().GetMaxRunSpeed() / _modspd)) //runSpeed
	buffer.WriteD(int32(stat.BaseMoveSpd.Walk))                              //walkspeed
	buffer.WriteD(int32(stat.BaseMoveSpd.SlowSwim))                          //swimRunSpeed
	buffer.WriteD(int32(stat.BaseMoveSpd.FastSwim))                          //swimWalkSpeed
	buffer.WriteD(25)                                                        //flyRunSpeed
	buffer.WriteD(25)                                                        //flyWalkSpeed
	buffer.WriteD(25)                                                        //flyRunSpeed again
	buffer.WriteD(0)                                                         //flyWalkSpeed again
	buffer.WriteF(_modspd)                                                   //moveMultipler

	buffer.WriteF(1.23) //atackSpeedMultiplier

	buffer.WriteF(8.0)  //collisionRadius
	buffer.WriteF(23.5) //collisionHeight

	buffer.WriteD(client.CurrentChar.HairStyle) //hairStyle
	buffer.WriteD(client.CurrentChar.HairColor) //hairColor
	buffer.WriteD(client.CurrentChar.Face)      //face

	if client.CurrentChar.IsAdmin {
		buffer.WriteD(1) //IsGM?
	} else {
		buffer.WriteD(0) //IsPlayer
	}

	buffer.WriteS(client.CurrentChar.Title) //title

	buffer.WriteD(client.CurrentChar.ClanId) //clanId
	buffer.WriteD(0)                         //clancrestId
	buffer.WriteD(0)                         //allyId
	buffer.WriteD(0)                         //allyCrestId
	buffer.WriteD(0)                         //RELATION CALCULATE ?

	buffer.WriteSingleByte(0)                  //mountType
	buffer.WriteSingleByte(0)                  //privateStoreType
	buffer.WriteSingleByte(0)                  //hasDwarfCraft
	buffer.WriteD(client.CurrentChar.PkKills)  //pk //TODO
	buffer.WriteD(client.CurrentChar.PvpKills) //pvp //TODO

	buffer.WriteH(0) //cubic size
	//FOR cubicks

	buffer.WriteSingleByte(0) //PartyRoom

	buffer.WriteD(0) //EFFECTS

	buffer.WriteSingleByte(0) //WATER FLY EARTH

	buffer.WriteD(0) //clanBitmask

	buffer.WriteH(0) // c2 recommendations remaining
	buffer.WriteH(0) // c2 recommendations received //TODO

	buffer.WriteD(0) //npcMountId

	buffer.WriteH(client.CurrentChar.GetInventoryLimit()) //inventoryLimit

	buffer.WriteD(client.CurrentChar.ClassId) //	classId
	buffer.WriteD(0)                          // special effects? circles around player...

	buffer.WriteD(int32(client.CurrentChar.GetMaxCP())) //MaxCP
	buffer.WriteD(int32(client.CurrentChar.CurCp))      //CurrentCp

	buffer.WriteSingleByte(0) //mounted air
	buffer.WriteSingleByte(0) //team Id

	buffer.WriteD(0) //ClanCrestLargeId

	buffer.WriteSingleByte(0) //isNoble
	buffer.WriteSingleByte(0) //isHero

	buffer.WriteSingleByte(0) //Fishing??
	buffer.WriteD(0)
	buffer.WriteD(0)
	buffer.WriteD(0)

	if client.CurrentChar.IsAdmin {
		buffer.WriteD(0x1a9112) //color name
	} else {
		var namecolor int32 = 0xffffff
		if client.CurrentChar.NameColor != "" {
			namecolor = data.StrToInt32(client.CurrentChar.NameColor)
		}
		buffer.WriteD(namecolor)
	}

	buffer.WriteSingleByte(1) //// changes the Speed display on Status Window

	buffer.WriteD(0) // changes the text above CP on Status Window
	buffer.WriteD(0) // plegue type

	if client.CurrentChar.IsAdmin {
		buffer.WriteD(0x6e071b) //titleColor
	} else {
		var titlecolor int32 = 0xffffff
		if client.CurrentChar.TitleColor != "" {
			titlecolor = data.StrToInt32(client.CurrentChar.TitleColor)
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

	buffer.WriteD(0)                           //FAME //TODO
	buffer.WriteD(0)                           //minimap or hellbound
	buffer.WriteD(client.CurrentChar.Vitality) //vitaliti Point
	buffer.WriteD(0)                           //abnormalEffects

	return buffer.Bytes()
}
