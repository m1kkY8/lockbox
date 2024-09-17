package notification

import (
	"log"
	"regexp"
	"strings"

	"github.com/gen2brain/beeep"
)

func Notify(msg string, username string) {
	reg := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleanedMessage := reg.ReplaceAllString(msg, "")
	partMsg := strings.SplitN(cleanedMessage, ": ", 2)
	getPartUser := strings.Split(partMsg[0], " ")
	fromUser := getPartUser[1]

	if fromUser == username {
		return
	}

	formatMsg := partMsg[1]
	err := beeep.Notify(fromUser, formatMsg, "assets/amogus.png")
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
