package main

import (
  "./telegramapi"
  "./pierc"
  "github.com/vaughan0/go-ini"
  "unicode/utf8"
  "io/ioutil"
  "strconv"
  "strings"
  "time"
  "fmt"
)

func offsetWriter(oc <-chan int) {
  for offset := range oc {
    err := ioutil.WriteFile("offset", []byte(strconv.Itoa(offset)), 0644)
    if err != nil {
      panic(err)
    }
  }
}

// http://maiyang.github.io/golang/%E5%AD%97%E7%AC%A6%E4%B8%B2/emoji/%E8%A1%A8%E6%83%85/2015/06/16/golang-character-length/
func FilterEmoji(content string) string {
    new_content := ""
    for _, value := range content {
        _, size := utf8.DecodeRuneInString(string(value))
        if size <= 3 {
            new_content += string(value)
        }
    }
    return new_content
  }

func main() {

  config, err := ini.LoadFile("config.ini")
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

  var offset int = 0
  data, err := ioutil.ReadFile("offset")
  if err == nil {
    offset, err = strconv.Atoi(string(data))
    if err != nil {
      offset = 0
    }
  }
  offsetChannel<-offset
  go offsetWriter(offsetWriterChannel)

  messageChannel := make(chan pierc.Message)
  go pierc.Writer(DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_BASE, messageChannel)

  for updateResult := range updateChannel {
    for _, update  := range updateResult.Result {
      author := update.Message.From
      name := strings.Join([]string{author.FirstName, author.LastName}, " ")
      if len(name) > 64 {
        name = name[0:64]
      }
      tm := time.Unix(update.Message.Date, 0)
      fmt.Printf("%d, %s: %s\n", update.UpdateID, name, update.Message.Text)
      offset = update.UpdateID + 1
      msg := pierc.Message{name, tm.Format("2006-01-02 15:04:05"), FilterEmoji(update.Message.Text)}
      messageChannel<-msg
    }
    offsetWriterChannel<-offset
    offsetChannel<-offset
  }

}
