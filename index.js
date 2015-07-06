var jsonfile = require('jsonfile');
var util = require('util');
var fetcher = require('fetch-er');

var token = jsonfile.readFileSync('token.json');

var chat_id = -5425916;

var entrypoint = 'https://api.telegram.org/bot';

var baseurl = `${entrypoint}${token}/`;

var getUpdateBase = `${baseurl}getUpdates`;
var sendStickerBase = `${baseurl}sendSticker`;

var knex = require('knex')({
  client: 'mysql',
  connection: jsonfile.readFileSync('db.json')
});

var offset = 0;

try {
  offset = jsonfile.readFileSync('offset');
} catch (error) {}

var updateHandler = function (data) {
  "use strict";
  var value = data[0];
  var lower = offset;
  if (value.ok) {
    var records = [];
    console.log((new Date()).getTime(), 'fetch update', offset)
    // console.log(JSON.stringify(value, null, 2));
    for (let update of value.result) {
      if (update.update_id < lower) { continue; }
      if (update.update_id >= offset) {
        offset = update.update_id + 1;
      }
      let message = update.message;
      let chat = message.chat;
      if (chat.id !== chat_id) { continue; }
      let author = `${message.from.first_name || ''} ${message.from.last_name || ''}`.trim();
      let text = message.text;
      if (!text) { continue; }

      if (/mosa/.test(text)) {
        sendSticker();
      }
      // console.log((new Date()).getTime(), update.update_id, author, text);
      records.push({
        channel: 'moztw-telegram',
        name: author,
        time: (new Date(message.date)).toISOString().slice(0, 19).replace('T', ' '),
        message: text,
        type: 'pubmsg',
        hidden: 'F'
      });
    }
    if (records.length) {
      knex('main').insert(records).then(function () {
      });
      jsonfile.writeFile('offset', offset);
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
