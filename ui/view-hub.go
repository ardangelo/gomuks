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
	//"strings"
	//"time"

	sync "github.com/sasha-s/go-deadlock"

	"go.mau.fi/mauview"
	"go.mau.fi/tcell"

	"maunium.net/go/gomuks/config"
	//"maunium.net/go/gomuks/matrix/rooms"
	"maunium.net/go/gomuks/ui/widget"
	//"maunium.net/go/mautrix/event"

	//"maunium.net/go/gomuks/debug"
)

type HubView struct {
	mauview.Component
	sync.RWMutex

	parent *MainView
}

func NewHubView(mainView *MainView) *HubView {

	// initialize hub
	rstr := &HubView{
		parent: mainView,
	}

	return rstr
}

func (hubView *HubView) Draw(screen mauview.Screen) {

	titleStyle := tcell.StyleDefault.Foreground(tcell.ColorDefault).Bold(true)
	widget.WriteLine(screen, mauview.AlignLeft, "GOMUKS", 2, 1, 0, titleStyle)
}

func (hubView *HubView) OnKeyEvent(event mauview.KeyEvent) bool {
	kb := config.Keybind{
		Key: event.Key(),
		Ch:  event.Rune(),
		Mod: event.Modifiers(),
	}

	switch hubView.parent.config.Keybindings.Roster[kb] {
	case "quit":
		hubView.parent.gmx.Stop(true)
	default:
		return false
	}
	return true
}

func (hubView *HubView) OnMouseEvent(event mauview.MouseEvent) bool {
	if event.HasMotion() {
		return false
	}

	return false
}
