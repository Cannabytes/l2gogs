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
		login       string
		objectId    int32
		playerName  string
		level       int32
		hp          float64
		CurHp       float64
		mp          float64
		CurMp       float64
		cp          float64
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
		exp         int32
		Sp          int32
		Karma       int32
		PvpKills    int32
		PkKills     int32
		ClanId      int32
		Race        race.Race
		classID     int32
		baseClassID int32
		title       string
		onlineTime  uint32
		Nobless     int32
		Vitality    int32
		isAdmin     bool
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
		SkillsItem              []Skill    //Скиллы, которые дает предметы, которые экиперованы на персонаже
		SkillsItemBonus         SkillBonus //Бонус предметов, которые добавляет статы, при экипировании
		SkillsBuffBonus         SkillBonus //Бонус предметов, которые добавляет статы, при  баффе
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
		sit                     bool
		buff                    []*BuffUser   //Баффы на персонаже
		BuffScheme              []*BuffScheme //Схемы баффов игрока
		OtherProperties         CharProperties
		Setting                 CharSetting
	}
	SkillBonus struct {
		MaxHP        float64
		MaxMP        float64
		MaxCP        float64
		Speed        float64 //Speed Run
		PDef         float64
		MDef         float64
		PAtk         float64
		MAtk         float64
		AtkSpd       float64
		MAtkSpd      float64
		CriticalRate float64
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

// HP Возвращает HP персонажа
func (c *Character) HP() float64 {
	return c.hp
}

// SetHP Установить HP
func (c *Character) SetHP(hp float64) {
	c.hp = hp
}

// MP Возвращает MP персонажа
func (c *Character) MP() float64 {
	return c.mp
}

// SetMP Установить MP
func (c *Character) SetMP(mp float64) {
	c.mp = mp
}

// SetMP Установить дополнительно MP
func (c *Character) SetAddMP(mp float64) {
	c.mp += mp
}

// CP Возвращает CP персонажа
func (c *Character) CP() float64 {
	return c.cp
}

// SetCP Установить CP
func (c *Character) SetCP(cp float64) {
	c.cp = cp
}

// EXP Возвращает EXP
func (c *Character) EXP() int32 {
	return c.exp
}

// ExpAdd Добавляет к Exp ещё nExp
func (c *Character) ExpAdd(nExp int32) {
	c.exp += nExp
}

// SetExp Установить новый Exp
func (c *Character) SetExp(nExp int32) {
	c.exp = nExp
}

// BaseClassID Вернуть базовый класс игрока
func (c *Character) BaseClassID() int32 {
	return c.baseClassID
}

// SetBaseClassID Установить базовый класс игрока
func (c *Character) SetBaseClassID(classID int) {
	c.baseClassID = int32(classID)
}

// ClassID Вернуть актуальный класс игрока
func (c *Character) ClassID() int32 {
	return c.classID
}

// SetClassID Установить класс игрока
func (c *Character) SetClassID(classID int) {
	c.classID = int32(classID)
}

// OnlineTime Время проведенное игроком в игре
func (c *Character) OnlineTime() uint32 {
	return c.onlineTime
}

// SetOnlineTime Добавление времени (в сек) к времени персонажа в игре
func (c *Character) SetOnlineTimeIncrement(nSecond int) {
	c.onlineTime += uint32(nSecond)
}

// SetOnlineTime Установить значение времени персонажа в игре
func (c *Character) SetOnlineTime(nSecond int) {
	c.onlineTime = uint32(nSecond)
}

// ID персонажа
func (c *Character) ObjectID() int32 {
	return c.objectId
}

// Установить ID персонажа
func (c *Character) SetObjectID(objectID int) {
	c.objectId = int32(objectID)
}

// Название аккаунта (Login)
func (c *Character) AccountName() string {
	return c.login
}

// SetAccountName Установить название аккаунта
func (c *Character) SetAccountName(login string) {
	c.login = login
}

// PlayerName Ник персонажа
func (c *Character) PlayerName() string {
	return c.playerName
}

// SetPlayerName Установить ник персонажу
func (c *Character) SetPlayerName(name string) {
	c.playerName = name
}

// Level Уровень игрока
func (c *Character) Level() int32 {
	return c.level
}

// SetLevel Установить уровень игрока
func (c *Character) SetLevel(level int) {
	c.level = int32(level)
}

// ResetSkillItemBonus Сбрасывает все бонусы скиллов предмета
func (c *Character) ResetSkillItemBonus() {
	c.SkillsItem = nil
	c.SkillsItemBonus = SkillBonus{
		MaxHP:        0,
		MaxMP:        0,
		MaxCP:        0,
		Speed:        0,
		PDef:         0,
		MDef:         0,
		PAtk:         0,
		MAtk:         0,
		AtkSpd:       0,
		MAtkSpd:      0,
		CriticalRate: 0,
	}
}

// Title Титул игрока
func (c *Character) Title() string {
	return c.title
}

// SetTitle Установить титул игроку
func (c *Character) SetTitle(title string) {
	c.title = title
}

// ResetSkillBuffBonus Сбрасывает все бонусы скиллов баффа
// Обязательно нужно перечислить все статы, которые необходимо сбросить
func (c *Character) resetSkillBuffBonus() {
	c.SkillsBuffBonus = SkillBonus{
		MaxHP: 0,
		MaxMP: 0,
		MaxCP: 0,
		Speed: 0,
		PDef:  0,
	}
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
	if !c.sit {
		c.sit = true
		return 0
	}
	c.sit = false
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
	c.ShortCut = RestoreMe(c.ObjectID(), c.ClassID())
	c.SetBuff(GetBuffSkill(c.ObjectID()))
	c.LoadSkills()
	c.SkillQueue = make(chan SkillHolder)
	c.Inventory = GetMyItems(c.ObjectID())
	c.Paperdoll = c.LoadingVisibleInventory()

	c.LoadCharactersMacros()

	for _, v := range &c.Paperdoll {
		if v.ObjId != 0 {
			c.AddBonusStat(v.BonusStats)
		}
	}

	//HP/MP/CP/REGEN соответствующий уровню
	LvlUpgain := AllStats[int(c.ClassID())].LvlUpgainData[c.Level()]
	c.SetHP(LvlUpgain.Hp)
	c.SetMP(LvlUpgain.Mp)
	c.SetCP(LvlUpgain.Cp)
	c.HpRegen = LvlUpgain.HpRegen
	c.MpRegen = LvlUpgain.MpRegen
	c.CpRegen = LvlUpgain.CpRegen

	c.ResetHpMpStatLevel()
	c.StatsRefresh()

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
	LvlUpgain := AllStats[int(c.ClassID())].LvlUpgainData[c.Level()]
	c.SetHP(LvlUpgain.Hp)
	c.SetMP(LvlUpgain.Mp)
	c.SetCP(LvlUpgain.Cp)
	c.CurHp = LvlUpgain.Cp
	c.HpRegen = LvlUpgain.HpRegen
	c.MpRegen = LvlUpgain.MpRegen
	c.CpRegen = LvlUpgain.CpRegen
}

// IsAdmin Имеет ли игрок статус админа
func (c *Character) IsAdmin() bool {
	return c.isAdmin
}

// SetAdmin Применить к пользователю статус админа
func (c *Character) SetAdmin(privilege bool) {
	c.isAdmin = privilege
}

// CounterTimeInGamePlayer Время игрока нахождения в игре
func (c *Character) CounterTimeInGamePlayer() {
	for {
		if c.InGame == false {
			logger.Info.Println("Счетчик времени в игре остановлен")
			return
		}
		c.SetOnlineTimeIncrement(1)
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

// SkillItemListRefresh Получение списка скиллов надетых предметов
// Возвращается true если скиллы были изменены
func (c *Character) SkillItemListRefresh() bool {
	c.ResetSkillItemBonus()
	for _, selectedItem := range c.Paperdoll {
		if selectedItem.ObjId != 0 {
			itemSkill := selectedItem.ItemSkill
			skill, ok := GetSkillName(itemSkill)
			if ok {
				c.AddBonusSkill(skill)
				c.bonusStatCalsSkills(skill)
			}
		}
	}
	if c.SkillsItem != nil {
		return true
	}
	return false
}

// Добавление бонусных статов от скиллов, которые надеты на персонажа
func (c *Character) bonusStatCalsSkills(skill Skill) {
	effect := skill.Effect
	if effect.PMaxHp != nil {
		c.SkillsItemBonus.MaxHP = float64(skills.CapMath(c.MaxHP(), skill.Effect.PMaxHp.Val, effect.PMaxHp.Cap))
	}
	if effect.PMaxMp != nil {
		c.SkillsItemBonus.MaxMP = float64(skills.CapMath(c.MaxMP(), skill.Effect.PMaxMp.Val, effect.PMaxMp.Cap))
	}
	if effect.PSpeed != nil {
		c.SkillsItemBonus.Speed = float64(skills.CapMath(c.Stats.BaseMoveSpd.Run, skill.Effect.PSpeed.Val, effect.PSpeed.Cap))
	}
}

// StatsRefresh Обновление счетчика всех статов персонажа
// По задумке творца, это необходимо использовать всякий раз при новом баффе, при надевании шмотки/оружия/бижи если имеются эффекты у него.
func (c *Character) StatsRefresh() {
	c.resetBonusStatCals()
	c.getRefreshStats()
	c.bonusStatCalsBuff()
}

// Подсчет статов от бафа
func (c *Character) bonusStatCalsBuff() {
	//сброс всех значений от баффа, для перерасчета
	c.resetSkillBuffBonus()
	for _, skillbuff := range c.Buff() {
		geteffect, _ := GetSkillDataInfo(skillbuff.Id, skillbuff.Level)
		c.SkillsBuffBonus = c.setStatPlayer(geteffect.Effect, c.SkillsBuffBonus)
	}
}

func (c *Character) setStatPlayer(effect Effect, sb SkillBonus) SkillBonus {

	if effect.PMaxHp != nil {
		if effect.PMaxHp.Cap == "per" {
			sb.MaxHP += c.HP() * effect.PMaxHp.Val / 100
		}
	}
	if effect.PMaxMp != nil {
		if effect.PMaxMp.Cap == "per" {
			sb.MaxMP += c.MaxMP() * effect.PMaxMp.Val / 100
		}
	}
	if effect.PSpeed != nil {
		if effect.PSpeed.Cap == "diff" {
			sb.Speed += effect.PSpeed.Val
		}
	}

	if effect.PPhysicalDefence != nil {
		if effect.PPhysicalDefence.Cap == "per" {
			sb.PDef += float64(c.PDef() * int(effect.PPhysicalDefence.Val) / 100)
		}
	}

	if effect.PMagicalDefence != nil {
		if effect.PMagicalDefence.Cap == "per" {
			sb.MDef += float64(c.MDef() * int(effect.PMagicalDefence.Val) / 100)
		}
	}

	if effect.PPhysicalAttack != nil {
		if effect.PPhysicalAttack.Cap == "per" {
			sb.PAtk += float64(c.PAtk() * int(effect.PPhysicalAttack.Val) / 100)
		}
	}

	if effect.PMagicalAttack != nil {
		if effect.PMagicalAttack.Cap == "per" {
			sb.MAtk += float64(c.MAtk() * int(effect.PMagicalAttack.Val) / 100)
		}
	}

	if effect.PAttackSpeed != nil {
		if effect.PAttackSpeed.Cap == "per" {
			sb.AtkSpd += float64(c.AttackSpeed() * int(effect.PAttackSpeed.Val) / 100)
		}
	}

	if effect.PCriticalRate != nil {
		if effect.PCriticalRate.Cap == "per" {
			sb.CriticalRate += float64(c.CriticalRate() * int(effect.PCriticalRate.Val) / 100)
		}
	}

	return sb
}

// GetRefreshStats Обновление статов персонажа
// Берется статы из оружия, брони.
// TODO: Скиллы и бижа не учитывается
func (c *Character) getRefreshStats() {
	c.resetBonusStatCals()
	c.BonusStats = nil
	for _, v := range &c.Paperdoll {
		if v.ObjId != 0 {
			c.AddBonusStat(v.BonusStats)
			c.BonusStatCals(v)
		}
	}

}

// BonusStatCals Сложение всех статов от предметов надетых на персонаже
func (c *Character) BonusStatCals(item MyItem) {
	for _, v := range item.BonusStats {
		if v.Type == "physical_damage" {
			c.Stats.BasePAtk += int(v.Val)
		}
		if v.Type == "physical_defense" {
			switch uint8(item.LocData) {
			case PAPERDOLL_HEAD:
				c.Stats.BasePDef.Head += v.Val
			case PAPERDOLL_CHEST:
				c.Stats.BasePDef.Chest += v.Val
			case PAPERDOLL_LEGS:
				c.Stats.BasePDef.Legs += v.Val
			case PAPERDOLL_GLOVES:
				c.Stats.BasePDef.Gloves += v.Val
			case PAPERDOLL_FEET:
				c.Stats.BasePDef.Feet += v.Val
			case PAPERDOLL_UNDER:
				c.Stats.BasePDef.Underwear += v.Val
			case PAPERDOLL_CLOAK:
				c.Stats.BasePDef.Cloak += v.Val
			}
		}
		if v.Type == "magical_damage" {
			c.Stats.BaseMAtk += int(v.Val)
		}
		if v.Type == "magical_defense" {
			logger.Info.Println("magical_defense", v.Val, item.LocData)
			switch uint8(item.LocData) {
			case PAPERDOLL_RFINGER:
				c.Stats.BaseMDef.Rfinger += v.Val
			case PAPERDOLL_LFINGER:
				c.Stats.BaseMDef.Lfinger += v.Val
			case PAPERDOLL_LBRACELET:
				c.Stats.BaseMDef.Lear += v.Val
			case PAPERDOLL_RBRACELET:
				c.Stats.BaseMDef.Rear += v.Val
			case PAPERDOLL_NECK:
				c.Stats.BaseMDef.Neck += v.Val
			}
			logger.Info.Println(c.MDef())
		}
		if v.Type == "critical" {
			c.Stats.BaseCritRate += int(v.Val)
		}
		if v.Type == "attack_speed" {
			c.Stats.BasePAtkSpd += int(v.Val)
		}
		if v.Type == "mp_bonus" {
			c.SetAddMP(v.Val)
		}
	}

}

func (c *Character) resetBonusStatCals() {
	c.Stats = AllStats[int(c.ClassID())].StaticData
}

func (c *Character) GetInventoryLimit() int16 {
	if c.Race == race.DWARF {
		return 100
	}
	return 80
}

// PDef Pdef всего
func (c *Character) PDef() int {
	return data.CalcFloat64(c.SkillsBuffBonus.PDef, c.Stats.BasePDef.Gloves, c.Stats.BasePDef.Legs, c.Stats.BasePDef.Underwear, c.Stats.BasePDef.Feet, c.Stats.BasePDef.Chest, c.Stats.BasePDef.Head)
}

func (c *Character) MDef() int {
	return data.CalcFloat64(c.SkillsBuffBonus.MDef, c.Stats.BaseMDef.Rfinger, c.Stats.BaseMDef.Lfinger, c.Stats.BaseMDef.Rear, c.Stats.BaseMDef.Lear, c.Stats.BaseMDef.Neck)
}

func (c *Character) MAtk() int {
	return data.CalcFloat64(c.SkillsBuffBonus.MAtk, float64(c.Stats.BaseMAtk))
}

func (c *Character) PAtk() int {
	return data.CalcFloat64(c.SkillsBuffBonus.PAtk, float64(c.Stats.BasePAtk))
}

func (c *Character) AttackSpeed() int {
	return data.CalcFloat64(c.SkillsBuffBonus.AtkSpd, float64(c.Stats.BasePAtkSpd))
}

func (c *Character) CriticalRate() int {
	return data.CalcFloat64(c.SkillsBuffBonus.CriticalRate, float64(c.Stats.BaseCritRate))
}

//TODO: Необходимо узнать информацию о суммировании маг.спида.
func (c *Character) MAttackSpeed() int {
	return data.CalcFloat64(c.SkillsBuffBonus.MAtkSpd)
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
				if v.ObjectID() == c.ObjectID() {
					continue
				}
				deleteObjectPkgTo = append(deleteObjectPkgTo, v.ObjectID())
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
				if v.ObjectID() == c.ObjectID() {
					continue
				}
				charInfoPkgTo = append(charInfoPkgTo, v.ObjectID())
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

//Возвращает общее кол-во ХП (дающее и с скиллы, с предметы и баффы)
func (c *Character) MaxHP() float64 {
	return c.HP() + c.SkillsItemBonus.MaxHP + c.SkillsBuffBonus.MaxHP
}

func (c *Character) MaxMP() float64 {
	return c.MP() + c.SkillsItemBonus.MaxMP + c.SkillsBuffBonus.MaxMP
}

func (c *Character) MaxCP() float64 {
	return c.CP() + c.SkillsItemBonus.MaxCP + c.SkillsBuffBonus.MaxCP
}

// GetMaxRunSpeed Возвращает скорость бега персонажа учитывая все баффы и .тд.
func (c *Character) MaxRunSpeed() float64 {
	return c.Stats.BaseMoveSpd.Run + c.SkillsItemBonus.Speed + c.SkillsBuffBonus.Speed
}

// Buff Список баффа
func (c *Character) Buff() []*BuffUser {
	return c.buff
}

// SetBuff Установить новый список баффа
func (c *Character) SetBuff(buff []*BuffUser) {
	c.buff = buff
}

// AddBuff Добавить бафф и отсекаем лишний
func (c *Character) AddBuff(id, level, second int) {
	c.buff = append(c.buff, &BuffUser{
		Id:     id,
		Level:  level,
		Second: second,
	})
	c.buffSifting()
}

// ClearBuff Очищение списка баффа
func (c *Character) ClearBuff() {
	c.buff = nil
}

// RemoveBuffSkill Удаление баффа
func (c *Character) RemoveBuffSkill(id int) {
	for index, buff := range c.Buff() {
		if buff.Id == id {
			c.buff = append(c.buff[:index], c.buff[index+1:]...)
		}
	}
}

// buffSifting Функция просеивания и сравнения баффов
// Удаление дубликатов, и оставляет бафф (если одинаковый лвл) который больше по времени будет действовать
// Если время баффов одинаковое, тогда применяется бафф больше по уровню
func (c *Character) buffSifting() {
	if len(c.Buff()) == 0 {
		return
	}
	var unique []*BuffUser
	buffGet := func(unique []*BuffUser, id int) (*BuffUser, int, bool) {
		for index, buff := range unique {
			if buff.Id == id {
				return buff, index, true
			}
		}
		return nil, 0, false
	}
	for _, buff := range c.Buff() {
		duplicateBuff, index, ok := buffGet(unique, buff.Id)
		if ok {
			if duplicateBuff.Second < buff.Second || duplicateBuff.Second == buff.Second && duplicateBuff.Level < buff.Level {
				unique = append(unique[:index], unique[index+1:]...)
				unique = append(unique, buff)
			}
		} else {
			unique = append(unique, buff)
			continue
		}
	}
	c.SetBuff(unique)
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

func (c *Character) GetSkillInfo(skill_id int) Skill {
	for _, skill := range c.Skills {
		if skill.SkillId == skill_id {
			return skill
		}
	}
	return Skill{}
}
