package main

import (
	"l2gogameserver/config"
	"l2gogameserver/data"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/buff"
	"l2gogameserver/gameserver/idfactory"
	"l2gogameserver/gameserver/models"
	"l2gogameserver/gameserver/models/items"
	"l2gogameserver/gameserver/models/multisell"
	"l2gogameserver/gameserver/models/teleport"
	"l2gogameserver/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//	gameserver.Load()
	//	gameserver.FindPath(-64072, 100856, -3584, -64072, 101048, -3584)

	setup()
	//defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	server.New().Start()
}

func setup() {
	config.LoadAllConfig()
	db.ConfigureDB()
	idfactory.Load()
	multisell.LoadMultisell()
	teleport.LoadLocationListTeleport()
	models.LoadStats()
	models.LoadSkills()
	models.LoadSkillsTrees()
	items.LoadItems()
	buff.LoadCommunityComboBuff()
	models.NewWorld()
	data.Load()
	models.LoadNpc()
}
