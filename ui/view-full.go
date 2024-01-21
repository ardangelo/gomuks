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
	"time"

	"go.mau.fi/mauview"
	"go.mau.fi/tcell"

	"maunium.net/go/gomuks/config"
	"maunium.net/go/gomuks/ui/widget"
)

type FullView struct {

	// Views
	roomView     *mauview.Box
	roomListView *RoomList
	flex *mauview.Flex

	focused      mauview.Focusable
	lastFocusTime time.Time

	mainView *MainView
}

func NewFullView(mainView *MainView) *FullView {
	fullView := &FullView{
		flex: mauview.NewFlex().SetDirection(mauview.FlexColumn),
		mainView: mainView,
	}

	fullView.flex.
		AddFixedComponent(mainView.roomListView, 25).
		AddFixedComponent(widget.NewBorder(), 1).
		AddProportionalComponent(mainView.roomView, 1)

	return fullView
}

func (view *FullView) Draw(screen mauview.Screen) {

	// Draw room view only
	if view.mainView.config.Preferences.HideRoomList {
		view.mainView.roomView.Draw(screen)
	
	// Draw entire flex view
	} else {
		view.flex.Draw(screen)
	}
}

func (view *FullView) OnKeyEvent(event mauview.KeyEvent) bool {

	kb := config.Keybind{
		Key: event.Key(),
		Ch:  event.Rune(),
		Mod: event.Modifiers(),
	}
	switch view.mainView.config.Keybindings.Main[kb] {
	case "add_newline":
		return view.flex.OnKeyEvent(tcell.NewEventKey(tcell.KeyEnter, '\n', event.Modifiers()|tcell.ModShift))
	default:
		goto defaultHandler
	}
	return true
defaultHandler:
	return view.flex.OnKeyEvent(event)
}

func (view *FullView) OnMouseEvent(event mauview.MouseEvent) bool {
	return view.flex.OnMouseEvent(event)
}

func (view *FullView) OnPasteEvent(event mauview.PasteEvent) bool {
	return view.flex.OnPasteEvent(event)
}

func (view *FullView) Focus() {
	if view.focused != nil {
		view.focused.Focus()
	}
}

func (view *FullView) Blur() {
	if view.focused != nil {
		view.focused.Blur()
	}
}

