package serverpackets

import (
	"context"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/packets"
)

func CharSelectionInfo(clientI interfaces.ReciverAndSender) []byte {
	client, ok := clientI.(*models.Client)
	if !ok {
		return []byte{}
	}

	buffer := packets.Get()

	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()
	rows, err := dbConn.Query(context.Background(), `SELECT "login", object_id, char_name, "level", cur_hp, cur_mp, face, hair_style, hair_color, sex, x, y, z, "exp", sp, karma, pvp_kills, pk_kills, clan_id, race, class_id, base_class, title, online_time, nobless, vitality, is_admin, name_color, title_color FROM characters WHERE Login = $1`, client.Account.Login)
	if err != nil {
		logger.Error.Panicln(err)
	}
	client.Account.Char = client.Account.Char[:0]
	for rows.Next() {
		var character = models.GetNewCharacterModel()
		var coord models.Coordinates
		isAdmin := false
		objectID := 0
		accountName := ""
		playerName := ""
		title := ""
		level := 0
		classID := 0
		baseClassID := 0
		onlineTime := 0
		var userExp int32 = 0
		err = rows.Scan(
			&accountName,
			&objectID,
			&playerName,
			&level,
			//&character.MaxHp, //Диприкейтед: мы макс ХП,МП получаем исходя из уровня, скиллов.
			&character.CurHp,
			//&character.MaxMp, //Диприкейтед: мы макс ХП,МП получаем исходя из уровня, скиллов.
			&character.CurMp,
			&character.Face,
			&character.HairStyle,
			&character.HairColor,
			&character.Sex,
			&coord.X,
			&coord.Y,
			&coord.Z,
			&userExp,
			&character.Sp,
			&character.Karma,
			&character.PvpKills,
			&character.PkKills,
			&character.ClanId,
			&character.Race,
			&classID,
			&baseClassID,
			&title,
			&onlineTime,
			&character.Nobless,
			&character.Vitality,
			&isAdmin,
			&character.NameColor,
			&character.TitleColor,
		)
		character.SetObjectID(objectID)
		character.SetAccountName(accountName)
		character.SetPlayerName(playerName)
		character.SetLevel(level)
		character.SetTitle(title)
		character.SetAdmin(isAdmin)
		character.SetBaseClassID(baseClassID)
		character.SetClassID(classID)
		character.SetOnlineTime(onlineTime)
		character.SetExp(userExp)

		if err != nil {
			logger.Error.Panicln(err)
		}
		character.Coordinates = &coord
		character.Conn = client
		client.Account.Char = append(client.Account.Char, character)
	}

	buffer.WriteSingleByte(0x09)
	buffer.WriteD(int32(len(client.Account.Char))) //size char in account

	// Can prevent players from creating new characters (if 0); (if 1, the client will ask if chars may be created (0x13) Response: (0x0D) )
	buffer.WriteD(7)          //char max number
	buffer.WriteSingleByte(0) // delim

	//todo блок который должен повторяться

	for _, char := range client.Account.Char {
		char.ResetHpMpStatLevel()

		buffer.WriteS(char.PlayerName()) // Pers name

		buffer.WriteD(char.ObjectID())    // objId
		buffer.WriteS(char.AccountName()) // loginName

		buffer.WriteD(0)           //TODO sessionId
		buffer.WriteD(char.ClanId) //clanId
		buffer.WriteD(0)           // Builder Level

		buffer.WriteD(char.Sex)           //sex
		buffer.WriteD(int32(char.Race))   // race
		buffer.WriteD(char.BaseClassID()) // baseclass

		buffer.WriteD(1) // active ??

		x, y, z := char.GetXYZ()
		buffer.WriteD(x) //x 53
		buffer.WriteD(y) //y 57
		buffer.WriteD(z) //z 61

		buffer.WriteF(float64(char.CurHp)) //currentHP
		buffer.WriteF(float64(char.CurMp)) //currentMP

		buffer.WriteD(char.Sp)                                                   // SP
		buffer.WriteQ(int64(char.EXP()))                                         // EXP
		buffer.WriteF(char.GetPercentFromCurrentLevel(char.EXP(), char.Level())) // percent
		buffer.WriteD(char.Level())                                              // level

		buffer.WriteD(char.Karma)    // karma
		buffer.WriteD(char.PkKills)  // pk
		buffer.WriteD(char.PvpKills) //pvp

		buffer.WriteD(0)
		buffer.WriteD(0)
		buffer.WriteD(0)
		buffer.WriteD(0)
		buffer.WriteD(0)
		buffer.WriteD(0)
		buffer.WriteD(0)

		paperdoll := char.LoadingVisibleInventory()

		for _, slot := range models.GetPaperdollOrder() {
			buffer.WriteD(int32(paperdoll[slot].Id))
		}

		buffer.WriteD(char.HairStyle) //hairStyle
		buffer.WriteD(char.HairColor) //hairColor
		buffer.WriteD(char.Face)      // face

		buffer.WriteF(char.MaxHP()) //max hp
		buffer.WriteF(char.MaxMP()) // max mp

		buffer.WriteD(0)              // days left before
		buffer.WriteD(char.ClassID()) //classId

		buffer.WriteD(1)          //auto-selected
		buffer.WriteSingleByte(0) // enchanted
		buffer.WriteD(0)          //augumented

		buffer.WriteD(0) // Currently on retail when you are on character select you don't see your transformation.

		// Implementing it will be waster of resources.
		buffer.WriteD(0)             // Pet ID
		buffer.WriteD(0)             // Pet Level
		buffer.WriteD(0)             // Pet Max Food
		buffer.WriteD(0)             // Pet Current Food
		buffer.WriteF(0)             // Pet Max HP
		buffer.WriteF(0)             // Pet Max MP
		buffer.WriteD(char.Vitality) // H5 Vitality

	}

	defer packets.Put(buffer)
	return buffer.Bytes()
}
