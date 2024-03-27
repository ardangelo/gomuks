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
	"strconv"

	"go.mau.fi/mauview"

	"maunium.net/go/gomuks/config"
)

type SettingsModal struct {
	mauview.Component

	container *mauview.Box
	form *mauview.Form

	roomListViewLabel *mauview.TextView
	roomListViewTagsButton *mauview.Button
	roomListViewHubButton *mauview.Button

	notifyLabel *mauview.TextView
	notifyOnButton *mauview.Button
	notifyOffButton *mauview.Button

	rewakeLabel *mauview.TextView
	rewakeEntry *mauview.InputField
	rewakeButton *mauview.Button

	parent *MainView
}

func formatCheckbox(format string, setting bool) string {
	if setting {
		return "(X) " + format
	} else {
		return "( ) " + format
	}
}

func NewSettingsModal(mainView *MainView, width int, height int) *SettingsModal {
	sm := &SettingsModal{
		form: mauview.NewForm(),
		parent: mainView,
	}

	sm.roomListViewLabel = mauview.NewTextView().SetText("Room list (restart req'd)")
	roomListViewTagsFormat := "Group rooms by tag"
	roomListViewHubFormat := "Sort by updated"
	sm.roomListViewTagsButton = mauview.NewButton(
		formatCheckbox(roomListViewTagsFormat, mainView.config.Preferences.TagGroupRooms)).
		SetOnClick(func() {
			sm.parent.config.Preferences.TagGroupRooms = true
			sm.roomListViewTagsButton.SetText(
				formatCheckbox(roomListViewTagsFormat,
					sm.parent.config.Preferences.TagGroupRooms))
			sm.roomListViewHubButton.SetText(
				formatCheckbox(roomListViewHubFormat,
					!sm.parent.config.Preferences.TagGroupRooms))
			sm.parent.gmx.Stop(true)
		})
	sm.roomListViewHubButton = mauview.NewButton(
		formatCheckbox(roomListViewHubFormat, !sm.parent.config.Preferences.TagGroupRooms)).
		SetOnClick(func() {
			sm.parent.config.Preferences.TagGroupRooms = false
			sm.roomListViewTagsButton.SetText(
				formatCheckbox(roomListViewTagsFormat,
					sm.parent.config.Preferences.TagGroupRooms))
			sm.roomListViewHubButton.SetText(
				formatCheckbox(roomListViewHubFormat,
					!sm.parent.config.Preferences.TagGroupRooms))
			sm.parent.gmx.Stop(true)
		})

	sm.notifyLabel = mauview.NewTextView().SetText("Notify on new messages")
	notifyOnFormat := "Flash LED until key pressed"
	notifyOffFormat := "Notification LED disabled"
	sm.notifyOnButton = mauview.NewButton(
		formatCheckbox(notifyOnFormat, sm.parent.config.NotifySound)).
		SetOnClick(func() {
			mainView.config.NotifySound = true
			sm.notifyOnButton.SetText(
				formatCheckbox(notifyOnFormat,
					sm.parent.config.NotifySound))
			sm.notifyOffButton.SetText(
				formatCheckbox(notifyOffFormat,
					!sm.parent.config.NotifySound))
		})
	sm.notifyOffButton = mauview.NewButton(
		formatCheckbox(notifyOffFormat, !sm.parent.config.NotifySound)).
		SetOnClick(func() {
			sm.parent.config.NotifySound = false
			sm.notifyOnButton.SetText(
				formatCheckbox(notifyOnFormat,
					sm.parent.config.NotifySound))
			sm.notifyOffButton.SetText(
				formatCheckbox(notifyOffFormat,
					!sm.parent.config.NotifySound))
		})

	sm.rewakeLabel = mauview.NewTextView().SetText("Rewake poll interval (minutes)")
	sm.rewakeEntry = mauview.NewInputField().
		SetTextAndMoveCursor(strconv.Itoa(
			sm.parent.config.Preferences.RewakeIntervalMins))
	sm.rewakeButton = mauview.NewButton("Apply").
		SetOnClick(func() {
			mins, err := strconv.Atoi(sm.rewakeEntry.GetText())
			if err != nil {
				sm.rewakeButton.SetText("Invalid")
			} else {
				sm.parent.config.Preferences.RewakeIntervalMins = mins
				sm.rewakeButton.SetText("Saved")
			}
		})

	sm.form.SetColumns([]int{1, width - 14, 1, 10, 1}).
		SetRows([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})

	sm.form.AddFormItem(sm.roomListViewLabel, 1, 1, 3, 1)
	sm.form.AddFormItem(sm.roomListViewTagsButton, 1, 2, 3, 1)
	sm.form.AddFormItem(sm.roomListViewHubButton, 1, 3, 3, 1)

	sm.form.AddFormItem(sm.notifyLabel, 1, 5, 3, 1)
	sm.form.AddFormItem(sm.notifyOnButton, 1, 6, 3, 1)
	sm.form.AddFormItem(sm.notifyOffButton, 1, 7, 3, 1)

	sm.form.AddFormItem(sm.rewakeLabel, 1, 9, 3, 1)
	sm.form.AddFormItem(sm.rewakeEntry, 1, 10, 1, 1)
	sm.form.AddFormItem(sm.rewakeButton, 3, 10, 1, 1)

	sm.form.FocusNextItem()

	sm.container = mauview.NewBox(sm.form).
		SetBorder(true).
		SetTitle("Settings").
		SetBlurCaptureFunc(func() bool {
			sm.parent.HideModal()
			return true
		})

	sm.Component = mauview.Center(sm.container, width, height).SetAlwaysFocusChild(true)

	return sm
}

func (fs *SettingsModal) Focus() {
	fs.container.Focus()
}

func (fs *SettingsModal) Blur() {
	fs.container.Blur()
}

func (sm *SettingsModal) OnKeyEvent(event mauview.KeyEvent) bool {
	kb := config.Keybind{
		Key: event.Key(),
		Ch:  event.Rune(),
		Mod: event.Modifiers(),
	}
	switch sm.parent.config.Keybindings.Modal[kb] {
	case "cancel":
		sm.parent.HideModal()
		return true
	}
	return sm.form.OnKeyEvent(event)
}
