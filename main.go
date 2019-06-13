/*
Copyright (C) Marck Tomack <marcktomack@tutanota.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"time"
)

var (
	bot, _ = tb.NewBot(tb.Settings{
		Token:  "YOUR TOKEN HERE",
		Poller: &tb.LongPoller{Timeout: 25 * time.Minute},
	})
)

func onUserJoined() {
	bot.Handle(tb.OnUserJoined, func(m *tb.Message) {
		convID := strconv.Itoa(m.Sender.ID)
		bot.Delete(m)
		inlineBtn := tb.InlineButton{
			Unique: "unique_unique",
			Data:   convID,
			Text:   "I'm not a bot",
		}
		inlineKeys := [][]tb.InlineButton{
			[]tb.InlineButton{inlineBtn},
		}
		User := &tb.ChatMember{
			User: m.Sender,
		}
		restrictUser := &tb.ChatMember{
			User:   m.Sender,
			Rights: tb.Rights{CanSendMessages: false},
		}

		promoteUser := &tb.ChatMember{
			User:   m.Sender,
			Rights: tb.Rights{CanSendMessages: true},
		}
		text := fmt.Sprintf(`<a href="tg://user?id=%v">%v %v</a>`+" <b>if you are not a bot, press the button below, otherwise you will be kicked</b>", m.Sender.ID, m.Sender.FirstName, m.Sender.LastName)
		bot.Restrict(m.Chat, restrictUser)
		msg, _ := bot.Send(m.Chat, text, &tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		}, tb.ModeHTML)

		bot.Handle(&inlineBtn, func(c *tb.Callback) {

			convData, _ := strconv.Atoi(inlineBtn.Data)
			if c.Sender.ID == convData {
				bot.Promote(m.Chat, promoteUser)
				bot.Delete(msg)

			}

		})

		time.AfterFunc(20*time.Second, func() {
			tryToDoThings := bot.Delete(msg)
			if tryToDoThings != nil {
				println("user clicked the button")
			} else if tryToDoThings == nil {
				bot.Ban(m.Chat, User)
				bot.Unban(m.Chat, m.Sender)
			}

		})

	})

}

func main() {
	onUserJoined()
	bot.Start()

}