package main

import (
  "./config"
  "./telegramapi"
  "io/ioutil"
  "strconv"
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

func main() {

  offsetWriterChannel := make(chan int, 10)
  offsetChannel := make(chan int)
  updateChannel := telegramapi.MakeUpdatesChannel()
  go telegramapi.TrackingUpdates(updateChannel, offsetChannel, config.API_KEY)

  var offset int = 0
  offsetChannel<-offset
  go offsetWriter(offsetWriterChannel)

  for updateResult := range updateChannel {
    for _, update  := range updateResult.Result {
      fmt.Printf("%d, %s: %s\n", update.UpdateID, update.Message.From.Username, update.Message.Text)
      offset = update.UpdateID + 1
    }
    offsetWriterChannel<-offset
    offsetChannel<-offset
  }

}
