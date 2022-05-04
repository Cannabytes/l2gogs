package models

import (
	"encoding/json"
	"l2gogameserver/config"
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver/interfaces"
	"os"
)

type Skill struct {
	SkillName   string `json:"skill_name"`
	SkillId     int    `json:"skill_id"`
	Level       int    `json:"level"`
	OperateType string `json:"operate_type"`
	MagicLevel  int    `json:"magic_level"`
	Effect      struct {
		IPAttack struct {
			Skillid int `json:"skillid"`
			Chance  int `json:"chance"`
			Val     int `json:"val"`
			Val2    int `json:"val2"`
		} `json:"i_p_attack"`
	} `json:"effect"`
	//OperateCond struct {
	//	EquipWeapon []string `json:"equip_weapon"`
	//} `json:"operate_cond"`
	IsMagic            bool `json:"is_magic"`
	MpConsume2         int  `json:"mp_consume2"`
	CastRange          int  `json:"cast_range"`
	EffectiveRange     int  `json:"effective_range"`
	SkillHitTime       int  `json:"skill_hit_time"`
	SkillCoolTime      int  `json:"skill_cool_time"`
	SkillHitCancelTime int  `json:"skill_hit_cancel_time"`
	ReuseDelay         int  `json:"reuse_delay"`
	Attribute          struct {
		Type  string `json:"type"`
		Power int    `json:"power"`
	} `json:"attribute"`
	TargetType  string `json:"target_type"`
	AffectScope string `json:"affect_scope"`
	AffectLimit struct {
		Pvp int `json:"pvp"`
		Pve int `json:"pve"`
	} `json:"affect_limit"`
	NextAction  string `json:"next_action"`
	MultiClass  bool   `json:"multi_class"`
	OlympiadUse bool   `json:"olympiad_use"`
}

var AllSkills []*Skill

func LoadSkills() {
	if config.Get().Debug.EnabledSkills == false {
		return
	}
	logger.Info.Println("Загрузка скиллов")

	file, err := os.Open("./datapack/data/stats/skills/skills.json")
	if err != nil {
		logger.Error.Panicln("Failed to load config file " + err.Error())
	}
	err = json.NewDecoder(file).Decode(&AllSkills)
	if err != nil {
		logger.Error.Panicln("Failed to decode config file " + file.Name() + " " + err.Error())
	}
}

func GetSkillDataInfo(skillId, skilllevel int) (*Skill, bool) {
	for _, skill := range AllSkills {
		if skill.SkillId == skillId && skill.Level == skilllevel {
			return skill, true
		}
	}
	return &Skill{}, false
}

type Trees struct {
	ClassId       int           `json:"classid"`                 // Id класса
	ParentClassId int           `json:"parentClassId,omitempty"` // -1 означает что это отсутствие родительского класса
	Skills        []TreesSkills `json:"skills"`                  //Скиллы класса
}
type TreesSkills struct {
	Name         string `json:"name"`
	SkillId      int    `json:"skillId"`
	SkillLvl     int    `json:"skillLvl"`
	GetLevel     int    `json:"getLevel,omitempty"`
	Sp           int    `json:"sp,omitempty"`
	AutoLearning bool   `json:"autoLearning,omitempty"`
	LearnedByNpc bool   `json:"learnedByNpc,omitempty"`
}

var SkillTrees []Trees

//Загрузка древа скиллов
func LoadSkillsTrees() {
	if config.Get().Debug.EnabledSkills == false {
		return
	}
	logger.Info.Println("Загрузка скиллов (SkillsTrees)")
	file, err := os.Open("./datapack/data/stats/skill_trees/treesSkills.json")
	if err != nil {
		logger.Error.Panicln("Failed to load config file " + err.Error())
	}
	err = json.NewDecoder(file).Decode(&SkillTrees)
	if err != nil {
		logger.Error.Panicln("Failed to decode config file " + file.Name() + " " + err.Error())
	}
}

//Удаление дубликатов скиллов
func dubpicateSkillList(SkillList []TreesSkills) []TreesSkills {
	//var ok bool
	var uniqueSkillChar []TreesSkills
	var userIdSkills []int
	for _, skill := range SkillList {
		userIdSkills = append(userIdSkills, skill.SkillId)
	}

	for _, skillId := range removeDuplicateInt(userIdSkills) {
		skillinfo := maxSkillLevel(SkillList, skillId)
		uniqueSkillChar = append(uniqueSkillChar, skillinfo)
	}
	return uniqueSkillChar
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func maxSkillLevel(SkillList []TreesSkills, skillid int) TreesSkills {
	msklvl := TreesSkills{}
	for _, skill := range SkillList {
		if skill.SkillId == skillid {
			if skill.SkillLvl > msklvl.SkillLvl {
				msklvl = skill
			}
		}
	}
	return msklvl
}

// GetLevelSkills Возвращает все скиллы персонажа, который соответствует уровню и классу
// Необходимо эту функцию юзать при повышении уровня
// Есть смысл сохранять это в БД только в случае если на сервере НЕ АВТОИЗУЧЕНИЕ скилов.
func GetLevelSkills(clientI interfaces.ReciverAndSender) {
	client, ok := clientI.(*Client)
	if !ok {
		panic(ok)
	}

	classId := int(client.CurrentChar.ClassId)
	charLevel := int(client.CurrentChar.Level)

	var all []TreesSkills
	userClassSkill, parentClassId := getSkillClassParent(classId, charLevel)
	all = append(all, userClassSkill...)
	if parentClassId != -1 {
		userClassSkill, parentClassId = getSkillClassParent(parentClassId, charLevel)
		all = append(all, userClassSkill...)
		if parentClassId != -1 {
			userClassSkill, parentClassId = getSkillClassParent(parentClassId, charLevel)
			all = append(all, userClassSkill...)
		}
	}
	all = append(all, userClassSkill...)

	for _, skills := range dubpicateSkillList(all) {
		s, _ := GetSkillDataInfo(skills.SkillId, skills.SkillLvl)
		client.CurrentChar.Skills = append(client.CurrentChar.Skills, *s)
	}

}

func GetSkillName(skillname string) (Skill, bool) {
	for _, skill := range AllSkills {
		if skill.SkillName == skillname {
			return *skill, true
		}
	}
	return Skill{}, false
}

// Возвращает скиллы класса
func getSkillClassParent(classId, char_level int) ([]TreesSkills, int) {
	var uniqueTreesSkills []Trees
	var uniqueSkills []TreesSkills
	parent := -1
	for _, trees := range SkillTrees {
		if trees.ClassId == classId {
			parent = trees.ParentClassId
			uniqueTreesSkills = append(uniqueTreesSkills, trees)
			break
		}
	}

	for _, uniq := range uniqueTreesSkills {
		for _, sk := range uniq.Skills {
			if sk.SkillLvl <= char_level {
				uniqueSkills = append(uniqueSkills, sk)
			}
		}
	}

	return uniqueSkills, parent
}

/*
func GetMySkills(charId int32) []Skill {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	rows, err := dbConn.Query(context.Background(), "SELECT skill_id, skill_level FROM character_skills WHERE char_id = $1", charId)
	if err != nil {
		logger.Error.Panicln(err)
	}

	var skills []Skill
	for rows.Next() {
		var skl Tuple

		err = rows.Scan(&skl.Id, &skl.Lvl)
		if err != nil {
			logger.Info.Println(err)
		}
		sk, ok := AllSkills[skl]
		if !ok {
			logger.Error.Panicln("not found Skill")
		}
		skills = append(skills, sk)
	}
	return skills
}
*/

func (c *Character) LoadSkills() {
	//c.Skills = []Skill{}
	//dbConn, err := db.GetConn()
	//if err != nil {
	//	logger.Error.Panicln(err)
	//}
	//defer dbConn.Release()
	//
	//rows, err := dbConn.Query(context.Background(), "SELECT skill_id,skill_level FROM character_skills WHERE char_id=$1 AND class_id=$2", c.ObjectId, c.ClassId)
	//if err != nil {
	//	logger.Error.Panicln(err)
	//}

	//for rows.Next() {
	//var t Tuple
	//err = rows.Scan(&t.Id, &t.Lvl)
	//if err != nil {
	//	logger.Error.Panicln(err)
	//}

	//c.Skills[t.Id] = t.
	//}

}
