package notification

import (
	"log"

	"github.com/gen2brain/beeep"
	"github.com/m1kkY8/gochat/src/message"
)

func Notify(msg message.Message, author string) {
	from := msg.Author
	content := msg.Content

	if from == author {
		return
	}

	err := beeep.Notify(from, content, "assets/amogus.png")
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
