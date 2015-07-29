package telegramapi

import (
  "encoding/json"
  "fmt"
  "net/http"
  "io/ioutil"
  "os"
  "strings"
)

type profileResult struct {
  Ok bool `json:"ok"`
  Result struct {
    ID int `json:"id"`
    FirstName string `json:"first_name"`
    Username string `json:"username"`
  } `json:"result"`
}
type updateResult struct {
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
      os.Exit(1)
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    fmt.Printf("%s", err)
    os.Exit(1)
  }
  return string(contents)
}

func GetUpdates(api_key string) *updateResult {
  json_string := fetch("https://api.telegram.org/" + api_key + "/getUpdates")

  dec := json.NewDecoder(strings.NewReader(json_string))
  result := new(updateResult)
  err := dec.Decode(result)
  if err != nil {
    fmt.Printf("%s\n", err)
  }
  return result
}

