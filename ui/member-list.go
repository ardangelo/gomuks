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
	"math"
	"sort"
	"strings"

	//"github.com/mattn/go-runewidth"

	"go.mau.fi/mauview"
	"go.mau.fi/tcell"

	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"maunium.net/go/gomuks/matrix/rooms"
	"maunium.net/go/gomuks/ui/widget"
)

type MemberList struct {

	flex *mauview.Flex
	listView *mauview.TextView

	memberList roomMemberList
}

func NewMemberList() *MemberList {
	ml := &MemberList{
		listView: mauview.NewTextView(),
		flex: mauview.NewFlex().SetDirection(mauview.FlexColumn),
	}

	ml.listView.
		SetScrollable(true).
		SetWrap(false).
		SetDynamicColors(true)

	ml.flex.
		AddFixedComponent(widget.NewBorder(), 1).
		AddProportionalComponent(ml.listView, 1)

	return ml
}

type memberListItem struct {
	rooms.Member
	PowerLevel int
	Sigil      rune
	UserID     id.UserID
	Color      tcell.Color
}

type roomMemberList []*memberListItem

func (rml roomMemberList) Len() int {
	return len(rml)
}

func (rml roomMemberList) Less(i, j int) bool {
	if rml[i].PowerLevel != rml[j].PowerLevel {
		return rml[i].PowerLevel > rml[j].PowerLevel
	}
	return strings.Compare(strings.ToLower(rml[i].Displayname), strings.ToLower(rml[j].Displayname)) < 0
}

func (rml roomMemberList) Swap(i, j int) {
	rml[i], rml[j] = rml[j], rml[i]
}

func (ml *MemberList) Update(data map[id.UserID]*rooms.Member, levels *event.PowerLevelsEventContent) *MemberList {
	ml.memberList = make(roomMemberList, len(data))
	i := 0
	highestLevel := math.MinInt32
	count := 0
	for _, level := range levels.Users {
		if level > highestLevel {
			highestLevel = level
			count = 1
		} else if level == highestLevel {
			count++
		}
	}
	for userID, member := range data {
		level := levels.GetUserLevel(userID)
		sigil := ' '
		if level == highestLevel && count == 1 {
			sigil = '~'
		} else if level > levels.StateDefault() {
			sigil = '&'
		} else if level >= levels.Ban() {
			sigil = '@'
		} else if level >= levels.Kick() || level >= levels.Redact() {
			sigil = '%'
		} else if level > levels.UsersDefault {
			sigil = '+'
		}
		ml.memberList[i] = &memberListItem{
			Member:     *member,
			UserID:     userID,
			PowerLevel: level,
			Sigil:      sigil,
			Color:      widget.GetHashColor(userID),
		}
		i++
	}
	sort.Sort(ml.memberList)

	// Create text view for member list
	// TODO: figure out TextView dynamic colors
	var sb strings.Builder
	for _, member := range ml.memberList {

		// Sigil
		if member.Sigil != ' ' {
			sb.WriteRune(member.Sigil)
		} else {
			sb.WriteRune(' ')
		}

		sb.WriteRune(' ')

		// Display name
		if member.Membership == "invite" {
			sb.WriteRune('(')
		}
		sb.WriteString(member.Displayname)
		if member.Membership == "invite" {
			sb.WriteRune('(')
		}
		sb.WriteRune('\n')
	}
	ml.listView.SetText(sb.String())

	return ml
}

func (ml *MemberList) Draw(screen mauview.Screen) {
	ml.flex.Draw(screen)
}
