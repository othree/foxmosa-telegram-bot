package main

import (
  "./telegramapi"
  "./pierc"
  "github.com/vaughan0/go-ini"
  "strconv"
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
  chat_id_str, ok := config.Get("telegram", "chat_id")
  chat_id, err := strconv.Atoi(chat_id_str)
  if !ok || err != nil {
    panic("Telegram Chat ID not available.")
  }
  DB_HOST, ok := config.Get("pierc", "db_host")
  if !ok || DB_HOST == "" {
    DB_HOST = "127.0.0.1"
  }
  DB_PORT, ok := config.Get("pierc", "db_port")
  if !ok || DB_PORT == "" {
    DB_PORT = "3306"
  }
  DB_USER, ok := config.Get("pierc", "db_user")
  if !ok {
    panic("Pierc db_user not set.")
  }
  DB_PASS, ok := config.Get("pierc", "db_pass")
  if !ok {
    panic("Pierc db_pass not set.")
  }
  DB_BASE, ok := config.Get("pierc", "db_base")
  if !ok {
    panic("Pierc db_base not set.")
  }


  offsetWriterChannel := make(chan int, 10)
  go offsetWriter(offsetWriterChannel)

  go telegram_to_offset(telegramapi.MakeUpdatesChannel(), offsetWriterChannel)


  messageChannel := make(chan pierc.Message)
  go pierc.Writer(DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_BASE, messageChannel)

  go telegram_to_pierc(telegramapi.MakeUpdatesChannel(), messageChannel)


  telegramapi.TrackingUpdates(token, chat_id, offsetInit())

}
