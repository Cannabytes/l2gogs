package interfaces

type Identifier interface {
	GetId() int32
}
type Spectator interface {
	MaxHP() float64
	MaxMP() float64
	MaxCP() float64
	MaxRunSpeed() float64
}
type UniquerId interface {
	ObjectID() int32
}
type Namer interface {
	PlayerName() string
}
type Other interface {
	SetStatusOffline()
}
type Positionable interface {
	SetX(int32)
	SetY(int32)
	SetZ(int32)
	SetXYZ(int32, int32, int32)
	SetHeading(int32)
	SetInstanceId(int32)
	GetX() int32
	GetY() int32
	GetZ() int32
	GetXYZ() (int32, int32, int32)
	GetCurrentRegion() WorldRegioner
	//setLocation(Location)
	//setXYZByLoc(ILocational)
}
type WorldRegioner interface {
	GetNeighbors() []WorldRegioner
	GetCharsInRegion() []CharacterI
	AddVisibleChar(CharacterI)
	GetNpcInRegion() []Npcer
	DeleteVisibleChar(CharacterI)
}
type Npcer interface {
	UniquerId
	Identifier
}

type CharacterI interface {
	Positionable
	Namer
	UniquerId
	Other
	Spectator
	EncryptAndSend(data []byte)
	CloseChannels()
	GetClassId() int32
}
type ReciverAndSender interface {
	Receive() (opcode byte, data []byte, e error)
	AddLengthAndSand(d []byte)
	Send(data []byte)
	EncryptAndSend(data []byte)
	CryptAndReturnPackageReadyToShip(data []byte) []byte
	Player() CharacterI
}
