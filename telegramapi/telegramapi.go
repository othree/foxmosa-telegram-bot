package telegramapi

import (
  "encoding/json"
  "fmt"
  "net/http"
  "io/ioutil"
  "time"
  "strings"
  "strconv"
)

func fetch(url string) string {
  response, err := http.Get(url)
  if err != nil {
      fmt.Printf("%s", err)
      return "{}"
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    fmt.Printf("%s", err)
  }
  return string(contents)
}

func GetUpdates(api_key string, offset int) *UpdateResult {
  json_string := fetch("https://api.telegram.org/bot" + api_key + "/getUpdates?offset=" + strconv.Itoa(offset) )

  dec := json.NewDecoder(strings.NewReader(json_string))
  result := new(UpdateResult)
  err := dec.Decode(result)
  if err != nil {
    fmt.Printf("%s\n", err)
  }
  return result
}

var channels []chan<- Update

func MakeUpdatesChannel() chan Update {
  uc := make(chan Update)
  channels = append(channels, uc)
  return uc
}

func TrackingUpdates(api_key string, chat_id int, offset int) {
  for {
    updateResult := GetUpdates(api_key, offset)
    for _, update  := range updateResult.Result {
      if update.Message.Chat.ID == chat_id && update.UpdateID >= offset {
        for _, channel := range channels {
          channel <- update
        }
        offset = update.UpdateID + 1
      }
    }
    time.Sleep(5 * time.Second)
  }
}

