var jsonfile = require('jsonfile');
var util = require('util');
var fetch = require('node-fetch');
var fetcher = require('fetch-er');

var token = jsonfile.readFileSync('token.json');

var chat_id = -5425916;

var entrypoint = 'https://api.telegram.org/bot';

var baseurl = `${entrypoint}${token}/`;

var getUpdateBase = `${baseurl}getUpdates`;
var sendStickerBase = `${baseurl}sendSticker`;

var offset = 0;

try {
  offset = jsonfile.readFileSync('offset');
} catch (error) {}

var updateHandler = function (data) {
  "use strict";
  var value = data[0];
  var lower = offset;
  if (value.ok) {
    console.log('fetch update', offset)
    // console.log(JSON.stringify(value, null, 2));
    for (let update of value.result) {
      if (update.update_id < lower) { continue; }
      if (update.update_id > offset) {
        offset = update.update_id + 1;
      }
      let message = update.message;
      let chat = message.chat;
      if (chat.id !== chat_id) { continue; }
      let author = message.from;
      let text = message.text;
      if (!text) { continue; }

      if (/mosa/.test(text)) {
        sendSticker();
      }
      console.log(author, text);
    }
  } else {
  }

  setTimeout(updateFetcher, 5000);
};

var sendSticker = function () {
  fetcher.get(sendStickerBase, {chat_id: chat_id, sticker: 'BQADBQADMgADMqsKAtWHZ4YwVJW-Ag'}).then(updateHandler);
};

var updateFetcher = function () {
  fetcher.get(getUpdateBase, {offset: offset}).then(updateHandler);
};

updateFetcher();
