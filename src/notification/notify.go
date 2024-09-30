package notification

import (
	"fmt"
	"log"

	"github.com/gen2brain/beeep"
	"github.com/m1kkY8/gochat/src/message"
)

func Notify(msg message.Message, author string) {
	from := msg.Author
	at := msg.Room
	content := msg.Content

	if from == author {
		return
	}

	err := beeep.Notify(fmt.Sprintf("%s at %s", from, at), string(content), "src/notification/assets/amogus.png")
	if err != nil {
		log.Println(err)
		return
	}
}
