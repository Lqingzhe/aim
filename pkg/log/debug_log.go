package newlog

import "log"

func LogRed(message string, input ...any) {
	log.Printf("\033[91m "+message+" \033[0m", input...)
}
