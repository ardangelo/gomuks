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
	sync "github.com/sasha-s/go-deadlock"

	"go.mau.fi/mauview"
	"go.mau.fi/tcell"

	"maunium.net/go/mautrix/id"

	"maunium.net/go/gomuks/config"
	"maunium.net/go/gomuks/debug"
	"maunium.net/go/gomuks/matrix/rooms"
)

type HubRoomListView struct {
	sync.RWMutex

	parent *MainView

	isFocused bool
	height       int
	width        int

	selected    *rooms.Room
	selectedTag string

	// The item main text color.
	mainTextColor tcell.Color
	// The text color for selected items.
	selectedTextColor tcell.Color
	// The background color for selected items.
	selectedBackgroundColor tcell.Color
}

func NewHubRoomListView(parent *MainView) *HubRoomListView {
	hrlv := &HubRoomListView{
		parent: parent,

		isFocused: false,

		mainTextColor:           tcell.ColorDefault,
		selectedTextColor:       tcell.ColorWhite,
		selectedBackgroundColor: tcell.ColorDarkGreen,
	}

	return hrlv
}

func (hrlv *HubRoomListView) GetView() mauview.FocusableComponent {
	return hrlv
}

func (hrlv *HubRoomListView) Contains(roomID id.RoomID) bool {
	hrlv.RLock()
	defer hrlv.RUnlock()

	return false
}

func (hrlv *HubRoomListView) Add(room *rooms.Room) {
	if room.IsReplaced() {
		debug.Print(room.ID, "is replaced by", room.ReplacedBy(), "-> not adding to room list")
		return
	}
}

func (hrlv *HubRoomListView) AddToTag(tag rooms.RoomTag, room *rooms.Room) {
	hrlv.Lock()
	defer hrlv.Unlock()
}

func (hrlv *HubRoomListView) Remove(room *rooms.Room) {
	hrlv.Lock()
	defer hrlv.Unlock()
}

func (hrlv *HubRoomListView) RemoveFromTag(tag string, room *rooms.Room) {
	hrlv.Lock()
	defer hrlv.Unlock()
}

func (hrlv *HubRoomListView) Bump(room *rooms.Room) {
	hrlv.RLock()
	defer hrlv.RUnlock()
}

func (hrlv *HubRoomListView) Clear() {
	hrlv.Lock()
	defer hrlv.Unlock()
}

func (hrlv *HubRoomListView) SetSelected(tag string, room *rooms.Room) {
	return
}

func (hrlv *HubRoomListView) HasSelected() bool {
	return hrlv.selected != nil
}

func (hrlv *HubRoomListView) Selected() (string, *rooms.Room) {
	return hrlv.selectedTag, hrlv.selected
}

func (hrlv *HubRoomListView) SelectedRoom() *rooms.Room {
	return hrlv.selected
}

func (hrlv *HubRoomListView) First() (string, *rooms.Room) {
	hrlv.RLock()
	defer hrlv.RUnlock()
	return "", nil
}

func (hrlv *HubRoomListView) Last() (string, *rooms.Room) {
	hrlv.RLock()
	defer hrlv.RUnlock()
	return "", nil
}

func (hrlv *HubRoomListView) Previous() (string, *rooms.Room) {
	hrlv.RLock()
	defer hrlv.RUnlock()
	return "", nil
}

func (hrlv *HubRoomListView) Next() (string, *rooms.Room) {
	hrlv.RLock()
	defer hrlv.RUnlock()
	return "", nil
}

// NextWithActivity Returns next room with activity.
//
// Sorted by (in priority):
//
// - Highlights
// - Messages
// - Other traffic (joins, parts, etc)
//
// TODO: Sorting. Now just finds first room with new messages.
func (hrlv *HubRoomListView) NextWithActivity() (string, *rooms.Room) {
	hrlv.RLock()
	defer hrlv.RUnlock()
	return "", nil
}

func (hrlv *HubRoomListView) OnKeyEvent(event mauview.KeyEvent) bool {

	kb := config.Keybind{
		Key: event.Key(),
		Ch:  event.Rune(),
		Mod: event.Modifiers(),
	}
	switch hrlv.parent.config.Keybindings.RoomList[kb] {
	case "next_room":
		hrlv.parent.SwitchRoom(hrlv.Next())
	case "prev_room":
		hrlv.parent.SwitchRoom(hrlv.Previous())
	case "search_rooms":
		hrlv.parent.ShowModal(NewFuzzySearchModal(hrlv.parent, 42, 12))
	case "scroll_up":
		break
	case "scroll_down":
		break
	case "select_room":
		if hrlv.parent.displayState == CompactRoomList {
			hrlv.parent.SetDisplayState(CompactRoom)
		} else {
			hrlv.parent.SetFlexFocused(hrlv.parent.roomView)
		}
	case "back":
		hrlv.parent.gmx.Stop(true)
	default:
		return true
	}
	return true
}

func (hrlv *HubRoomListView) OnPasteEvent(_ mauview.PasteEvent) bool {
	return false
}

func (hrlv *HubRoomListView) OnMouseEvent(event mauview.MouseEvent) bool {
	if event.HasMotion() {
		return false
	}
	switch event.Buttons() {
	case tcell.WheelUp:
		return true
	case tcell.WheelDown:
		return true
	case tcell.Button1:
		return true
	}
	return false
}

func (hrlv *HubRoomListView) Focus() {
	hrlv.isFocused = true
}

func (hrlv *HubRoomListView) Blur() {
	hrlv.isFocused = false
}

// Draw draws this primitive onto the screen.
func (hrlv *HubRoomListView) Draw(screen mauview.Screen) {
	hrlv.RLock()
	defer hrlv.RUnlock()
	hrlv.width, hrlv.height = screen.Size()
}
