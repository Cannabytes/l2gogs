package serverpackets

import (
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/packets"
)

const (
	LEVEL    int32 = 0x01
	EXP      int32 = 0x02
	STR      int32 = 0x03
	DEX      int32 = 0x04
	CON      int32 = 0x05
	INT      int32 = 0x06
	WIT      int32 = 0x07
	MEN      int32 = 0x08
	CUR_HP   int32 = 0x09
	MAX_HP   int32 = 0x0a
	CUR_MP   int32 = 0x0b
	MAX_MP   int32 = 0x0c
	SP       int32 = 0x0d
	CUR_LOAD int32 = 0x0e
	MAX_LOAD int32 = 0x0f
	P_ATK    int32 = 0x11
	ATK_SPD  int32 = 0x12
	P_DEF    int32 = 0x13
	EVASION  int32 = 0x14
	ACCURACY int32 = 0x15
	CRITICAL int32 = 0x16
	M_ATK    int32 = 0x17
	CAST_SPD int32 = 0x18
	M_DEF    int32 = 0x19
	PVP_FLAG int32 = 0x1a
	KARMA    int32 = 0x1b
	CUR_CP   int32 = 0x21
	MAX_CP   int32 = 0x22
)

/**
 * Даный параметр отсылается оффом в паре с MAX_HP
 * Сначала CUR_HP, потом MAX_HP
 */
//var CUR_HP uint8 = 0x09
//var MAX_HP uint8 = 0x0a

/**
 * Даный параметр отсылается оффом в паре с MAX_MP
 * Сначала CUR_MP, потом MAX_MP
 */
//var CUR_MP uint8 = 0x0b
//var MAX_MP uint8 = 0x0c

/**
 * Меняется отображение только в инвентаре, для статуса требуется UserInfo
 */
//var CUR_LOAD uint8 = 0x0e

/**
 * Меняется отображение только в инвентаре, для статуса требуется UserInfo
 */
//var MAX_LOAD uint8 = 0x0f
//var PVP_FLAG uint8 = 0x1a
//var KARMA uint8 = 0x1b

/**
 * Даный параметр отсылается оффом в паре с MAX_CP
 * Сначала CUR_CP, потом MAX_CP
 */
//var CUR_CP uint8 = 0x21
//var MAX_CP uint8 = 0x22

func StatusUpdate(clientI interfaces.ReciverAndSender) []byte {
	char := clientI.(*models.Client).CurrentChar

	buffer := packets.Get()
	defer packets.Put(buffer)

	buffer.WriteSingleByte(0x18)

	buffer.WriteD(char.ObjectId) //Object id
	buffer.WriteD(6)

	buffer.WriteD(CUR_HP)
	buffer.WriteD(char.CurHp)

	buffer.WriteD(MAX_HP)
	buffer.WriteD(char.MaxHp)

	buffer.WriteD(CUR_MP)
	buffer.WriteD(char.CurMp)

	buffer.WriteD(MAX_MP)
	buffer.WriteD(char.MaxMp)

	buffer.WriteD(CUR_CP)
	buffer.WriteD(char.CurCp)

	buffer.WriteD(MAX_CP)
	buffer.WriteD(char.MaxCp)

	return buffer.Bytes()
}
