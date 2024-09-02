package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	store, err := LoadDefaultChatsStore()
	defer func(store *ChatsStore) {
		err := store.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(store)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tgbot := NewTgBot(os.Getenv("WB_TG_TOKEN"), ctx)
	tgbot.store = store
	tgbot.RegisterHandlers()

	/*scheduler := NewScheduler()
	  sendfrog := NewSendFrogTask(tgbot, store)
	  err = scheduler.ScheduleTaskAndStart(sendfrog.SendFrog)
	  if err != nil {
	  	panic(err)
	  }*/

	tgbot.Start()
}
