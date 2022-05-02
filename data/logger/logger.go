package logger

//Логирование всех данных
//TODO: в будущем ERROR сохранять в файл.

import (
	"encoding/json"
	"log"
	"os"
)

var (
	Warning *log.Logger
	Info    *log.Logger
	Error   *log.Logger
	debug   *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "ERROR: ", log.Ltime|log.Lshortfile)
	debug = log.New(os.Stdout, "DEBUG: ", log.Ltime|log.Lshortfile)
}

func Debug(varstrct interface{}) {
	s, _ := json.MarshalIndent(varstrct, "", "\t")
	debug.Println(string(s))
}
