package buffdata

//Баффы/Денцы/Сонги/Дебаффы персонажа
type BuffUser struct {
	Id     int //id skill
	Level  int //skill level
	Second int //Время баффа в секундах (обратный счет)
}
