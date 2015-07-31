package main

import (
  "./telegramapi"
  "./pierc"
  "github.com/vaughan0/go-ini"
)

func main() {

  config, err := ini.LoadFile("config.ini")
  if err != nil {
    panic("Config file not loaded.")
  }
  token, ok := config.Get("telegram", "token")
  if !ok {
    panic("Telegram API token not available.")
  }
  DB_USER, ok := config.Get("pierc", "db_user")
  if !ok {
    panic("Pierc db_user not set.")
  }
  DB_PASS, ok := config.Get("pierc", "db_pass")
  if !ok {
    panic("Pierc db_pass not set.")
  }
  DB_HOST, ok := config.Get("pierc", "db_host")
  if !ok {
    panic("Pierc db_host not set.")
  }
  DB_BASE, ok := config.Get("pierc", "db_base")
  if !ok {
    panic("Pierc db_base not set.")
  }


  offset := offsetInit()
  offsetWriterChannel := make(chan int, 10)
  go offsetWriter(offsetWriterChannel)

  go telegram_to_offset(telegramapi.MakeUpdatesChannel(), offsetWriterChannel)


  messageChannel := make(chan *pierc.Message)
  go pierc.Writer(DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_BASE, messageChannel)

  go telegram_to_pierc(telegramapi.MakeUpdatesChannel(), messageChannel)


  telegramapi.TrackingUpdates(token, offset)

}
