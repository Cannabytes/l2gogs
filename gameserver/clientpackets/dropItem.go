package clientpackets

import (
	"l2gogameserver/gameserver/drops"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/serverpackets"
	"l2gogameserver/packets"
)

func DropItem(client *models.Client, data []byte) ([]byte, models.MyItem) {
	var read = packets.NewReader(data)
	objectId := read.ReadInt32()
	count := int64(read.ReadInt32())
	_ = read.ReadInt32() // хз
	x := read.ReadInt32()
	y := read.ReadInt32()
	z := read.ReadInt32()

	item := drops.DropItemCharacter(client, objectId, count, x, y, z)

	buffer := packets.Get()
	defer packets.Put(buffer)

	pkg := serverpackets.DropItem(client, objectId, count, x, y, z)
	buffer.WriteSlice(client.CryptAndReturnPackageReadyToShip(pkg))

	return buffer.Bytes(), item
}
