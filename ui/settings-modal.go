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

	"maunium.net/go/gomuks/config"
)

type SettingsModal struct {
	mauview.Component

	container *mauview.Box
	form *mauview.Form

	settingButton *mauview.Button

	parent *MainView
}

func NewSettingsModal(mainView *MainView, width int, height int) *SettingsModal {
	sm := &SettingsModal{
		form: mauview.NewForm(),
		parent: mainView,
	}

	sm.settingButton = mauview.NewButton("Setting on").
		SetOnClick(func() { sm.settingButton.SetText("Setting off") })

	sm.form.SetColumns([]int{1, 5, 1, 30, 1}).
		SetRows([]int{1, 1, 1, 1, 1, 1, 1, 1, 1})
	sm.form.AddFormItem(sm.settingButton, 3, 1, 1, 1)
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
