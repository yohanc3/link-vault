package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	util "wisdom-hoard/util"

	"github.com/rs/zerolog"
)

func InitLogger(name string) zerolog.Logger {

	var fileName string 

	dir, err := os.Getwd()

	if err != nil{
		panic("Not able to get current working directory")
	}

	fileName = dir + "/"+ name +".log"

	file, err := os.OpenFile(
	 	fileName,	
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
    0664,
	)
	
	if err != nil{
		panic("Not able to open file at " + fileName)
	}

	var logger zerolog.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out: file,
		NoColor: true,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatTimestamp: func(i interface{}) string {
			return util.CurrentTime()
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s |", i)
		},
		FormatCaller: func(i interface{}) string {
			return filepath.Base(fmt.Sprintf("%s", i))
		},

}).
	With().
	Timestamp().
	Caller().
	Logger()

	return logger
	
}

var GeneralLogger zerolog.Logger = InitLogger("general")
var StorageLogger zerolog.Logger = InitLogger("storage")
