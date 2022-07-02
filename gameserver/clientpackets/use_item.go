package clientpackets

import (
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/models/items"
	"l2gogameserver/gameserver/models/race"
	"l2gogameserver/gameserver/models/sysmsg"
	"l2gogameserver/gameserver/serverpackets"
	"l2gogameserver/packets"
)

const formalWearId = 6408
const fortFlagId = 9819

func UseItem(clientI interfaces.ReciverAndSender, data []byte) {
	client, ok := clientI.(*models.Client)
	if !ok {
		return
	}
	var packet = packets.NewReader(data)

	objId := packet.ReadInt32() // targetObjId
	ctrlPressed := packet.ReadInt32() != 0
	_ = ctrlPressed

	var selectedItem *models.MyItem

	find := false
	for i := range client.CurrentChar.Inventory.Items {
		item := client.CurrentChar.Inventory.Items[i]
		if item.ObjId == objId {
			selectedItem = item
			find = true
			break
		}
	}

	// если предмет не найден в инвентаре, то выходим
	if !find {
		return
	}

	//Сюда попадают предметы при двойном клике, которые не являются надеваемыми, к примеру адена, свитки, ресурсы
	if selectedItem.ItemType == items.Other || selectedItem.ItemType == items.Money || selectedItem.ItemType == items.Quest {
		logger.Warning.Printf("Предмет %s id %d", selectedItem.Name, selectedItem.Id)
		return
	}

	if selectedItem.IsEquipable() {
		// нельзя надевать Formal Wear с проклятым оружием
		if client.CurrentChar.IsCursedWeaponEquipped() && objId == formalWearId {
			return
		}

		// todo тут еще 2 проверки

		switch selectedItem.SlotBitType {
		case items.SlotLrHand, items.SlotLHand, items.SlotRHand:

			// если в руке Combat flag
			if client.CurrentChar.IsActiveWeapon() && models.GetActiveWeapon(client.CurrentChar.Inventory.Items, client.CurrentChar.Paperdoll).Item.Id == fortFlagId {
				pkg := serverpackets.SystemMessage(sysmsg.CannotEquipItemDueToBadCondition, client)
				client.EncryptAndSend(pkg)
				return
			}
			//todo тут 2 проврки на  isMounted  и isDisarmed

			// нельзя менять оружие/щит если в руках проклятое оружие
			if client.CurrentChar.IsCursedWeaponEquipped() {
				return
			}

			//  запрет носить НЕ камаелям эксклюзивное оружие  камаелей
			if selectedItem.IsEquipped() == 0 && selectedItem.IsWeapon() { // todo еще проверка && !activeChar.canOverrideCond(ITEM_CONDITIONS))

				switch client.CurrentChar.Race {
				case race.KAMAEL:
					if selectedItem.IsWeaponTypeNone() {
						pkg := serverpackets.SystemMessage(sysmsg.CannotEquipItemDueToBadCondition, client)
						client.EncryptAndSend(pkg)

						return
					}
				case race.HUMAN, race.DWARF, race.ELF, race.DARK_ELF, race.ORC:
					if selectedItem.IsOnlyKamaelWeapon() {
						pkg := serverpackets.SystemMessage(sysmsg.CannotEquipItemDueToBadCondition, client)
						client.EncryptAndSend(pkg)
						return
					}
				}
			}
		// камаель не может носить тяжелую или маг броню
		// они могут носить только лайт, может проверять на !LIGHT ?
		case items.SlotChest, items.SlotBack, items.SlotGloves, items.SlotFeet, items.SlotHead, items.SlotFullArmor, items.SlotLegs:
			if client.CurrentChar.Race == race.KAMAEL && (selectedItem.IsHeavyArmor() || selectedItem.IsMagicArmor()) {
				pkg := serverpackets.SystemMessage(sysmsg.CannotEquipItemDueToBadCondition, client)
				client.EncryptAndSend(pkg)
				return
			}
		case items.SlotDeco:
			//todo проверка !item.isEquipped() && (activeChar.getInventory().getTalismanSlots() == 0

		}

	}

	//Если выбранный предмет надет и пользователь кликает по нем - снимаем его и кладем в инвентарь
	if selectedItem.IsEquipped() == 1 {
		//Нужно сделать проверку, занят ли слот предмета, если занят, то снять и надеть шмотку новую
		client.CurrentChar.ItemTakeOff(selectedItem, client.CurrentChar.GetFirstEmptySlot())
		clientI.EncryptAndSend(serverpackets.InventoryUpdate(selectedItem, models.UpdateTypeModify))
	} else {
		//Когда предмет не надет на персонаже, и персонаж кликнул надеть предмет
		//Сначала проверим, свободный ли слот, туда куда наденется шмот
		slotBitType := client.CurrentChar.SlotItemInfo(selectedItem)
		if slotBitType == 255 {
			logger.Error.Panicln("Ошибка, не найден слот предмета")
		}
		//Поиск занятого слота
		busySlotItem, ok := client.CurrentChar.EmptyPaperdollSlot(slotBitType)
		logger.Info.Println(busySlotItem.Name, ok)
		if ok { //Опусташаем слот перед надеванием
			client.CurrentChar.ItemTakeOff(busySlotItem, selectedItem.LocData)
			clientI.EncryptAndSend(serverpackets.InventoryUpdate(busySlotItem, models.UpdateTypeModify))
		}
		//Надеваем шмот в пустой слот
		client.CurrentChar.ItemPutOn(selectedItem, uint8(selectedItem.LocData))

		clientI.EncryptAndSend(serverpackets.InventoryUpdate(selectedItem, models.UpdateTypeModify))
	}

	//client.CurrentChar.Paperdoll[selectedItem.LocData] = *selectedItem

	//client.CurrentChar.ShowItemsEquipped()

	//Проверка скиллов предмета
	client.CurrentChar.SkillItemListRefresh()
	clientI.EncryptAndSend(serverpackets.SkillList(client))

	clientI.EncryptAndSend(serverpackets.UserInfo(client))

}
