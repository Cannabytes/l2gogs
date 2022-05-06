package models

import (
	"context"
	"l2gogameserver/data"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/dto"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models/items"
	"l2gogameserver/gameserver/models/race"
	"l2gogameserver/gameserver/skills"
	"l2gogameserver/utils"
	"net"

	"sync"
	"time"
)

type (
	Character struct {
		Login       string
		ObjectId    int32
		CharName    string
		Level       int32
		MaxHp       float64
		CurHp       float64
		MaxMp       float64
		CurMp       float64
		MaxCp       float64
		CurCp       float64
		HpRegen     float64
		MpRegen     float64
		CpRegen     float64
		Face        int32
		HairStyle   int32
		HairColor   int32
		Sex         int32
		Coordinates *Coordinates
		Heading     int32
		Exp         int32
		Sp          int32
		Karma       int32
		PvpKills    int32
		PkKills     int32
		ClanId      int32
		Race        race.Race
		ClassId     int32
		BaseClass   int32
		Title       string
		OnlineTime  uint32
		Nobless     int32
		Vitality    int32
		IsAdmin     bool
		NameColor   string
		TitleColor  string

		CurrentRegion *WorldRegion
		Conn          *Client
		SockConn      net.Conn
		AttackEndTime int64
		// Paperdoll - массив всех слотов которые можно одеть
		Paperdoll       [26]MyItem
		Stats           StaticData
		pvpFlag         bool
		ShortCut        map[int32]dto.ShortCutDTO
		ActiveSoulShots []int32
		IsDead          bool
		IsFakeDeath     bool
		// Skills todo: проверить слайс или мапа лучше для скилов
		Skills                  []Skill
		SkillsItem              []Skill        //Скиллы, которые дает предметы, которые экиперованы на персонаже
		SkillsItemBonus         SkillItemBonus //Бонус предметов, которые добавляет статы, при экипировании
		IsCastingNow            bool
		SkillQueue              chan SkillHolder
		CurrentSkill            *SkillHolder // todo А может быть без * попробовать?
		Inventory               Inventory
		CursedWeaponEquippedId  int
		BonusStats              []items.ItemBonusStat
		ChannelUpdateShadowItem chan IUP
		InGame                  bool //Если игрок онлайн
		Target                  int32
		Macros                  []Macro
		CharInfoTo              chan []int32
		DeleteObjectTo          chan []int32
		NpcInfo                 chan []interfaces.Npcer
		IsMoving                bool
		Sit                     bool
		Buff                    []*BuffUser   //Баффы на персонаже
		BuffScheme              []*BuffScheme //Схемы баффов игрока
		OtherProperties         CharProperties
		Setting                 CharSetting
	}
	SkillItemBonus struct {
		MaxHP float64
		MaxMP float64
		MaxCP float64
		Speed float64 //Speed Run
	}
	BuffUser struct {
		Id     int //id skill
		Level  int //skill level
		Second int //Время баффа в секундах (обратный счет)
	}
	BuffScheme struct {
		Id     int
		CharId int
		Name   string
		Buffs  []BuffSchemeSkill
	}
	BuffSchemeSkill struct {
		SchemeId   int
		SkillId    int
		SkillLevel int
	}
	CharProperties struct {
		InventorySlot int32 //Кол-во слотов не постоянное, меняется в зависимости от определенных умений или статуса персонажа
		BuffSlot      int32 //Кол-во слотов не постоянное, меняется в зависимости от определенных умений
	}
	CharSetting struct {
		Language       int  // TODO:потом заменить на список
		EnableExp      bool //Получение опыта
		EnableSP       bool //Получение SP
		EnableAutoLoot bool //Автоматический подбор дропа

		AutoTradeParty  bool //Автоматически принимать трейд от членов пати
		AutoTradeClan   bool //Автоматически принимать трейд от членов клана
		AutoTradeSelfIP bool //Автоматически принимать трейд от всех с таким же IP
		AutoTradeAll    bool //Автоматически принимать трейд от всех игроков

		AutoPartyAll    bool //Автоматически принимать пати от всех
		AutoPartyClan   bool //Автоматически принимать пати от членов клана
		AutoPartySelfIP bool //Автоматически принимать пати от всех с таким же IP

		EnableSoulShotHalo bool //Вкл/Откл. сияния сосок

	}
	SkillHolder struct {
		Skill        Skill
		CtrlPressed  bool
		ShiftPressed bool
	}
	Coordinates struct {
		mu sync.Mutex
		X  int32
		Y  int32
		Z  int32
	}
	ToSendInfo struct {
		To   []int32
		Info utils.PacketByte
	}

	IUP struct {
		ObjId      int32
		UpdateType int16
	}
)

//Сбрасывает все бонусы скиллов предмета
func (c *Character) ResetSkillItemBonus() {
	c.SkillsItem = nil

	c.SkillsItemBonus.MaxHP = 0
	c.SkillsItemBonus.MaxMP = 0
	c.SkillsItemBonus.MaxCP = 0

	c.SkillsItemBonus.Speed = 0
}

func GetNewCharacterModel() *Character {
	character := new(Character)
	var sk []Skill
	character.Skills = sk
	character.ChannelUpdateShadowItem = make(chan IUP, 10)
	character.InGame = false
	return character
}

// SetSitStandPose Меняет положение персонажа от сидячего к стоячему и на оборот
//Возращает значение нового положения
func (c *Character) SetSitStandPose() int32 {
	if !c.Sit {
		c.Sit = true
		return 0
	}
	c.Sit = false
	return 1
}

func (c *Character) ListenSkillQueue() {
	for {
		select {
		case res := <-c.SkillQueue:
			logger.Info.Println("SKILL V QUEUE")
			logger.Info.Println(res.Skill.SkillId)
		}
	}
}

func (c *Character) SetSkillToQueue(skill Skill, ctrlPressed, shiftPressed bool) {
	s := SkillHolder{
		Skill:        skill,
		CtrlPressed:  ctrlPressed,
		ShiftPressed: shiftPressed,
	}
	c.SkillQueue <- s
}

// IsActiveWeapon есть ли у персонажа оружие в руках
func (c *Character) IsActiveWeapon() bool {
	x := c.Paperdoll[PAPERDOLL_RHAND]
	if x.ObjId != 0 {
		return true
	}
	return false
}

// GetPercentFromCurrentLevel получить % опыта на текущем уровне
func (c *Character) GetPercentFromCurrentLevel(exp, level int32) float64 {
	expPerLevel, expPerLevel2 := data.GetExpData(level)
	return float64(int64(exp)-expPerLevel) / float64(expPerLevel2-expPerLevel)
}

// GetBuffSkill Получение из БД всех сохраненных баффов
func GetBuffSkill(charId int32) []*BuffUser {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()
	var buffs []*BuffUser

	rows, err := dbConn.Query(context.Background(), "SELECT skill_id, level, second FROM character_buffs WHERE char_id = $1", charId)
	if err != nil {
		logger.Error.Panicln(err)
	}
	for rows.Next() {
		var buff BuffUser
		err = rows.Scan(&buff.Id, &buff.Level, &buff.Second)
		if err != nil {
			logger.Error.Panicln(err)
		}
		buffs = append(buffs, &buff)
	}
	return buffs
}

// Load загрузка персонажа
func (c *Character) Load() {
	c.InGame = true
	c.ShortCut = RestoreMe(c.ObjectId, c.ClassId)
	c.Buff = GetBuffSkill(c.ObjectId)
	c.LoadSkills()
	c.SkillQueue = make(chan SkillHolder)
	c.Inventory = GetMyItems(c.ObjectId)
	c.Paperdoll = c.LoadingVisibleInventory()

	c.LoadCharactersMacros()

	for _, v := range &c.Paperdoll {
		if v.ObjId != 0 {
			c.AddBonusStat(v.BonusStats)
		}
	}

	//HP/MP/CP/REGEN соответствующий уровню
	LvlUpgain := AllStats[int(c.ClassId)].LvlUpgainData[c.Level]
	c.MaxMp = LvlUpgain.Mp
	c.MaxHp = LvlUpgain.Hp
	c.MaxCp = LvlUpgain.Cp
	c.HpRegen = LvlUpgain.HpRegen
	c.MpRegen = LvlUpgain.MpRegen
	c.CpRegen = LvlUpgain.CpRegen

	c.ResetHpMpStatLevel()
	//Установка классовых статов
	//c.Stats = AllStats[int(c.ClassId)].StaticData //todo а для чего BaseClass ??
	c.GetRefreshStats()

	reg := GetRegion(c.Coordinates.X, c.Coordinates.Y, c.Coordinates.Z)
	c.CharInfoTo = make(chan []int32, 2)
	c.DeleteObjectTo = make(chan []int32, 2)
	c.NpcInfo = make(chan []interfaces.Npcer, 2)
	c.setWorldRegion(reg)

	reg.AddVisibleChar(c)

	go c.Shadow()
	go c.ListenSkillQueue()
	go c.checkRegion()
	go c.CounterTimeInGamePlayer()
}

// ResetHpMpStatLevel Установка значений на ХП,МП,ЦП, и реген по уровню
func (c *Character) ResetHpMpStatLevel() {
	LvlUpgain := AllStats[int(c.ClassId)].LvlUpgainData[c.Level]
	c.MaxMp = (LvlUpgain.Mp)
	c.MaxHp = (LvlUpgain.Hp)
	c.CurHp = (LvlUpgain.Cp)
	c.HpRegen = LvlUpgain.HpRegen
	c.MpRegen = LvlUpgain.MpRegen
	c.CpRegen = LvlUpgain.CpRegen
}

// CounterTimeInGamePlayer Время игрока нахождения в игре
func (c *Character) CounterTimeInGamePlayer() {
	for {
		if c.InGame == false {
			logger.Info.Println("Счетчик времени в игре остановлен")
			return
		}
		c.OnlineTime++
		time.Sleep(time.Second)
	}
}

func (c *Character) Shadow() {
	for {
		for i := range c.Inventory.Items {
			v := c.Inventory.Items[i]
			if v.Item.Durability > 0 && v.Loc == PaperdollLoc {
				var iup IUP
				iup.ObjId = v.ObjId
				switch c.Inventory.Items[i].Mana {

				case 0:
					iup.UpdateType = UpdateTypeRemove
					c.ChannelUpdateShadowItem <- iup
					DeleteItem(v, c)
				default:
					c.Inventory.Items[i].Mana -= 1
					iup.UpdateType = UpdateTypeModify
					c.ChannelUpdateShadowItem <- iup
				}

			}
		}

		time.Sleep(time.Second)
	}

}

func (c *Character) checkSoulShot() {
	if len(c.ActiveSoulShots) == 0 {
		return
	}
}

func (c *Character) IsCursedWeaponEquipped() bool {
	return c.CursedWeaponEquippedId != 0
}

func (c *Character) AddBonusStat(s []items.ItemBonusStat) {
	c.BonusStats = append(c.BonusStats, s...)
}

func (c *Character) RemoveBonusStat(s []items.ItemBonusStat) {
	//for i,v := range c.BonusStats {
	//	for _,vv := range s {
	//		if v == vv {
	//			c.BonusStats[i] = c.BonusStats[len(c.BonusStats)-1] //todo переделать на безопасный вариант ) или еще что нить придумать
	//			c.BonusStats = c.BonusStats[:len(c.BonusStats)-1]
	//		}
	//	}
	//
	//}

	news := make([]items.ItemBonusStat, 0, len(c.BonusStats))
	for _, v := range c.BonusStats {
		flag := false
		for _, vv := range s {
			if v == vv {
				flag = true
				break
			}
		}
		if !flag {
			news = append(news, v)
		}
	}
	c.BonusStats = news
}

func (c *Character) AddBonusSkill(s Skill) {
	c.SkillsItem = append(c.SkillsItem, s)
}

// Получение списка скиллов надетых предметов
// Возвращается true если скиллы были изменены
func (c *Character) SkillItemListRefresh() bool {
	c.ResetSkillItemBonus()
	for _, selectedItem := range c.Paperdoll {
		if selectedItem.ObjId != 0 {
			itemSkill := selectedItem.ItemSkill
			skill, ok := GetSkillName(itemSkill)
			if ok {
				c.AddBonusSkill(skill)
				c.BonusStatCalsSkills(skill)
			}
		}
	}
	if c.SkillsItem != nil {
		return true
	}
	return false
}

func (c *Character) BonusStatCalsSkills(skill Skill) {
	effect := skill.Effect
	if effect.PMaxMp != nil {
		c.SkillsItemBonus.MaxMP = float64(skills.CapMath(c.MaxMp, skill.Effect.PMaxMp.Val, effect.PMaxMp.Cap))
	}
	if effect.PSpeed != nil {
		c.SkillsItemBonus.Speed = float64(skills.CapMath(c.Stats.BaseMoveSpd.Run, skill.Effect.PSpeed.Val, effect.PSpeed.Cap))
	}

}

// GetRefreshStats Обновление статов персонажа
// Берется статы из оружия, брони.
// TODO: Скиллы и бижа не учитывается
func (c *Character) GetRefreshStats() {
	c.ResetBonusStatCals()
	c.BonusStats = nil
	for _, v := range &c.Paperdoll {
		if v.ObjId != 0 {
			c.AddBonusStat(v.BonusStats)
			c.BonusStatCals(v)
		}
	}

}

// Кол-во п.дефа который дает весь надетый шмот
//func (c *Character) GetPDefEquip() int32 {
//var base float64
////if c.Paperdoll[PAPERDOLL_FEET].IsEquipped() == 1 {
////	base = float64(c.Stats.BasePDef.Feet)
////}
////if c.Paperdoll[PAPERDOLL_CHEST].IsEquipped() == 1 {
////	base += float64(c.Stats.BasePDef.Chest)
////}
////if c.Paperdoll[PAPERDOLL_CLOAK].IsEquipped() == 1 {
////	base += float64(c.Stats.BasePDef.Cloak)
////}
////if c.Paperdoll[PAPERDOLL_HEAD].IsEquipped() == 1 {
////	base += float64(c.Stats.BasePDef.Head)
////}
////if c.Paperdoll[PAPERDOLL_GLOVES].IsEquipped() == 1 {
////	base += float64(c.Stats.BasePDef.Gloves)
////}
////if c.Paperdoll[PAPERDOLL_LEGS].IsEquipped() == 1 {
////	base += float64(c.Stats.BasePDef.Legs)
////}
////if c.Paperdoll[PAPERDOLL_UNDER].IsEquipped() == 1 {
////	base += float64(c.Stats.BasePDef.Underwear)
////}
//
//for _, v := range c.BonusStats {
//	if v.Type == "physical_defense" {
//		base += v.Val
//	}
//}
//
//base *= float64(c.Level+89) / 100
//
//return int32(base)
//}

// Сложение всех статов
func (c *Character) BonusStatCals(item MyItem) {
	for _, v := range item.BonusStats {
		if v.Type == "physical_damage" {
			c.Stats.BasePAtk += int(v.Val)
		}
		if v.Type == "physical_defense" {
			c.Stats.BasePDef.PDef += int(v.Val)
		}
		if v.Type == "magical_damage" {
			c.Stats.BaseMAtk += int(v.Val)
		}
		if v.Type == "magical_defense" {
			c.Stats.BaseMDef.MDef += int(v.Val)
		}
		if v.Type == "critical" {
			c.Stats.BaseCritRate += int(v.Val)
		}
		if v.Type == "attack_speed" {
			c.Stats.BasePAtkSpd += int(v.Val)
		}
		if v.Type == "mp_bonus" {
			c.MaxMp += (v.Val)
		}
	}

}

func (c *Character) ResetBonusStatCals() {
	c.Stats = AllStats[int(c.ClassId)].StaticData
}

func (c *Character) GetInventoryLimit() int16 {
	if c.Race == race.DWARF {
		return 100
	}
	return 80
}

func (c *Character) setWorldRegion(newRegion interfaces.WorldRegioner) {
	var oldAreas []interfaces.WorldRegioner

	currReg := c.GetCurrentRegion().(*WorldRegion)
	if currReg != nil {
		c.CurrentRegion.DeleteVisibleChar(c)
		oldAreas = currReg.GetNeighbors()
	}

	var newAreas []interfaces.WorldRegioner
	if newRegion != nil {
		newRegion.AddVisibleChar(c)
		newAreas = newRegion.GetNeighbors()
	}

	// кому отправить charInfo
	deleteObjectPkgTo := make([]int32, 0, 64)
	for _, region := range oldAreas {
		if !Contains(newAreas, region) {

			for _, v := range region.GetCharsInRegion() {
				if v.GetObjectId() == c.GetObjectId() {
					continue
				}
				deleteObjectPkgTo = append(deleteObjectPkgTo, v.GetObjectId())
			}
		}
	}
	if len(deleteObjectPkgTo) > 0 {
		c.DeleteObjectTo <- deleteObjectPkgTo
	}

	// кому отправить charInfo
	charInfoPkgTo := make([]int32, 0, 64)
	npcPkgTo := make([]interfaces.Npcer, 0, 64)
	for _, region := range newAreas {
		if !Contains(oldAreas, region) {
			for _, v := range region.GetCharsInRegion() {
				if v.GetObjectId() == c.GetObjectId() {
					continue
				}
				charInfoPkgTo = append(charInfoPkgTo, v.GetObjectId())
			}

			npcPkgTo = append(npcPkgTo, region.GetNpcInRegion()...)

		}
	}
	if len(charInfoPkgTo) > 0 {
		c.CharInfoTo <- charInfoPkgTo
	}
	c.CurrentRegion = newRegion.(*WorldRegion)

	if len(npcPkgTo) > 0 {
		c.NpcInfo <- npcPkgTo
	}

}

func (c *Character) checkRegion() {
	for {
		time.Sleep(time.Second)
		if c.CurrentRegion != nil {
			curReg := c.CurrentRegion
			x, y, z := c.GetXYZ()
			ncurReg := GetRegion(x, y, z)
			if curReg != ncurReg {
				c.setWorldRegion(ncurReg)
			}
		}
	}

}

//ExistItemInInventory Возвращает ссылку на Item если он есть в инвентаре
func (c *Character) ExistItemInInventory(objectItemId int32) *MyItem {
	for i := range c.Inventory.Items {
		item := c.Inventory.Items[i]
		if item.ObjId == objectItemId {
			return item
		}
	}
	return nil
}

func (c *Character) GetObjectId() int32 {
	return c.ObjectId
}
func (c *Character) GetName() string {
	return c.CharName
}

//Возвращает общее кол-во ХП (дающее и с скиллы, с предметы и баффы)
func (c *Character) GetMaxHP() float64 {
	return c.MaxHp + c.SkillsItemBonus.MaxHP
}

func (c *Character) GetMaxMP() float64 {
	return c.MaxMp + c.SkillsItemBonus.MaxMP
}

func (c *Character) GetMaxCP() float64 {
	return c.MaxCp + c.SkillsItemBonus.MaxCP
}

//Возвращает скорость бега персонажа учитывая все баффы и .тд.
func (c *Character) GetMaxRunSpeed() float64 {
	return c.Stats.BaseMoveSpd.Run + c.SkillsItemBonus.Speed
}

func (c *Character) GetBuff() []*BuffUser {
	return c.Buff
}
func (c *Character) SetStatusOffline() {
	c.InGame = false
}
func (c *Character) SetX(x int32) {
	c.Coordinates.X = x
}
func (c *Character) SetY(y int32) {
	c.Coordinates.Y = y
}
func (c *Character) SetZ(z int32) {
	c.Coordinates.Z = z
}
func (c *Character) SetXYZ(x, y, z int32) {
	c.Coordinates.X = x
	c.Coordinates.Y = y
	c.Coordinates.Z = z
}
func (c *Character) SetHeading(h int32) {
	c.Heading = h
}
func (c *Character) SetInstanceId(i int32) {
	_ = i
	//TODO release
}
func (c *Character) GetXYZ() (x, y, z int32) {
	return c.Coordinates.X, c.Coordinates.Y, c.Coordinates.Z
}
func (c *Character) GetX() int32 {
	return c.Coordinates.X
}
func (c *Character) GetY() int32 {
	return c.Coordinates.Y
}
func (c *Character) GetZ() int32 {
	return c.Coordinates.Z
}

func (c *Character) EncryptAndSend(data []byte) {
	c.Conn.EncryptAndSend(data)
}
func (c *Character) GetCurrentRegion() interfaces.WorldRegioner {
	return c.CurrentRegion
}

func (c *Character) CloseChannels() {
	c.ChannelUpdateShadowItem = nil
	c.NpcInfo = nil
	c.CharInfoTo = nil
	c.DeleteObjectTo = nil
}

func (c *Character) GetClassId() int32 {
	return c.ClassId
}

func (c *Character) GetSkillInfo(skill_id int) Skill {
	for _, skill := range c.Skills {
		if skill.SkillId == skill_id {
			return skill
		}
	}
	return Skill{}
}
