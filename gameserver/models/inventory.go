package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/idfactory"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models/items"
	"l2gogameserver/gameserver/models/items/armorType"
	"l2gogameserver/gameserver/models/items/attribute"
	"l2gogameserver/gameserver/models/items/consumeType"
	"l2gogameserver/gameserver/models/items/etcItemType"
	"l2gogameserver/gameserver/models/items/weaponType"
	"strconv"
	"strings"
)

const (
	PAPERDOLL_UNDER      uint8 = 0
	PAPERDOLL_HEAD       uint8 = 1
	PAPERDOLL_HAIR       uint8 = 2
	PAPERDOLL_HAIR2      uint8 = 3
	PAPERDOLL_NECK       uint8 = 4
	PAPERDOLL_RHAND      uint8 = 5
	PAPERDOLL_CHEST      uint8 = 6
	PAPERDOLL_LHAND      uint8 = 7
	PAPERDOLL_REAR       uint8 = 8
	PAPERDOLL_LEAR       uint8 = 9
	PAPERDOLL_GLOVES     uint8 = 10
	PAPERDOLL_LEGS       uint8 = 11
	PAPERDOLL_FEET       uint8 = 12
	PAPERDOLL_RFINGER    uint8 = 13
	PAPERDOLL_LFINGER    uint8 = 14
	PAPERDOLL_LBRACELET  uint8 = 15
	PAPERDOLL_RBRACELET  uint8 = 16
	PAPERDOLL_DECO1      uint8 = 17
	PAPERDOLL_DECO2      uint8 = 18
	PAPERDOLL_DECO3      uint8 = 19
	PAPERDOLL_DECO4      uint8 = 20
	PAPERDOLL_DECO5      uint8 = 21
	PAPERDOLL_DECO6      uint8 = 22
	PAPERDOLL_CLOAK      uint8 = 23
	PAPERDOLL_BELT       uint8 = 24
	PAPERDOLL_TOTALSLOTS uint8 = 25

	PaperdollLoc string = "PAPERDOLL"
	InventoryLoc string = "INVENTORY"

	UpdateTypeUnchanged int16 = 0
	UpdateTypeAdd       int16 = 1
	UpdateTypeModify    int16 = 2
	UpdateTypeRemove    int16 = 3
)

type MyItem struct {
	items.Item
	ObjId               int32
	Enchant             int
	LocData             int32 // ID слота, который занят в инвентаре
	Count               int64
	Loc                 string
	Time                int
	AttackAttributeType attribute.Attribute
	AttackAttributeVal  int
	Mana                int
	AttributeDefend     [6]int16
}

type Inventory struct {
	Items   []*MyItem
	IsEquip MyItem
}

// NewItemCreate Создает новый экземпляр предмета
func NewItemCreate(character *Character, itemid int, count int64) (*MyItem, bool) {
	itemData, ok := items.GetItemInfo(itemid)
	if !ok {
		logger.Error.Println("Предмет в инвентаре не найден")
		return &MyItem{}, false
	}
	newItem := MyItem{
		Item:    itemData,
		ObjId:   idfactory.GetNext(),
		LocData: character.GetFirstEmptySlot(),
		Count:   count,
		Loc:     InventoryLoc,
	}
	return &newItem, true
}

// IsEquipWeapon Возращает информацию о экиперованном оружии
func (i Inventory) IsEquipWeapon() (*MyItem, bool) {
	for _, item := range i.Items {
		if item.Loc == PaperdollLoc && item.ItemType == items.Weapon {
			return item, true
		}
	}
	return &MyItem{}, false
}

// AddItem Добавление предмета в инвентарь, возвращает модификатор действий
func (i *Inventory) AddItem(item *MyItem) (*MyItem, int16) {
	_, indexObject, isItemInventory := i.ExistItemID(item.Id)

	if isItemInventory == false { //Если предмет не найден в инвентаре
		i.Items = append(i.Items, item)
		return item, UpdateTypeAdd
	}

	//Если предмет стакуем
	if item.ConsumeType == consumeType.Stackable || item.ConsumeType == consumeType.Asset {
		i.Items[indexObject].Count += item.Count
		return i.Items[indexObject], UpdateTypeModify
	} else { //Если не стакуется, тогда отдельно добавляем предмет
		i.Items = append(i.Items, item)
		return item, UpdateTypeAdd
	}

}

//При загрузке персонажа, получаем все данные "одетые" на персонаже
func (c Character) LoadingVisibleInventory() [26]MyItem {

	dbConn, err := db.GetConn()
	if err != nil {
		//TODO: почему-то, нужно будет проанализировать
		//БД запущена, всё нормально. Однако зажал Enter после авторизации, чтоб сразу перейти к загрузке.
		/*
			ERROR: 00:02:48 inventory.go:129: failed to connect to `host=localhost user=postgres database=postgres`: dial error (dial tcp [::1]:5432: connectex:
			No connection could be made because the target machine actively refused it.)
			panic: failed to connect to `host=localhost user=postgres database=postgres`: dial error (dial tcp [::1]:5432: connectex: No connection could be made
			 because the target machine actively refused it.)


			goroutine 36 [running]:
			log.(*Logger).Panicln(0xffffd7200002c224?, {0xc0034070a0?, 0x0?, 0x0?})
			        C:/Program Files/Go/src/log/log.go:262 +0x69
			l2gogameserver/gameserver/models.Character.LoadingVisibleInventory({{0xc004fb065c, 0x4}, 0x9613, {0xc004fb0660, 0xa}, 0x4b, 0x409b200000000000, 0x3ff
			0000000000000, 0x4089400000000000, 0x3ff0000000000000, ...})
			        C:/go/l2gogameserver/gameserver/models/inventory.go:129 +0xa5
			l2gogameserver/gameserver/serverpackets.CharSelectionInfo({0xc75a08?, 0xc00010e1b0?})

		*/
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	rows, err := dbConn.Query(context.Background(), "SELECT object_id, item, loc_data, enchant_level FROM items WHERE owner_id= $1 AND loc= $2", c.ObjectID(), PaperdollLoc)
	if err != nil {
		logger.Error.Panicln(err)
	}

	var mts [26]MyItem

	for rows.Next() {
		var objId int
		var itemId int
		var enchantLevel int
		var locData int
		err = rows.Scan(&objId, &itemId, &locData, &enchantLevel)
		if err != nil {
			logger.Info.Println(err)
		}
		item, ok := items.GetItemFromStorage(itemId)
		if !ok {
			logger.Error.Panicln("Предмет не найден")
		}
		mt := MyItem{
			Item:    item,
			ObjId:   int32(objId),
			Enchant: enchantLevel,
			Count:   1,
			Loc:     PaperdollLoc,
		}
		mts[int32(locData)] = mt
	}
	return mts
}

// RestoreVisibleInventory Получение всех ID предметов, которые экиперованы на персонаже
func (c *Character) RestoreVisibleInventory() [26]MyItem {

	var mts [26]MyItem

	for _, item := range c.Inventory.Items {
		if item.Loc == PaperdollLoc {
			mt := MyItem{
				Item:    item.Item,
				ObjId:   item.ObjId,
				Enchant: item.Enchant,
				Count:   item.Count,
				Loc:     PaperdollLoc,
			}
			mts[item.LocData] = mt
		}
	}

	return mts

}

// IsEquipable Можно ли надеть предмет
func (i *MyItem) IsEquipable() bool {
	return !((i.SlotBitType == items.SlotNone) || (i.EtcItemType == etcItemType.ARROW) || (i.EtcItemType == etcItemType.BOLT) || (i.EtcItemType == etcItemType.LURE))
}
func (i *MyItem) IsHeavyArmor() bool {
	return i.ArmorType == armorType.HEAVY
}
func (i *MyItem) IsMagicArmor() bool {
	return i.ArmorType == armorType.MAGIC
}
func (i *MyItem) IsArmor() bool {
	return i.ItemType == items.ShieldOrArmor
}
func (i *MyItem) IsOnlyKamaelWeapon() bool {
	return i.WeaponType == weaponType.RAPIER || i.WeaponType == weaponType.CROSSBOW || i.WeaponType == weaponType.ANCIENTSWORD
}
func (i *MyItem) IsWeapon() bool {
	return i.ItemType == items.Weapon
}
func (i *MyItem) IsWeaponTypeNone() bool {
	return i.WeaponType == weaponType.NONE
}
func GetMyItems(charId int32) Inventory {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	sqlString := "SELECT items.object_id, item, loc_data, enchant_level, count, loc, time, mana_left FROM items WHERE owner_id = $1"
	rows, err := dbConn.Query(context.Background(), sqlString, charId)
	if err != nil {
		logger.Error.Panicln(err)
	}

	var inventory Inventory

	for rows.Next() {
		var itm MyItem
		var id int

		err := rows.Scan(&itm.ObjId, &id, &itm.LocData, &itm.Enchant, &itm.Count, &itm.Loc, &itm.Time, &itm.Mana)
		if err != nil {
			logger.Error.Panicln(err)
		}

		it, ok := items.GetItemFromStorage(id)
		if ok {
			itm.Item = it

			if itm.IsWeapon() {

				itm.AttackAttributeType, itm.AttackAttributeVal = getAttributeForWeapon(itm.ObjId)
			} else if itm.IsArmor() {
				itm.AttributeDefend = getAttributeForArmor(itm.ObjId)
			}
			inventory.Items = append(inventory.Items, &itm)
		}
	}
	//RefreshLocData(inventory.Items)
	return inventory
}

// RefreshLocData Сброс LocData всех предметов
//func RefreshLocData(Items []*MyItem) []*MyItem {
//for _, item := range Items {
//	item.LocData = -1
//item.LocData = character.GetFirstEmptySlot()
//}
//return Items
//}

func getAttributeForWeapon(objId int32) (attribute.Attribute, int) {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()
	el := attribute.Attribute(-2) // None

	var elementType, elementValue int
	err = dbConn.QueryRow(context.Background(), "SELECT element_type,element_value FROM item_elementals WHERE item_id = $1", objId).
		Scan(&elementType, &elementValue)

	if err == pgx.ErrNoRows {
		return el, 0
	} else if err != nil {
		logger.Error.Panicln(err)
	}

	el = attribute.Attribute(elementType)

	return el, elementValue
}

func getAttributeForArmor(objId int32) [6]int16 {
	var att [6]int16
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	rows, err := dbConn.Query(context.Background(), "SELECT element_type,element_value FROM item_elementals WHERE item_id = $1", objId)

	if err == pgx.ErrNoRows {
		return att
	} else if err != nil {
		logger.Error.Panicln(err)
	}

	for rows.Next() {
		var atType, atVal int
		err = rows.Scan(&atType, &atVal)
		if err != nil {
			logger.Error.Panicln(err)
		}
		att[atType] = int16(atVal)
	}

	return att
}

func (i *MyItem) IsEquipped() int16 {
	if i.Loc == InventoryLoc {
		return 0
	}
	return 1
}

func SaveInventoryInDB(inventory []*MyItem) {
	if len(inventory) == 0 {
		return
	}

	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	var sb strings.Builder
	sb.WriteString("UPDATE items SET loc_data = mylocdata, loc = myloc FROM ( VALUES ")

	for i := range inventory {
		v := inventory[i]

		sb.WriteString("(" + strconv.Itoa(int(v.LocData)) + ",'" + v.Loc + "'," + strconv.Itoa(int(v.ObjId)) + ")")

		if len(inventory)-1 != i {
			sb.WriteString(",")
		}
	}
	sb.WriteString(") as myval (mylocdata,myloc,myobjid) WHERE items.object_id = myval.myobjid")
	_, err = dbConn.Exec(context.Background(), sb.String())
	if err != nil {
		logger.Info.Println(err.Error())
	}
}

func GetActiveWeapon(inventory []*MyItem, paperdoll [26]MyItem) *MyItem {
	q := paperdoll[PAPERDOLL_RHAND]
	for i := range inventory {
		v := inventory[i]
		if v.ObjId == q.ObjId {
			return v
		}
	}
	return nil
}

// UseEquippableItem использовать предмет который можно надеть на персонажа
func UseEquippableItem(selectedItem *MyItem, character *Character) (*MyItem, bool) {
	//todo надо как то обновлять paperdoll, или возвращать массив или же  вынести это в другой пакет
	if selectedItem.IsEquipped() == 1 {
		logger.Info.Println("Предмет снимаем")
		return unEquipAndRecord(selectedItem, character)
	} else {
		logger.Info.Println("Надеваем этот:", selectedItem.Name)
		return equipItemAndRecord(selectedItem, character)
	}
}

// unEquipAndRecord cнять предмет
func unEquipAndRecord(selectedItem *MyItem, character *Character) (*MyItem, bool) {
	switch selectedItem.SlotBitType {
	case items.SlotLEar:
		return setPaperdollItemToInventary(PAPERDOLL_LEAR, selectedItem, character)
	case items.SlotREar:
		return setPaperdollItemToInventary(PAPERDOLL_REAR, selectedItem, character)
	case items.SlotNeck:
		return setPaperdollItemToInventary(PAPERDOLL_NECK, selectedItem, character)
	case items.SlotRFinger:
		return setPaperdollItemToInventary(PAPERDOLL_RFINGER, selectedItem, character)
	case items.SlotLFinger:
		return setPaperdollItemToInventary(PAPERDOLL_LFINGER, selectedItem, character)
	case items.SlotHair:
		return setPaperdollItemToInventary(PAPERDOLL_HAIR, selectedItem, character)
	case items.SlotHair2:
		return setPaperdollItemToInventary(PAPERDOLL_HAIR2, selectedItem, character)
	case items.SlotHairall: //todo Разобраться что тут на l2j
		return setPaperdollItemToInventary(PAPERDOLL_HAIR, selectedItem, character)
	case items.SlotHead:
		return setPaperdollItemToInventary(PAPERDOLL_HEAD, selectedItem, character)
	case items.SlotRHand, items.SlotLrHand:
		return setPaperdollItemToInventary(PAPERDOLL_RHAND, selectedItem, character)
	case items.SlotLHand:
		return setPaperdollItemToInventary(PAPERDOLL_LHAND, selectedItem, character)
	case items.SlotGloves:
		return setPaperdollItemToInventary(PAPERDOLL_GLOVES, selectedItem, character)
	case items.SlotChest, items.SlotAlldress, items.SlotFullArmor:
		return setPaperdollItemToInventary(PAPERDOLL_CHEST, selectedItem, character)
	case items.SlotLegs:
		return setPaperdollItemToInventary(PAPERDOLL_LEGS, selectedItem, character)
	case items.SlotBack:
		return setPaperdollItemToInventary(PAPERDOLL_CLOAK, selectedItem, character)
	case items.SlotFeet:
		return setPaperdollItemToInventary(PAPERDOLL_FEET, selectedItem, character)
	case items.SlotUnderwear:
		return setPaperdollItemToInventary(PAPERDOLL_UNDER, selectedItem, character)
	case items.SlotLBracelet:
		return setPaperdollItemToInventary(PAPERDOLL_LBRACELET, selectedItem, character)
	case items.SlotRBracelet:
		return setPaperdollItemToInventary(PAPERDOLL_RBRACELET, selectedItem, character)
	case items.SlotDeco:
		return setPaperdollItemToInventary(PAPERDOLL_DECO1, selectedItem, character)
	case items.SlotBelt:
		return setPaperdollItemToInventary(PAPERDOLL_BELT, selectedItem, character)
	}
	return nil, false
}

// SlotItemInfo Возвращает ID позиции надетого предмета
func (c Character) SlotItemInfo(selectedItem *MyItem) uint8 {
	paperdoll := c.Paperdoll
	switch selectedItem.SlotBitType {
	case items.SlotLrHand:
		//setPaperdollItem(PAPERDOLL_LHAND, nil, character)
		return PAPERDOLL_RHAND
	case items.SlotLEar, items.SlotREar, items.SlotLrEar:
		if paperdoll[PAPERDOLL_LEAR].ObjId == 0 {
			return PAPERDOLL_LEAR
		} else if paperdoll[PAPERDOLL_REAR].ObjId == 0 {
			return PAPERDOLL_REAR
		} else {
			return PAPERDOLL_LEAR
		}
	case items.SlotNeck:
		return PAPERDOLL_NECK
	case items.SlotRFinger, items.SlotLFinger, items.SlotLrFinger:
		if paperdoll[PAPERDOLL_LFINGER].ObjId == 0 {
			return PAPERDOLL_LFINGER
		} else if paperdoll[PAPERDOLL_RFINGER].ObjId == 0 {
			return PAPERDOLL_RFINGER
		} else {
			return PAPERDOLL_LFINGER
		}
	case items.SlotHair:
		hair := paperdoll[PAPERDOLL_HAIR]
		if hair.ObjId != 0 && hair.SlotBitType == items.SlotHairall {
			return PAPERDOLL_HAIR2
		} else {
			return PAPERDOLL_HAIR
		}
		return PAPERDOLL_HAIR
	case items.SlotHair2:
		hair2 := paperdoll[PAPERDOLL_HAIR]
		if hair2.ObjId != 0 && hair2.SlotBitType == items.SlotHairall {
			return PAPERDOLL_HAIR
		} else {
			return PAPERDOLL_HAIR2
		}
		return PAPERDOLL_HAIR2
	case items.SlotHairall:
		return PAPERDOLL_HAIR2
		return PAPERDOLL_HAIR
	case items.SlotHead:
		return PAPERDOLL_HEAD
	case items.SlotRHand:
		//todo снять стрелы
		return PAPERDOLL_RHAND
	case items.SlotLHand:
		rh := paperdoll[PAPERDOLL_RHAND]
		if (rh.ObjId != 0) && (rh.SlotBitType == items.SlotLrHand) && !(((rh.WeaponType == weaponType.BOW) && (selectedItem.EtcItemType == etcItemType.ARROW)) || ((rh.WeaponType == weaponType.CROSSBOW) && (selectedItem.EtcItemType == etcItemType.BOLT)) || ((rh.WeaponType == weaponType.FISHINGROD) && (selectedItem.EtcItemType == etcItemType.LURE))) {
			return PAPERDOLL_RHAND
		}
		return PAPERDOLL_LHAND
	case items.SlotGloves:
		return PAPERDOLL_GLOVES
	case items.SlotChest:
		return PAPERDOLL_CHEST
	case items.SlotLegs:
		chest := paperdoll[PAPERDOLL_CHEST]
		if chest.ObjId != 0 && chest.SlotBitType == items.SlotFullArmor {
			return PAPERDOLL_CHEST
		}
		return PAPERDOLL_LEGS
	case items.SlotBack:
		return PAPERDOLL_CLOAK
	case items.SlotFeet:
		return PAPERDOLL_FEET
	case items.SlotUnderwear:
		return PAPERDOLL_UNDER
	case items.SlotLBracelet:
		return PAPERDOLL_LBRACELET
	case items.SlotRBracelet:
		return PAPERDOLL_RBRACELET
	case items.SlotDeco:
		return PAPERDOLL_DECO1
	case items.SlotBelt:
		return PAPERDOLL_BELT
	case items.SlotFullArmor:
		return PAPERDOLL_LEGS
		return PAPERDOLL_CHEST
	case items.SlotAlldress:
		logger.Error.Panicln("Слот все адреса")
		//return PAPERDOLL_LEGS
		//return PAPERDOLL_LHAND
		//return PAPERDOLL_RHAND
		//return PAPERDOLL_HEAD
		//return PAPERDOLL_FEET
		//return PAPERDOLL_GLOVES
		//return PAPERDOLL_CHEST
	default:
		logger.Error.Println("Не определен Slot для itemId: "+strconv.Itoa(selectedItem.Id), "вероятно это не относится к шмоту")
	}
	return 255
}

// equipItemAndRecord одеть предмет
func equipItemAndRecord(selectedItem *MyItem, character *Character) (*MyItem, bool) {
	//todo проверка на приват Store, надо будет передавать character?
	// еще проверка на ITEM_CONDITIONS
	formal := character.Paperdoll[PAPERDOLL_CHEST]
	// Проверка надето ли офф. одежда и предмет не является букетом(id=21163)
	if (selectedItem.Id != 21163) && (formal.ObjId != 0) && (formal.SlotBitType == items.SlotAlldress) {
		// только chest можно
		switch selectedItem.SlotBitType {
		case items.SlotLrHand, items.SlotLHand, items.SlotRHand, items.SlotLegs, items.SlotFeet, items.SlotGloves, items.SlotHead:
			return nil, false
		}
	}

	paperdoll := character.Paperdoll
	switch selectedItem.SlotBitType {
	case items.SlotLrHand:
		//setPaperdollItem(PAPERDOLL_LHAND, nil, character)
		return setPaperdollItem(PAPERDOLL_RHAND, selectedItem, character)
	case items.SlotLEar, items.SlotREar, items.SlotLrEar:
		if paperdoll[PAPERDOLL_LEAR].ObjId == 0 {
			return setPaperdollItem(PAPERDOLL_LEAR, selectedItem, character)
		} else if paperdoll[PAPERDOLL_REAR].ObjId == 0 {
			return setPaperdollItem(PAPERDOLL_REAR, selectedItem, character)
		} else {
			return setPaperdollItem(PAPERDOLL_LEAR, selectedItem, character)
		}

	case items.SlotNeck:
		return setPaperdollItem(PAPERDOLL_NECK, selectedItem, character)
	case items.SlotRFinger, items.SlotLFinger, items.SlotLrFinger:
		if paperdoll[PAPERDOLL_LFINGER].ObjId == 0 {
			return setPaperdollItem(PAPERDOLL_LFINGER, selectedItem, character)
		} else if paperdoll[PAPERDOLL_RFINGER].ObjId == 0 {
			return setPaperdollItem(PAPERDOLL_RFINGER, selectedItem, character)
		} else {
			return setPaperdollItem(PAPERDOLL_LFINGER, selectedItem, character)
		}

	case items.SlotHair:
		hair := paperdoll[PAPERDOLL_HAIR]
		if hair.ObjId != 0 && hair.SlotBitType == items.SlotHairall {
			return setPaperdollItem(PAPERDOLL_HAIR2, selectedItem, character)
		} else {
			return setPaperdollItem(PAPERDOLL_HAIR, selectedItem, character)
		}
		return setPaperdollItem(PAPERDOLL_HAIR, selectedItem, character)
	case items.SlotHair2:
		hair2 := paperdoll[PAPERDOLL_HAIR]
		if hair2.ObjId != 0 && hair2.SlotBitType == items.SlotHairall {
			return setPaperdollItem(PAPERDOLL_HAIR, selectedItem, character)
		} else {
			return setPaperdollItem(PAPERDOLL_HAIR2, selectedItem, character)
		}
		return setPaperdollItem(PAPERDOLL_HAIR2, selectedItem, character)
	case items.SlotHairall:
		return setPaperdollItem(PAPERDOLL_HAIR2, selectedItem, character)
		return setPaperdollItem(PAPERDOLL_HAIR, selectedItem, character)
	case items.SlotHead:
		return setPaperdollItem(PAPERDOLL_HEAD, selectedItem, character)
	case items.SlotRHand:
		//todo снять стрелы
		return setPaperdollItem(PAPERDOLL_RHAND, selectedItem, character)
	case items.SlotLHand:
		rh := paperdoll[PAPERDOLL_RHAND]
		if (rh.ObjId != 0) && (rh.SlotBitType == items.SlotLrHand) && !(((rh.WeaponType == weaponType.BOW) && (selectedItem.EtcItemType == etcItemType.ARROW)) || ((rh.WeaponType == weaponType.CROSSBOW) && (selectedItem.EtcItemType == etcItemType.BOLT)) || ((rh.WeaponType == weaponType.FISHINGROD) && (selectedItem.EtcItemType == etcItemType.LURE))) {
			return setPaperdollItem(PAPERDOLL_RHAND, selectedItem, character)
		}
		return setPaperdollItem(PAPERDOLL_LHAND, selectedItem, character)
	case items.SlotGloves:
		return setPaperdollItem(PAPERDOLL_GLOVES, selectedItem, character)
	case items.SlotChest:
		return setPaperdollItem(PAPERDOLL_CHEST, selectedItem, character)
	case items.SlotLegs:
		chest := paperdoll[PAPERDOLL_CHEST]
		if chest.ObjId != 0 && chest.SlotBitType == items.SlotFullArmor {
			return setPaperdollItem(PAPERDOLL_CHEST, selectedItem, character)
		}
		return setPaperdollItem(PAPERDOLL_LEGS, selectedItem, character)
	case items.SlotBack:
		return setPaperdollItem(PAPERDOLL_CLOAK, selectedItem, character)
	case items.SlotFeet:
		return setPaperdollItem(PAPERDOLL_FEET, selectedItem, character)
	case items.SlotUnderwear:
		return setPaperdollItem(PAPERDOLL_UNDER, selectedItem, character)
	case items.SlotLBracelet:
		return setPaperdollItem(PAPERDOLL_LBRACELET, selectedItem, character)
	case items.SlotRBracelet:
		return setPaperdollItem(PAPERDOLL_RBRACELET, selectedItem, character)
	case items.SlotDeco:
		return setPaperdollItem(PAPERDOLL_DECO1, selectedItem, character)
	case items.SlotBelt:
		return setPaperdollItem(PAPERDOLL_BELT, selectedItem, character)
	case items.SlotFullArmor:
		return setPaperdollItem(PAPERDOLL_LEGS, selectedItem, character)
		return setPaperdollItem(PAPERDOLL_CHEST, selectedItem, character)
	case items.SlotAlldress:
		return setPaperdollItem(PAPERDOLL_LEGS, selectedItem, character)
		return setPaperdollItem(PAPERDOLL_LHAND, selectedItem, character)
		return setPaperdollItem(PAPERDOLL_RHAND, selectedItem, character)
		return setPaperdollItem(PAPERDOLL_HEAD, selectedItem, character)
		return setPaperdollItem(PAPERDOLL_FEET, selectedItem, character)
		return setPaperdollItem(PAPERDOLL_GLOVES, selectedItem, character)
		return setPaperdollItem(PAPERDOLL_CHEST, selectedItem, character)
	default:
		logger.Error.Panicln("Не определен Slot для itemId: " + strconv.Itoa(selectedItem.Id))
	}
	return nil, false
}

//Снять предмет и поместить его в инвентарь
func setPaperdollItemToInventary(slot uint8, selectedItem *MyItem, character *Character) (*MyItem, bool) {
	character.ItemTakeOff(selectedItem, int32(slot))
	character.RemoveBonusStat(selectedItem.BonusStats)
	return selectedItem, true
}

//TODO: Всю эту функцию нужно переписать по человечески, она явно тупит и в ней нехватает обработки locdata
func setPaperdollItem(slot uint8, selectedItem *MyItem, character *Character) (*MyItem, bool) {
	// eсли selectedItem nil, то ищем предмет которых находиться в slot
	// переносим его в инвентарь, убираем бонусы этого итема у персонажа
	//for i := range character.Inventory.Items {
	//	itemInInventory := character.Inventory.Items[i]
	//	if itemInInventory.LocData == int32(slot) && itemInInventory.Loc == PaperdollLoc {
	//		itemInInventory.LocData = getFirstEmptySlot(character.Inventory.Items)
	//		itemInInventory.Loc = InventoryLoc
	//		character.Inventory.Items[i] = itemInInventory
	//		logger.Info.Println(itemInInventory.Loc, itemInInventory.LocData)
	//		character.RemoveBonusStat(itemInInventory.BonusStats)
	//		return character.Inventory.Items[i], true
	//	}
	//}

	oldItem, ok := character.GetSlotItem(slot)

	//Надеваем на слот
	character.ItemPutOn(selectedItem, slot)
	// добавить бонусы предмета персонажу
	//character.AddBonusStat(selectedItem.BonusStats)

	logger.Info.Println(ok, oldItem.Name)
	if ok { //Если слот не свободем, снимаем то что там
		character.ItemTakeOff(oldItem, character.GetFirstEmptySlot())
		logger.Info.Println("Слот НЕ свободен, нужно снять предмет который на нем и надеть новый")
		character.RemoveBonusStat(oldItem.BonusStats)
		logger.Info.Println(ok, oldItem.Name)
		//Это на случай, если пользователь снимает тот шмот, который надет
		//if oldItem.ObjId == selectedItem.ObjId {
		return oldItem, true
		//}
	}
	return selectedItem, false
	//character.Inventory.Items[keyCurrentItem] = selectedItem
}

// GetSlotItem находим предмет, который стоит на N слоте
func (c *Character) GetSlotItem(slotLocID uint8) (*MyItem, bool) {
	for _, myItem := range c.Inventory.Items {
		logger.Info.Println(myItem.LocData, myItem.Name)
		if myItem.LocData == int32(slotLocID) {
			return myItem, true
		}
	}
	return &MyItem{}, false
}

// EmptyPaperdollSlot Занят ли слот экипировки
func (c *Character) EmptyPaperdollSlot(checkSlot uint8) (*MyItem, bool) {
	for _, slot := range GetPaperdollOrder() {
		if slot == checkSlot {
			for _, myItem := range c.Inventory.Items {
				if uint8(myItem.LocData) == checkSlot {
					return myItem, true
				}
			}
		}
	}
	return &MyItem{}, false
}

// ItemPutOn Надеть вещь
func (c *Character) ItemPutOn(selectedItem *MyItem, slot uint8) {
	selectedItem.Loc = PaperdollLoc
	selectedItem.LocData = int32(slot)
}

// ItemTakeOff Снять предмет
func (c *Character) ItemTakeOff(selectedItem *MyItem, slot int32) {
	logger.Info.Println("Снимается вещь: ", selectedItem.Name)
	selectedItem.Loc = InventoryLoc
	selectedItem.LocData = slot
	c.Paperdoll[c.SlotItemInfo(selectedItem)] = MyItem{}
	logger.Info.Println(c.SlotItemInfo(selectedItem))
	//c.ShowItemsEquipped()
}

func (c *Character) GetFirstEmptySlot() int32 {
	limit := int32(80) // todo дефолтно 80 , но может быть больше
	//todo:(c)logan22, может быть больше и во время игры меняться,
	//следовательно лучше вывести в отдельную структуру с дополнительными параметры персонажа
	myItems := c.Inventory.Items
	for i := int32(0); i < limit; i++ {
		flag := false
		for j := range myItems {
			v := myItems[j]
			if v.Loc == InventoryLoc && v.LocData == i {
				flag = true
				break
			}
		}
		if !flag {
			return i
		}
	}
	logger.Error.Panicln("не нашёл куда складывать итем")
	return 0
}

func (i *MyItem) GetAttackElement() attribute.Attribute {
	el := attribute.Attribute(-2) // none
	if i.IsWeapon() {
		el = i.AttackAttributeType
	}

	if el == attribute.None {
		if i.BaseAttributeAttack.Val > 0 {
			return i.getBaseAttributeElement()
		}
	}

	return el
}

func (i *MyItem) getBaseAttributeElement() attribute.Attribute {
	return i.BaseAttributeAttack.Type
}

//TODO: Я хз че это за функция, переделай её, запрос он тут не нужен вроде.
func DeleteItem(selectedItem *MyItem, character *Character) {
	//TODO переделать, не надо создавать новый inventiry
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	if selectedItem.Loc == PaperdollLoc {
		character.Paperdoll[selectedItem.LocData] = MyItem{}
	}
	var inventory Inventory
	for _, v := range character.Inventory.Items {
		if v.ObjId != selectedItem.ObjId {
			inventory.Items = append(inventory.Items, v)
		}
	}
	character.Inventory = inventory
}

func GetPaperdollOrder() []uint8 {
	return []uint8{
		PAPERDOLL_UNDER,
		PAPERDOLL_REAR,
		PAPERDOLL_LEAR,
		PAPERDOLL_NECK,
		PAPERDOLL_RFINGER,
		PAPERDOLL_LFINGER,
		PAPERDOLL_HEAD,
		PAPERDOLL_RHAND,
		PAPERDOLL_LHAND,
		PAPERDOLL_GLOVES,
		PAPERDOLL_CHEST,
		PAPERDOLL_LEGS,
		PAPERDOLL_FEET,
		PAPERDOLL_CLOAK,
		PAPERDOLL_RHAND,
		PAPERDOLL_HAIR,
		PAPERDOLL_HAIR2,
		PAPERDOLL_RBRACELET,
		PAPERDOLL_LBRACELET,
		PAPERDOLL_DECO1,
		PAPERDOLL_DECO2,
		PAPERDOLL_DECO3,
		PAPERDOLL_DECO4,
		PAPERDOLL_DECO5,
		PAPERDOLL_DECO6,
		PAPERDOLL_BELT,
	}
}

// AddItem Добавление предмета
func AddItem(selectedItem MyItem, character *Character) Inventory {
	//Прежде чем просто добавить, необходимо проверить на существование предмета в инвентаре
	//Если он есть, тогда просто добавим к имеющимся предмету.
	//TODO: Однако, есть предметы (кроме оружия, брони, бижи), которые не стакуются, к примеру 7832
	//TODO: потом нужно определить тип предметов которые не стыкуются.
	for i := range character.Inventory.Items {
		itemInventory := character.Inventory.Items[i]
		if selectedItem.Item.Id == itemInventory.Item.Id {
			character.Inventory.Items[i].Count = itemInventory.Count + character.Inventory.Items[i].Count
			return character.Inventory
		}
	}

	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	nitem := MyItem{
		Item:                selectedItem.Item,
		ObjId:               selectedItem.ObjId,
		Enchant:             selectedItem.Enchant,
		LocData:             selectedItem.LocData,
		Count:               selectedItem.Count,
		Loc:                 "",
		Time:                selectedItem.Time,
		AttackAttributeType: selectedItem.AttackAttributeType,
		AttackAttributeVal:  selectedItem.AttackAttributeVal,
		Mana:                selectedItem.Mana,
		AttributeDefend:     [6]int16{},
	}
	character.Inventory.Items = append(character.Inventory.Items, &nitem)

	_, err = dbConn.Exec(context.Background(), `INSERT INTO "items" ("owner_id", "object_id", "item", "count", "enchant_level", "loc", "loc_data", "time_of_use", "custom_type1", "custom_type2", "mana_left", "time", "agathion_energy") VALUES ($1, $2, $3, $4, 0, 'INVENTORY', 0, 0, 0, 0, '-1', 0, 0)`, character.ObjectID(), selectedItem.ObjId, selectedItem.Item.Id, selectedItem.Count)
	if err != nil {
		logger.Error.Panicln(err)
	}

	return character.Inventory
}

//RemoveItemCharacter Удаление предмета из инвентаря персонажа
// count - сколько надо удалить
func RemoveItemCharacter(character *Character, item *MyItem, count int64) {
	logger.Info.Println("Удаление предмета из инвентаря")
	if item.Count < count || item.Count == 0 || count == 0 {
		logger.Info.Println("Неверное количество предметов для удаления")
	}
	if item.Count == count {
		DeleteItem(item, character)
		item = nil
	} else {
		newCount := item.Count - count
		item.Count = newCount
	}
}

func ExistItemObject(characterI interfaces.CharacterI, objectId int32, count int64) (*MyItem, bool) {
	character, ok := characterI.(*Character)
	if !ok {
		logger.Error.Panicln("ExistItemObject not character")
	}
	for _, item := range character.Inventory.Items {
		if item.ObjId == objectId && item.Count >= count {
			return item, true
		}
	}
	return nil, false
}

// AddInventoryItem Добавление предмета в инвентарь пользователю
// Возращаемые параметры
// 1.Ссылка на предмет
// 2.Количество
// 3.Тип обновления/удаления/добавления
// 4.True если предмет найден
func AddInventoryItem(character *Character, item MyItem, count int64) (*MyItem, int64, int16, bool) {
	for index, inv := range character.Inventory.Items {
		if inv.Item.Id == item.Id {
			if inv.IsEquipable() {
				logger.Info.Println("Нельзя передавать надетый предмет")
				return &MyItem{}, 0, UpdateTypeUnchanged, false
			}
			//Если предмет стакуемый, тогда изменим его значение
			if inv.ConsumeType == consumeType.Stackable || inv.ConsumeType == consumeType.Asset {
				inv.Count = inv.Count + count
				character.Inventory.Items[index].Count = inv.Count
				inv.Loc = "INVENTORY"
				return inv, inv.Count, UpdateTypeModify, true
			} else { //Если предмет не стакуемый, тогда добавим новое значение
				item.ObjId = idfactory.GetNext()
				item.Count = count
				item.LocData = character.GetFirstEmptySlot()
				character.Inventory.Items = append(character.Inventory.Items, &item)
				return &item, count, UpdateTypeAdd, true
			}
		}
	}
	item.ObjId = idfactory.GetNext()
	item.Count = count
	item.LocData = character.GetFirstEmptySlot()
	character.Inventory.Items = append(character.Inventory.Items, &item)
	return &item, count, UpdateTypeAdd, true
}

// RemoveItem Удаление предмета игрока
// 1.Возвращаемые параметры ссылка на предмет
// 2.Оставшейся кол-во предметов после удаления
// 3.Type удаления (Remove/Update)
// 4.Возращаемт False если предмет не был найден в инвентаре
func RemoveItem(character *Character, item *MyItem, count int64) (*MyItem, int64, int16, bool) {
	for index, itm := range character.Inventory.Items {
		if itm.Id == item.Id {
			if itm.ConsumeType == consumeType.Stackable || itm.ConsumeType == consumeType.Asset {
				itm.Count -= count
				if itm.Count <= 0 {
					character.Inventory.Items = append(character.Inventory.Items[:index], character.Inventory.Items[index+1:]...)
					return &MyItem{}, itm.Count, UpdateTypeRemove, true
				} else {
					character.Inventory.Items[index].Count = itm.Count
					return character.Inventory.Items[index], itm.Count, UpdateTypeModify, true
				}
			} else {
				character.Inventory.Items = append(character.Inventory.Items[:index], character.Inventory.Items[index+1:]...)
				return &MyItem{}, 0, UpdateTypeRemove, true
			}
		}
	}
	return &MyItem{}, 0, UpdateTypeModify, false
}

// ExistItemID Функция проверяет есть ли данный предмет в инвентаре, если есть, возращает сам объект, и его индекс
func (i *Inventory) ExistItemID(itemid int) (*MyItem, int, bool) {
	for index, item := range i.Items {
		if item.Id == itemid {
			return item, index, true
		}
	}
	return &MyItem{}, 0, false
}

// Сохранение инвентаря в базе данных
func (i Inventory) Save(charId int) {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()
	countItems := len(i.Items)
	if countItems == 0 {
		return
	}
	dbConn.Exec(context.Background(), `DELETE FROM "items" WHERE "owner_id" = $1`, charId)
	sql := `INSERT INTO "items" ("owner_id", "object_id", "item", "count", "enchant_level", "loc", "loc_data", "time_of_use", "custom_type1", "custom_type2",  "time", "agathion_energy") VALUES `
	for index, item := range i.Items {
		sql += fmt.Sprintf("(%d, %d, %d, %d, %d, '%s', %d, %d, %d, %d, '-1', %d)", charId, item.ObjId, item.Id, item.Count, item.Enchant, item.Loc, item.LocData, item.Time, item.ConsumeType, item.ConsumeType, item.Mana)
		if countItems != index+1 {
			sql += ","
		}
	}
	dbConn.Exec(context.Background(), sql)
}
