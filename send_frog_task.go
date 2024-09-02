package main

import (
	"log"
	"os"
)

type SendFrogTask struct {
	tgbot      *TgBot
	chatsStore *ChatsStore
}

func (s *SendFrogTask) SendFrog() {
	chats := s.chatsStore.Chats()
	if len(chats) == 0 {
		return
	}

	file, err := os.CreateTemp(".", "*.jpeg")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		name := file.Name()
		file.Close()
		os.Remove(name)
	}()

	err = DownloadFrogPhoto(file)
	if err != nil {
		log.Println(err)
		return
	}

	for _, chat := range chats {
		file.Seek(0, 0)
		s.tgbot.SendImage(chat, "Среда, мои чуваки!", file)
	}
}

func NewSendFrogTask(tgbot *TgBot, store *ChatsStore) *SendFrogTask {
	return &SendFrogTask{tgbot, store}
}
