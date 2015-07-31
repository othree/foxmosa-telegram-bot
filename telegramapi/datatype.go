package telegramapi

type ProfileResult struct {
  Ok bool `json:"ok"`
  Result struct {
    ID int `json:"id"`
    FirstName string `json:"first_name"`
    Username string `json:"username"`
  } `json:"result"`
}
type Sticker struct {
  Width int `json:"width"`
  Height int `json:"height"`
  Thumb struct {
    FileID string `json:"file_id"`
    FileSize int `json:"file_size"`
    Width int `json:"width"`
    Height int `json:"height"`
  } `json:"thumb"`
  FileID string `json:"file_id"`
  FileSize int `json:"file_size"`
}
type Document struct {
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
}
type Photo struct {
  FileID string `json:"file_id"`
  FileSize int `json:"file_size"`
  Width int `json:"width"`
  Height int `json:"height"`
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
    Sticker Sticker `json:"sticker"`
    Document Document `json:"document"`
    Text string `json:"text"`
    Photo []Photo `json:"photo"`
    Caption string `json:"caption"`
  } `json:"message"`
}
type UpdateResult struct {
  Ok bool `json:"ok"`
  Result []Update `json:"result"`
}

