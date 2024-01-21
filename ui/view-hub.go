// gomuks - A terminal Matrix client written in Go.
// Copyright (C) 2020 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ui

import (
	"go.mau.fi/mauview"
	"go.mau.fi/tcell"

	"maunium.net/go/gomuks/matrix/rooms"
	"maunium.net/go/gomuks/config"
)

type HubView struct {
	mainView *MainView
}

func NewHubView(mainView *MainView) *HubView {
	hubView := &HubView{
		mainView: mainView,
	}

	return hubView
}

func (view *HubView) SwitchRoom(room *rooms.Room) {
	if room == nil {
		return
	}
}

func (view *HubView) Draw(screen mauview.Screen) {

	view.mainView.roomView.Draw(screen)
}

func (view *HubView) OnKeyEvent(event mauview.KeyEvent) bool {

	kb := config.Keybind{
		Key: event.Key(),
		Ch:  event.Rune(),
		Mod: event.Modifiers(),
	}
	switch view.mainView.config.Keybindings.Main[kb] {
	case "add_newline":
		return view.mainView.roomView.OnKeyEvent(tcell.NewEventKey(tcell.KeyEnter, '\n', event.Modifiers()|tcell.ModShift))
	default:
		goto defaultHandler
	}
	return true
defaultHandler:
	return view.mainView.roomView.OnKeyEvent(event)
}

func (view *HubView) OnMouseEvent(event mauview.MouseEvent) bool {
	return view.mainView.roomView.OnMouseEvent(event)
}

func (view *HubView) OnPasteEvent(event mauview.PasteEvent) bool {
	return view.mainView.roomView.OnPasteEvent(event)
}

func (view *HubView) Focus() {
}

func (view *HubView) Blur() {
}

