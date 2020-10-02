package clientpackets

import (
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/packets"
)

type CharCreate struct {
	Name      pgtype.Bytea
	Race      int32
	Sex       uint8
	ClassId   int32
	Int       int32
	Str       int32
	Con       int32
	Men       int32
	Dex       int32
	Wit       int32
	HairStyle uint8
	HairColor uint8
	Face      uint8
	X         int32
	Y         int32
	Z         int32
	MaxHp     int32
	CutHp     int32
	MaxMp     int32
	CurMp     int32
}

func NewCharacterCreate(data []byte, db *pgx.Conn, login string) int32 {
	var packet = packets.NewReader(data)
	var charCreate CharCreate

	charCreate.Name.Bytes = []byte(packet.ReadString())

	charCreate.Race = packet.ReadInt32()
	charCreate.Sex = byte(packet.ReadInt32())
	charCreate.ClassId = packet.ReadInt32()
	charCreate.Int = packet.ReadInt32()
	charCreate.Str = packet.ReadInt32()
	charCreate.Con = packet.ReadInt32()
	charCreate.Men = packet.ReadInt32()
	charCreate.Dex = packet.ReadInt32()
	charCreate.Wit = packet.ReadInt32()
	charCreate.HairStyle = byte(packet.ReadInt32())
	charCreate.HairColor = byte(packet.ReadInt32())
	charCreate.Face = byte(packet.ReadInt32())
	reason := charCreate.validate(db, login)

	return reason
}

var (
	ReasonCreationFailed      int32 = 0x00
	ReasonTooManyCharacters   int32 = 0x01
	ReasonNameAlreadyExists   int32 = 0x02
	Reason16EngChars          int32 = 0x03
	ReasonIncorrectName       int32 = 0x04
	ReasonCreateNotAllowed    int32 = 0x05
	REASON_CHOOSE_ANOTHER_SVR int32 = 0x06
	ReasonOk                  int32 = 99
)

func (cc *CharCreate) validate(db *pgx.Conn, login string) int32 {
	lenName := len(cc.Name.Bytes)
	if (lenName < 1) || (lenName > 16) {
		return Reason16EngChars
	}

	if cc.Face > 2 {
		return ReasonCreationFailed
	}

	if ((cc.Sex == 0) && (cc.HairStyle > 4)) || ((cc.Sex) != 0 && (cc.HairStyle > 6)) {
		return ReasonCreationFailed
	}

	if cc.HairStyle > 3 {
		return ReasonCreationFailed
	}

	row := db.QueryRow("(SELECT exists(SELECT char_name from characters WHERE char_name = $1))", cc.Name.Bytes)
	var exist bool
	err := row.Scan(&exist)
	if err != nil {
		return ReasonCreateNotAllowed
	}
	if exist {
		return ReasonNameAlreadyExists
	}

	row = db.QueryRow("SELECT count(*) FROM characters where login = $1", []byte(login))
	var i int
	err = row.Scan(&i)
	if err != nil {
		return ReasonCreateNotAllowed
	}
	if i > 6 {
		return ReasonTooManyCharacters
	}
	spawn := models.GetCreationCoordinates(cc.ClassId)
	_, err = db.Exec("INSERT INTO characters (char_name, race, sex, class_id, hair_style, hair_color, face,x,y,z,login, base_class) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
		cc.Name.Bytes,
		cc.Race,
		cc.Sex,
		cc.ClassId,
		cc.HairStyle,
		cc.HairColor,
		cc.Face,
		spawn.X,
		spawn.Y,
		spawn.Z,
		[]byte(login),
		cc.ClassId)
	if err != nil {
		return ReasonCreateNotAllowed
	}
	return ReasonOk
}
