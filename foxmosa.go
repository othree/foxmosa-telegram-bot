package main

import (
  "./config"
  "./telegramapi"
  "fmt"
)

func main() {

  result := telegramapi.GetUpdates(config.API_KEY)

  for _, update := range result.Result {
    fmt.Printf("%d, %s: %s\n", update.UpdateID, update.Message.From.Username, update.Message.Text)
  }
}
