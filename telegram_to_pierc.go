package main

import (
  "./telegramapi"
  "./pierc"
  "github.com/othree/noemoji"
  "strings"
  "time"
  "fmt"
)

func telegram_to_pierc(updateChannel <-chan *telegramapi.Update, messageChannel chan<- *pierc.Message) {
  for update := range updateChannel {
    message := update.Message
    author := message.From
    name := strings.TrimSpace(strings.Join([]string{author.FirstName, author.LastName}, " "))
    if len(name) > 64 {
      name = name[0:64]
    }
    text := noemoji.Noemojitize(message.Text)
    if len(message.Photo) > 0 {
      text = strings.Join([]string{"[photo]", noemoji.Noemojitize(message.Caption)}, " ")
    }
    if message.Document.FileName != "" {
      document := message.Document
      text = strings.Join([]string{"[file:"+ document.MimeType +"]", noemoji.Noemojitize(document.FileName)}, " ")
    }
    if message.Sticker.FileID != "" {
      sticker := message.Sticker
      val, ok := foxmosaStickerMap[sticker.FileID]
      if ok {
        text = "[foxmosa] " + val
      } else {
        text = "[sticker]" + sticker.FileID
      }
    }
    tm := time.Unix(update.Message.Date, 0).Format("2006-01-02 15:04:05")
    fmt.Printf("[%d] %s %s: %s\n", update.UpdateID, tm, name, text)
    msg := pierc.Message{name, tm, text}
    messageChannel<-&msg
  }
}

