package core

import (
	"io"
	"log"
	"os"
	"sync"
)

type LogHandler struct {
}

var Log_Ins *LogHandler
var Log_once sync.Once

func GetLogManager() *LogHandler {
	Log_once.Do(func() {
		Log_Ins = &LogHandler{}
	})
	return Log_Ins
}

func (l *LogHandler) SetLogFile() {
	logFilePath := "Screen_Log.log"
	// log 파일이 없을 경우 log 파일 생성
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		os.Create(logFilePath)
	}
	// log 파일 열기
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// log 패키지를 활요하여 작성할 경우 log 파일에 작성되도록 설정
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
