package main

import (
  "./config"
  "./telegramapi"
  "fmt"
)

func main() {

  offsetChannel := make(chan int)
  updateChannel := telegramapi.MakeUpdatesChannel()
  go telegramapi.TrackingUpdates(updateChannel, offsetChannel, config.API_KEY)

  var offset int = 0
  offsetChannel<-offset

  for updateResult := range updateChannel {
    for _, update  := range updateResult.Result {
      fmt.Printf("%d, %s: %s\n", update.UpdateID, update.Message.From.Username, update.Message.Text)
      offset = update.UpdateID + 1
    }
    offsetChannel<-offset
  }

}
