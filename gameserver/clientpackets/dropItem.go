package clientpackets

import (
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
	"l2gogameserver/packets"
)

func DropItem(clientI interfaces.ReciverAndSender, data []byte) ([]byte, models.MyItem, int16) {
	client, ok := clientI.(*models.Client)
	if !ok {
		return nil, models.MyItem{}, 0
	}

	var read = packets.NewReader(data)
	objectId := read.ReadInt32()
	count := int64(read.ReadInt32())
	_ = read.ReadInt32() // хз
	x := read.ReadInt32()
	y := read.ReadInt32()
	z := read.ReadInt32()

	item, ok := models.ExistItemObject(client.CurrentChar, objectId, count)
	if !ok {
		logger.Error.Println("Объект предмета не найден ID OBJECT:", objectId)
		return []byte{}, models.MyItem{}, 2
	}
	nItem, _, updtype, ok := models.RemoveItem(client.CurrentChar, item, count)

	pkg := serverpackets.DropItem(clientI, objectId, count, x, y, z)
	client.EncryptAndSend(pkg)

	return nil, nItem, updtype
}
