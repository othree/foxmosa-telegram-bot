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
type UpdateResult struct {
  Ok bool `json:"ok"`
  Result []struct {
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
      Date int `json:"date"`
      Text string `json:"text"`
    } `json:"message"`
  } `json:"result"`
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

func TrackingUpdates(uc chan<- *UpdateResult, oc <-chan int, api_key string) {
  for offset := range oc {
    uc <- GetUpdates(api_key, offset)
    time.Sleep(5 * time.Second)
  }
}

func MakeUpdatesChannel() chan *UpdateResult {
  return make(chan *UpdateResult, 10)
}

