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
    author := update.Message.From
    name := strings.TrimSpace(strings.Join([]string{author.FirstName, author.LastName}, " "))
    if len(name) > 64 {
      name = name[0:64]
    }
    text := noemoji.Noemojitize(update.Message.Text)
    if len(update.Message.Photo) > 0 {
      text = strings.Join([]string{"[photo]", noemoji.Noemojitize(update.Message.Caption)}, " ")
    }
    tm := time.Unix(update.Message.Date, 0).Format("2006-01-02 15:04:05")
    fmt.Printf("[%d] %s %s: %s\n", update.UpdateID, tm, name, text)
    msg := pierc.Message{name, tm, text}
    messageChannel<-&msg
  }
}

