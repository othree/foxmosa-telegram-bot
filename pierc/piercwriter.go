package pierc

import (
  "database/sql"
)
import _ "github.com/go-sql-driver/mysql"

type Message struct {
  Name string
  Time string
  Text string
}

func Writer(schema string, mc <-chan *Message) {
  db, err := sql.Open("mysql",  schema )
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

