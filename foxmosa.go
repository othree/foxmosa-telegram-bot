package main

import (
  "./config"
  "./telegramapi"
  "io/ioutil"
  "strconv"
  "strings"
  "time"
  "fmt"
  "database/sql"
)
import _ "github.com/go-sql-driver/mysql"

type Message struct {
  Name string
  Time string
  Text string
}

func dbWriter(mc <-chan Message) {
  db, err := sql.Open("mysql",  config.DB_USER + ":" + config.DB_PASS + "@" + config.DB_HOST + "/" + config.DB_BASE )
  if err != nil {
    panic(err)
  }
  stmtIns, err := db.Prepare("INSERT INTO main (channel, name, time, message, type, hidden) VALUES( ?, ?, ?, ?, ?, ? )")
  if err != nil {
    panic(err)
  }
  defer stmtIns.Close()
  for message := range mc {
    _, err = stmtIns.Exec("moztw-telegram", message.Name, message.Time, message.Text, "pubmsg", "F")
    if err != nil {
      panic(err)
    }
  }
}


func offsetWriter(oc <-chan int) {
  for offset := range oc {
    err := ioutil.WriteFile("offset", []byte(strconv.Itoa(offset)), 0644)
    if err != nil {
      panic(err)
    }
  }
}

func main() {

  offsetWriterChannel := make(chan int, 10)
  offsetChannel := make(chan int)
  updateChannel := telegramapi.MakeUpdatesChannel()
  go telegramapi.TrackingUpdates(updateChannel, offsetChannel, config.API_KEY)

  var offset int = 0
  offsetChannel<-offset
  go offsetWriter(offsetWriterChannel)

  messageChannel := make(chan Message)
  go dbWriter(messageChannel)

  for updateResult := range updateChannel {
    for _, update  := range updateResult.Result {
      author := update.Message.From
      name := strings.Join([]string{author.FirstName, author.LastName}, " ")
      tm := time.Unix(update.Message.Date, 0)
      fmt.Printf("%d, %s: %s\n", update.UpdateID, name, update.Message.Text)
      offset = update.UpdateID + 1
      msg := Message{name, tm.Format("%Y-%m-%d %H:%M:%S"), update.Message.Text}
      messageChannel<-msg
    }
    offsetWriterChannel<-offset
    offsetChannel<-offset
  }

}
