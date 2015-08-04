package main

import (
  "./telegramapi"
  "io/ioutil"
  "strconv"
)

func offsetInit() int {
  offset := 0
  data, err := ioutil.ReadFile("offset")
  if err == nil {
    offset, err = strconv.Atoi(string(data))
    if err != nil {
      offset = 0
    }
  }
  return offset
}

func offsetWriter(oc <-chan int) {
  for offset := range oc {
    err := ioutil.WriteFile("offset", []byte(strconv.Itoa(offset)), 0644)
    if err != nil {
      panic(err)
    }
  }
}

func telegram_to_offset(updateChannel <-chan telegramapi.Update, offsetWriterChannel chan<- int) {
  offset := 0
  for update := range updateChannel {
    offset = update.UpdateID + 1
    offsetWriterChannel<-offset
  }
}
