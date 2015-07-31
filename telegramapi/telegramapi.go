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

type ProfileResult struct {
  Ok bool `json:"ok"`
  Result struct {
    ID int `json:"id"`
    FirstName string `json:"first_name"`
    Username string `json:"username"`
  } `json:"result"`
}
type Update struct {
  UpdateID int `json:"update_id"`
  Message struct {
    MessageID int `json:"message_id"`
    From struct {
      ID int `json:"id"`
      FirstName string `json:"first_name"`
      LastName string `json:"last_name"`
      Username string `json:"username"`
    } `json:"from"`
    Chat struct {
      ID int `json:"id"`
      Title string `json:"title"`
    } `json:"chat"`
    Date int64 `json:"date"`
    Document struct {
      FileName string `json:"file_name"`
      MimeType string `json:"mime_type"`
      Thumb struct {
        FileID string `json:"file_id"`
        FileSize int `json:"file_size"`
        Width int `json:"width"`
        Height int `json:"height"`
      } `json:"thumb"`
      FileID string `json:"file_id"`
      FileSize int `json:"file_size"`
    } `json:"document"`
    Text string `json:"text"`
    Photo []struct {
      FileID string `json:"file_id"`
      FileSize int `json:"file_size"`
      Width int `json:"width"`
      Height int `json:"height"`
    } `json:"photo"`
    Caption string `json:"caption"`
  } `json:"message"`
}
type UpdateResult struct {
  Ok bool `json:"ok"`
  Result []Update `json:"result"`
}

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
  json_string := fetch("https://api.telegram.org/" + api_key + "/getUpdates?offset=" + strconv.Itoa(offset) )

  dec := json.NewDecoder(strings.NewReader(json_string))
  result := new(UpdateResult)
  err := dec.Decode(result)
  if err != nil {
    fmt.Printf("%s\n", err)
  }
  return result
}

var channels []chan<- *Update

func MakeUpdatesChannel() chan *Update {
  uc := make(chan *Update, 10)
  channels = append(channels, uc)
  return uc
}

func TrackingUpdates(api_key string, offset int) {
  for {
    updateResult := GetUpdates(api_key, offset)
    for _, update  := range updateResult.Result {
      for _, channel := range channels {
        channel <- &update
      }
      offset = update.UpdateID + 1
    }
    time.Sleep(5 * time.Second)
  }
}

