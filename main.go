package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	store, err := LoadDefaultChatsStore()
	defer store.Close()
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	go func() {
		<-ch
		store.Close()
		os.Exit(1488)
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tgbot := NewTgBot(os.Getenv("WB_TG_TOKEN"), ctx)
	tgbot.store = store
	tgbot.RegisterHandlers()

	scheduler := NewScheduler()
	sendfrog := NewSendFrogTask(tgbot, store)
	err = scheduler.ScheduleTaskAndStart(sendfrog.SendFrog)
	if err != nil {
		panic(err)
	}

	tgbot.Start()
}
