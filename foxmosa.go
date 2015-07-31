package main

import (
  "./telegramapi"
  "./pierc"
  "github.com/othree/noemoji"
  "github.com/vaughan0/go-ini"
  "strings"
  "time"
  "fmt"
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


  offsetWriterChannel := make(chan int, 10)
  offsetChannel := make(chan int)
  updateChannel := telegramapi.MakeUpdatesChannel()
  go telegramapi.TrackingUpdates(updateChannel, offsetChannel, token)

  offset := offsetInit()
  offsetChannel<-offset
  go offsetWriter(offsetWriterChannel)

  messageChannel := make(chan pierc.Message)
  go pierc.Writer(DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_BASE, messageChannel)

  for updateResult := range updateChannel {
    for _, update  := range updateResult.Result {
      author := update.Message.From
      name := strings.TrimSpace(strings.Join([]string{author.FirstName, author.LastName}, " "))
      text := noemoji.Noemojitize(update.Message.Text)
      if len(name) > 64 {
        name = name[0:64]
      }
      tm := time.Unix(update.Message.Date, 0)
      fmt.Printf("[%d] %s %s: %s\n", update.UpdateID, tm.Format("2006-01-02 15:04:05"), name, text)
      offset = update.UpdateID + 1
      msg := pierc.Message{name, tm.Format("2006-01-02 15:04:05"), text}
      messageChannel<-msg
    }
    offsetWriterChannel<-offset
    offsetChannel<-offset
  }

}
