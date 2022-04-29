package serverpackets

import (
	"l2gogameserver/data"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/packets"
)

func UserInfo(clientI interfaces.CharacterI) []byte {
	character, ok := clientI.(*models.Character)
	if !ok {
		return []byte{}
	}
	character.GetRefreshStats()
	stat := &character.Stats

	buffer := packets.Get()
	defer packets.Put(buffer)

	x, y, z := character.GetXYZ()

	buffer.WriteSingleByte(0x32)
	buffer.WriteD(x)
	buffer.WriteD(y)
	buffer.WriteD(z)

	buffer.WriteD(0) // Vehicle

	buffer.WriteD(character.ObjectId) //objId

	buffer.WriteS(character.CharName) //name //TODO

	buffer.WriteD(int32(character.Race)) //race ordinal //TODO
	buffer.WriteD(character.Sex)         //sex
	buffer.WriteD(character.BaseClass)   //baseClass

	buffer.WriteD(character.Level)                                                      //level //TODO
	buffer.WriteQ(int64(character.Exp))                                                 //exp
	buffer.WriteF(character.GetPercentFromCurrentLevel(character.Exp, character.Level)) //percent

	buffer.WriteD(int32(stat.STR)) //str
	buffer.WriteD(int32(stat.DEX)) //dex
	buffer.WriteD(int32(stat.CON)) //con
	buffer.WriteD(int32(stat.INT)) //int
	buffer.WriteD(int32(stat.WIT)) //wit
	buffer.WriteD(int32(stat.MEN)) //men

	buffer.WriteD(character.MaxHp) //Max hp //TODO

	buffer.WriteD(character.CurHp) //hp currnebt

	buffer.WriteD(character.MaxMp) //max mp
	buffer.WriteD(character.CurMp) //mp

	buffer.WriteD(character.Sp) //sp //TODO
	buffer.WriteD(0)            //currentLoad

	buffer.WriteD(109020) //maxLoad

	if character.IsActiveWeapon() {
		buffer.WriteD(20) //no weapon
	} else {
		buffer.WriteD(40) //equiped weapon
	}

	for _, slot := range models.GetPaperdollOrder() {
		buffer.WriteD(character.Paperdoll[slot].ObjId) //objId
	}
	for _, slot := range models.GetPaperdollOrder() {
		buffer.WriteD(int32(character.Paperdoll[slot].Id)) //itemId
	}
	for _, slot := range models.GetPaperdollOrder() {
		buffer.WriteD(int32(character.Paperdoll[slot].Enchant)) //enchant
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
	buffer.WriteD(25)                                                                                                                                                                      //M atackSpped

	buffer.WriteD(330) //patackSpeed again?

	buffer.WriteD(47) //mdef

	buffer.WriteD(character.PvpKills) //pvp
	buffer.WriteD(character.Karma)    //karma

	buffer.WriteD(999) //runSpeed
	buffer.WriteD(80)  //walkspeed
	buffer.WriteD(50)  //swimRunSpeed
	buffer.WriteD(50)  //swimWalkSpeed
	buffer.WriteD(0)   //flyRunSpeed
	buffer.WriteD(0)   //flyWalkSpeed
	buffer.WriteD(0)   //flyRunSpeed again
	buffer.WriteD(0)   //flyWalkSpeed again
	buffer.WriteF(1.1) //moveMultipler

	buffer.WriteF(1.21) //atackSpeedMultiplier

	buffer.WriteF(8.0)  //collisionRadius
	buffer.WriteF(23.5) //collisionHeight

	buffer.WriteD(character.HairStyle) //hairStyle
	buffer.WriteD(character.HairColor) //hairColor
	buffer.WriteD(character.Face)      //face

	buffer.WriteD(1) //IsGM?

	buffer.WriteS(character.Title) //title

	buffer.WriteD(character.ClanId) //clanId
	buffer.WriteD(0)                //clancrestId
	buffer.WriteD(0)                //allyId
	buffer.WriteD(0)                //allyCrestId
	buffer.WriteD(0)                //RELATION CALCULATE ?

	buffer.WriteSingleByte(0)         //mountType
	buffer.WriteSingleByte(0)         //privateStoreType
	buffer.WriteSingleByte(0)         //hasDwarfCraft
	buffer.WriteD(character.PkKills)  //pk //TODO
	buffer.WriteD(character.PvpKills) //pvp //TODO

	buffer.WriteH(0) //cubic size
	//FOR cubicks

	buffer.WriteSingleByte(0) //PartyRoom

	buffer.WriteD(0) //EFFECTS

	buffer.WriteSingleByte(0) //WATER FLY EARTH

	buffer.WriteD(0) //clanBitmask

	buffer.WriteH(0) // c2 recommendations remaining
	buffer.WriteH(0) // c2 recommendations received //TODO

	buffer.WriteD(0) //npcMountId

	buffer.WriteH(character.GetInventoryLimit()) //inventoryLimit

	buffer.WriteD(character.ClassId) //	classId
	buffer.WriteD(0)                 // special effects? circles around player...

	buffer.WriteD(character.MaxCp) //MaxCP
	buffer.WriteD(character.CurCp) //CurrentCp

	buffer.WriteSingleByte(0) //mounted air
	buffer.WriteSingleByte(0) //team Id

	buffer.WriteD(0) //ClanCrestLargeId

	buffer.WriteSingleByte(0) //isNoble
	buffer.WriteSingleByte(0) //isHero

	buffer.WriteSingleByte(0) //Fishing??
	buffer.WriteD(0)
	buffer.WriteD(0)
	buffer.WriteD(0)

	buffer.WriteD(16777215)

	buffer.WriteSingleByte(1) //// changes the Speed display on Status Window

	buffer.WriteD(0) // changes the text above CP on Status Window
	buffer.WriteD(0) // plegue type

	title := (255 & 0xFF) + ((168 & 0xFF) << 8) + ((00 & 0xFF) << 16)
	buffer.WriteD(int32(title)) //titleColor

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

	buffer.WriteD(0)                  //FAME //TODO
	buffer.WriteD(0)                  //minimap or hellbound
	buffer.WriteD(character.Vitality) //vitaliti Point
	buffer.WriteD(0)                  //abnormalEffects

	return buffer.Bytes()
}
