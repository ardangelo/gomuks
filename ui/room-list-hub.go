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
	"fmt"

	sync "github.com/sasha-s/go-deadlock"

	"go.mau.fi/mauview"
	"go.mau.fi/tcell"

	"maunium.net/go/mautrix/id"
	"maunium.net/go/mautrix/event"

	"maunium.net/go/gomuks/config"
	"maunium.net/go/gomuks/debug"
	"maunium.net/go/gomuks/matrix/rooms"

	"maunium.net/go/gomuks/ui/widget"
)

type HubRoom struct {
	mxRoom *rooms.Room
	initialHistoryLoaded bool
	latestMessage string
}

type HubRoomListView struct {
	sync.RWMutex

	parent *MainView

	isFocused bool
	height       int
	width        int

	hubRooms []*HubRoom
	selected *HubRoom
	selectedTag string
	renderStartIndex int

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
		renderStartIndex: 0,

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
	hrlv.Lock()
	defer hrlv.Unlock()

	return false
}

func (hrlv *HubRoomListView) Add(room *rooms.Room) {
	if room.IsReplaced() {
		return
	}

	hrlv.Lock()
	defer hrlv.Unlock()

	insertAt := len(hrlv.hubRooms)
	for i := 0; i < len(hrlv.hubRooms); i++ {
		if hrlv.hubRooms[i].mxRoom == room {
			return
		} else if room.LastReceivedMessage.After(hrlv.hubRooms[i].mxRoom.LastReceivedMessage) {
			insertAt = i
			break
		}
	}
	hrlv.hubRooms = append(hrlv.hubRooms, nil)
	copy(hrlv.hubRooms[insertAt+1:], hrlv.hubRooms[insertAt:len(hrlv.hubRooms)-1])
	hrlv.hubRooms[insertAt] = &HubRoom{
		mxRoom: room,
		initialHistoryLoaded: false,
		latestMessage: "",
	}
}

func (hrlv *HubRoomListView) AddToTag(tag rooms.RoomTag, room *rooms.Room) {
	return
}

func (hrlv *HubRoomListView) indexOfRoom(room *rooms.Room) int {
	if room == nil {
		return -1
	}

	for index, hubRoom := range hrlv.hubRooms {
		if hubRoom.mxRoom == room {
			return index
		}
	}

	return -1
}

func (hrlv *HubRoomListView) Remove(room *rooms.Room) {
	hrlv.Lock()
	defer hrlv.Unlock()

	index := hrlv.indexOfRoom(room)
	if index < 0 || index > len(hrlv.hubRooms) {
		return
	}

	last := len(hrlv.hubRooms) - 1
	if index < last {
		copy(hrlv.hubRooms[index:], hrlv.hubRooms[index+1:])
	}
	hrlv.hubRooms[last] = nil
	hrlv.hubRooms = hrlv.hubRooms[:last]
}

func (hrlv *HubRoomListView) RemoveFromTag(tag string, room *rooms.Room) {
	return
}

func (hrlv *HubRoomListView) Bump(room *rooms.Room) {
	hrlv.Remove(room)
	hrlv.Add(room)
}

func (hrlv *HubRoomListView) Clear() {
	hrlv.hubRooms = hrlv.hubRooms[:0]
}

func (hrlv *HubRoomListView) SetSelected(tag string, room *rooms.Room) {
	index := hrlv.indexOfRoom(room)
	if index < 0 || index >= len(hrlv.hubRooms) {
		return
	}
	hrlv.selected = hrlv.hubRooms[index]

	// Update render start index
	renderableIndices := (hrlv.height - 2) / 2 - 1
	if renderableIndices < 0 {
		renderableIndices = 1
	}
	if index < hrlv.renderStartIndex {
		hrlv.renderStartIndex = index
	} else if hrlv.renderStartIndex + renderableIndices <= index {
		hrlv.renderStartIndex = index - renderableIndices
	}
	for hrlv.renderStartIndex >= len(hrlv.hubRooms) {
		hrlv.renderStartIndex -= renderableIndices
	}
	if hrlv.renderStartIndex < 0 {
		hrlv.renderStartIndex = 0
	}
	debug.Print("Render start index for ", index, ": ", hrlv.renderStartIndex)
}

func (hrlv *HubRoomListView) HasSelected() bool {
	return hrlv.selected != nil
}

func (hrlv *HubRoomListView) Selected() (string, *rooms.Room) {
	if hrlv.selected == nil {
		return "", nil
	}
	return hrlv.selectedTag, hrlv.selected.mxRoom
}

func (hrlv *HubRoomListView) SelectedRoom() *rooms.Room {
	if hrlv.selected == nil {
		return nil
	}
	return hrlv.selected.mxRoom
}

func (hrlv *HubRoomListView) First() (string, *rooms.Room) {
	hrlv.Lock()
	defer hrlv.Unlock()
	if len(hrlv.hubRooms) > 0 {
		return "", hrlv.hubRooms[0].mxRoom
	}
	return "", nil
}

func (hrlv *HubRoomListView) Last() (string, *rooms.Room) {
	hrlv.Lock()
	defer hrlv.Unlock()
	if len(hrlv.hubRooms) > 0 {
		return "", hrlv.hubRooms[len(hrlv.hubRooms) - 1].mxRoom
	}
	return "", nil
}

// Get HubRoom at index, clamping to valid indices
func (hrlv *HubRoomListView) seekRoom(index int) *HubRoom {

	if len(hrlv.hubRooms) == 0 {
		return nil
	}

	if index < 0 {
		index = 0
	} else if index >= len(hrlv.hubRooms) {
		index = len(hrlv.hubRooms) - 1
	}

	return hrlv.hubRooms[index]
}

func (hrlv *HubRoomListView) Previous() (string, *rooms.Room) {
	hrlv.Lock()
	defer hrlv.Unlock()

	index := 0
	if hrlv.selected != nil {
		index = hrlv.indexOfRoom(hrlv.selected.mxRoom)
	}

	if hrlv.selected = hrlv.seekRoom(index - 1); hrlv.selected != nil {
		return "", hrlv.selected.mxRoom
	}
	return "", nil
}

func (hrlv *HubRoomListView) Next() (string, *rooms.Room) {
	hrlv.Lock()
	defer hrlv.Unlock()

	index := 0
	if hrlv.selected != nil {
		index = hrlv.indexOfRoom(hrlv.selected.mxRoom)
	}

	if hrlv.selected = hrlv.seekRoom(index + 1); hrlv.selected != nil {
		return "", hrlv.selected.mxRoom
	}
	return "", nil
}

func (hrlv *HubRoomListView) NextWithActivity() (string, *rooms.Room) {
	hrlv.RLock()
	defer hrlv.RUnlock()
	for _, hubRoom := range hrlv.hubRooms {
		if hubRoom.mxRoom.HasNewMessages() {
			return "", hubRoom.mxRoom
		}
	}
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
		debug.Print("Switching to next room")
		hrlv.parent.SwitchRoom(hrlv.Next())
	case "prev_room":
		debug.Print("Switching to previous room")
		hrlv.parent.SwitchRoom(hrlv.Previous())
	case "search_rooms":
		hrlv.parent.ShowModal(NewFuzzySearchModal(hrlv.parent, 30, 12))
	case "scroll_up":
		if hrlv.selected != nil {
			index := hrlv.indexOfRoom(hrlv.selected.mxRoom)
			renderableIndices := (hrlv.height - 2) / 2 - 1
			hrlv.parent.SwitchRoom("", hrlv.seekRoom(index - renderableIndices).mxRoom)
		}
	case "scroll_down":
		if hrlv.selected != nil {
			index := hrlv.indexOfRoom(hrlv.selected.mxRoom)
			renderableIndices := (hrlv.height - 2) / 2 - 1
			hrlv.parent.SwitchRoom("", hrlv.seekRoom(index + renderableIndices).mxRoom)
		}
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
	hrlv.width, hrlv.height = screen.Size()

	if hrlv.width == 0 {
		return
	}

	now := time.Now()
	tm := now.Format("15:04")

	headerStyle := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorDarkGreen)
	headerText := fmt.Sprintf("GOMUKS%*s", hrlv.width - 6, tm)
	widget.WriteLine(screen, mauview.AlignLeft, headerText, 0, 0, hrlv.width, headerStyle)

	y := 1
	if len(hrlv.hubRooms) == 0 {
		return
	}

	for _, hubRoom := range hrlv.hubRooms[hrlv.renderStartIndex:] {
		if hubRoom.mxRoom.IsReplaced() {
			continue
		}

		renderHeight := 2
		if y+renderHeight >= hrlv.height {
			renderHeight = hrlv.height - y
		}

		isSelected := hubRoom == hrlv.selected

		style := tcell.StyleDefault.
			Foreground(tcell.ColorDefault).
			Bold(hubRoom.mxRoom.HasNewMessages())
		if isSelected {
			style = style.Italic(true)
			if hrlv.isFocused {
				style = style.Foreground(tcell.ColorBlack).
					Background(tcell.ColorWhite)
			}
		}

		timestamp := hubRoom.mxRoom.LastReceivedMessage
		tm := timestamp.Format("15:04")
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if timestamp.Before(today) {
			if timestamp.Before(today.AddDate(0, 0, -6)) {
				tm = timestamp.Format("2006-01-02")
			} else {
				tm = timestamp.Format("Monday")
			}
		}

		lastMessage, received := hubRoom.GetLatestMessage(hrlv)
		msgStyle := style.Foreground(tcell.ColorGray).Italic(!received)
		startingX := 0

		if isSelected {
			msgStyle = msgStyle.Italic(true)
			if hrlv.isFocused {
				msgStyle = msgStyle.
					Foreground(tcell.ColorBlack).
					Background(tcell.ColorWhite)
			}
			startingX += 2

			widget.WriteLine(screen, mauview.AlignLeft, string(tcell.RuneDiamond)+" ", 0, y, 4, style)
		}

		tmX := hrlv.width - len(tm)
		widget.WriteLinePadded(screen, mauview.AlignLeft, hubRoom.mxRoom.GetTitle(), startingX, y, tmX, style)
		widget.WriteLine(screen, mauview.AlignLeft, tm, tmX, y, startingX+len(tm), style)
		widget.WriteLinePadded(screen, mauview.AlignLeft, "  " + lastMessage, 0, y+1,
			hrlv.width, msgStyle)

		y += renderHeight
		if y >= hrlv.height {
			break
		}
	}
}

func (hubRoom *HubRoom) GetLatestMessage(hrlv *HubRoomListView) (string, bool) {

	roomView, _ := hrlv.parent.getRoomView(hubRoom.mxRoom.ID, true)

	if msgView := roomView.MessageView(); len(msgView.messages) < 20 && !msgView.initialHistoryLoaded {
		msgView.initialHistoryLoaded = true
		go hrlv.parent.LoadHistory(hubRoom.mxRoom.ID)
	}

	if len(roomView.content.messages) > 0 {
		for index := len(roomView.content.messages) - 1; index >= 0; index-- {
			if roomView.content.messages[index].Type == event.MsgText {
				return roomView.content.messages[index].PlainText(), true
			}
		}
	}

	return "", false
}
