package models

import (
	"context"
	"encoding/json"
	"l2gogameserver/config"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
	"l2gogameserver/gameserver/interfaces"
	"os"
)

/*
type Skill struct {
	ID          int                `json:"id"`
	Levels      int                `json:"levels"`
	Name        string             `json:"name"`
	Power       int                `json:"power"`
	CastRange   int                `json:"castRange"`
	CoolTime    int                `json:"coolTime"`
	HitTime     int                `json:"hitTime"`
	OverHit     bool               `json:"overHit"`
	ReuseDelay  int                `json:"reuseDelay"`
	OperateType skills.OperateType `json:"operateType"`
	TargetType  targets.TargetType `json:"targetType"`
	IsMagic     int                `json:"isMagic"`
	MagicLvl    int                `json:"magicLvl"`
	MpConsume1  int                `json:"mpConsume1"`
	MpConsume2  int                `json:"mpConsume2"`
}

type SkillForParseJSON struct {
	ID          int                `json:"id"`
	Levels      int                `json:"levels"`
	Name        string             `json:"name"`
	Power       []int              `json:"power"`
	CastRange   int                `json:"castRange"`
	CoolTime    int                `json:"coolTime"`
	HitTime     int                `json:"hitTime"`
	OverHit     bool               `json:"overHit"`
	ReuseDelay  int                `json:"reuseDelay"`
	OperateType skills.OperateType `json:"operateType"`
	TargetType  targets.TargetType `json:"targetType"`
	IsMagic     int                `json:"isMagic"`
	MagicLvl    []int              `json:"magicLvl"`
	MpConsume1  []int              `json:"mpConsume1"`
	MpConsume2  []int              `json:"mpConsume2"`
}
*/
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

//var AllSkills map[Tuple]AllSkill

//var AllSkills map[Tuple]Skill

type Tuple struct {
	Id  int
	Lvl int
}

func LoadSkills() {
	if config.Get().Debug.EnabledSkills == false {
		return
	}
	logger.Info.Println("Загрузка скиллов")

	file, err := os.Open("./datapack/data/stats/skills/skills.json")
	if err != nil {
		logger.Error.Panicln("Failed to load config file " + err.Error())
	}
	var skillsJson []Skill
	err = json.NewDecoder(file).Decode(&skillsJson)
	if err != nil {
		logger.Error.Panicln("Failed to decode config file " + file.Name() + " " + err.Error())
	}

	//AllSkills = make(map[Tuple]AllSkill)
	//for _, v := range skillsJson {
	//	fSkill := v
	//	if v.Level > 1 {
	//		for i := 0; i < v.Level; i++ {
	//			fSkill.Level = i
	//			AllSkills[Tuple{v.SkillId, i}] = fSkill
	//		}
	//	} else {
	//		AllSkills[Tuple{v.SkillId, v.Level}] = fSkill
	//	}
	//}
	//
	//qw := AllSkills
	//_ = qw
}

type Trees struct {
	ClassId       int           `json:"classid"`
	ParentClassId int           `json:"parent_class_id,omitempty"`
	Skills        []TreesSkills `json:"skills"`
}
type TreesSkills struct {
	Name         string `json:"name"`
	SkillId      int    `json:"skillId"`
	SkillLvl     int    `json:"skillLvl"`
	GetLevel     int    `json:"getLevel"`
	Sp           int    `json:"sp"`
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
	var unique []TreesSkills

	skillGet := func(unique []TreesSkills, id int) (TreesSkills, int, bool) {
		for index, skill := range unique {
			if skill.SkillId == id {
				return skill, index, true
			}
		}
		return TreesSkills{}, 0, false
	}

	for _, skill := range SkillList {
		uniskill, index, ok := skillGet(unique, skill.SkillId)
		if ok { //Если найден такой скилл уже
			if uniskill.SkillLvl > skill.SkillLvl {
				unique = append(unique[:index], unique[index+1:]...)
				unique = append(unique, skill)
			}
		} else {
			unique = append(unique, skill)
		}
	}

	return unique
}

// GetLevelSkills Возвращает все скиллы персонажа, который соответствует уровню и классу
func GetLevelSkills(clientI interfaces.ReciverAndSender) {
	client, ok := clientI.(*Client)
	if !ok {
		panic(ok)
	}

	classId := int(client.CurrentChar.ClassId)
	charLevel := int(client.CurrentChar.Level)

	var all []TreesSkills
	classSkill, tr, ok := getSkillClassParent(classId, charLevel)
	if ok {
		all = append(all, classSkill...)
	}
	if tr.ParentClassId == -1 {
		return
	}

	classSkill, tr, ok = getSkillClassParent(tr.ParentClassId, charLevel)
	if ok {
		all = append(all, classSkill...)
	}
	if tr.ParentClassId == -1 {
		return
	}

	classSkill, tr, ok = getSkillClassParent(tr.ParentClassId, charLevel)
	if ok {
		all = append(all, classSkill...)
	}
	if tr.ParentClassId == -1 {
		return
	}

	all = dubpicateSkillList(all)

	for _, skills := range all {
		AllSkills := Skill{
			SkillId: skills.SkillId,
			Level:   skills.SkillLvl,
		}

		client.CurrentChar.Skills = append(client.CurrentChar.Skills, AllSkills)
	}
}

func getSkillClassParent(classId, char_level int) ([]TreesSkills, Trees, bool) {
	var t []TreesSkills
	for _, trees := range SkillTrees {
		if trees.ClassId == classId {
			if trees.ParentClassId == -1 {
				return t, trees, false
			}
			for _, skill := range trees.Skills {
				if skill.GetLevel >= char_level {
					t = append(t, skill)
				}
			}
			return t, trees, true
		}
	}
	return t, Trees{}, false
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
	c.Skills = []Skill{}
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	rows, err := dbConn.Query(context.Background(), "SELECT skill_id,skill_level FROM character_skills WHERE char_id=$1 AND class_id=$2", c.ObjectId, c.ClassId)
	if err != nil {
		logger.Error.Panicln(err)
	}

	for rows.Next() {
		var t Tuple
		err = rows.Scan(&t.Id, &t.Lvl)
		if err != nil {
			logger.Error.Panicln(err)
		}

		//c.Skills[t.Id] = t.
	}

}
