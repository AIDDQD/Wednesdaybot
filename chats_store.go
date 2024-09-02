package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	DefaultStoreFile = "chats.store"
)

type ChatsStore struct {
	file  *os.File
	chats []int64
}

func (c *ChatsStore) Close() error {
	defer c.file.Close()
	err := c.writeChats()
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatsStore) readChats() {
	scanner := bufio.NewScanner(c.file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		parsedInt, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		c.chats = append(c.chats, parsedInt)
	}
}

func (c *ChatsStore) initFile(filename string) error {
	var err error
	c.file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return err
	}

	c.readChats()

	return err
}

func (c *ChatsStore) writeChats() error {
	if len(c.chats) == 0 {
		return nil
	}

	c.file.Truncate(0)
	c.file.Seek(0, 0)
	writer := bufio.NewWriter(c.file)
	for _, chat := range c.chats {
		_, err := writer.WriteString(fmt.Sprintf("%d\n", chat))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func (c *ChatsStore) Chats() []int64 {
	return c.chats
}

func (c *ChatsStore) HasChat(chat int64) bool {
	for _, ch := range c.chats {
		if chat == ch {
			return true
		}
	}
	return false
}

func (c *ChatsStore) AddChat(chat int64) bool {
	if c.HasChat(chat) {
		return false
	}
	c.chats = append(c.chats, chat)
	return true
}

func LoadChatsStore(storeFile string) (*ChatsStore, error) {
	store := &ChatsStore{}
	err := store.initFile(storeFile)
	return store, err
}

func LoadDefaultChatsStore() (*ChatsStore, error) {
	return LoadChatsStore(DefaultStoreFile)
}
