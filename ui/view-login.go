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

	"maunium.net/go/mautrix"

	"maunium.net/go/gomuks/beeper"
	"maunium.net/go/gomuks/config"
	"maunium.net/go/gomuks/debug"
	ifc "maunium.net/go/gomuks/interface"
)

type LoginView struct {
	*mauview.Form

	container *mauview.Centerer
	box *mauview.Box

	emailLabel *mauview.TextField
	email *mauview.InputField
	requestCodeButton *mauview.Button

	codeLabel  *mauview.TextField
	code  *mauview.InputField

	error *mauview.TextView

	loginButton *mauview.Button
	quitButton  *mauview.Button

	session string

	matrix ifc.MatrixContainer
	config *config.Config
	parent *GomuksUI
}

func (ui *GomuksUI) NewLoginView() mauview.Component {
	view := &LoginView{
		Form: mauview.NewForm(),

		emailLabel: mauview.NewTextField().SetText("Email"),
		email: mauview.NewInputField(),
		requestCodeButton: mauview.NewButton("Enter email to request code"),

		codeLabel:  mauview.NewTextField().SetText("Code"),
		code:  mauview.NewInputField(),

		loginButton: mauview.NewButton("Enter code to log in"),
		quitButton:  mauview.NewButton("Quit"),

		matrix: ui.gmx.Matrix(),
		config: ui.gmx.Config(),
		parent: ui,
	}

	view.email.SetPlaceholder("example@example.com").
		SetTextColor(tcell.ColorWhite).
		SetChangedFunc(view.emailChanged)
	view.code.SetPlaceholder("123456").
		SetTextColor(tcell.ColorWhite).
		SetChangedFunc(view.codeChanged)

	view.requestCodeButton.
		SetBackgroundColor(tcell.ColorDarkCyan).
		SetForegroundColor(tcell.ColorWhite).
		SetFocusedForegroundColor(tcell.ColorWhite)
	view.loginButton.
		SetBackgroundColor(tcell.ColorDarkCyan).
		SetForegroundColor(tcell.ColorWhite).
		SetFocusedForegroundColor(tcell.ColorWhite)
	view.quitButton.
		SetOnClick(func() { ui.gmx.Stop(true) }).
		SetBackgroundColor(tcell.ColorDarkCyan).
		SetForegroundColor(tcell.ColorWhite).
		SetFocusedForegroundColor(tcell.ColorWhite)

	view.
		SetColumns([]int{1, 5, 1, 30, 1}).
		SetRows([]int{1, 1, 1, 1, 1, 1, 1, 1, 1})
	view.
		AddFormItem(view.email, 3, 1, 1, 1).
		AddFormItem(view.requestCodeButton, 3, 3, 1, 1).
		AddFormItem(view.code, 3, 5, 1, 1).
		AddFormItem(view.loginButton, 1, 7, 3, 1).
		AddFormItem(view.quitButton, 1, 9, 3, 1).
		AddComponent(view.emailLabel, 1, 1, 1, 1).
		AddComponent(view.codeLabel, 1, 5, 1, 1)
	view.FocusNextItem()
	ui.loginView = view

	view.box = mauview.NewBox(view).SetTitle("Log in to Matrix")
	view.container = mauview.Center(view.box, 40, 12)
	view.container.SetAlwaysFocusChild(true)
	return view.container
}

func (view *LoginView) emailAuthFlow() {

	resp, err := beeper.StartLogin()
	if err != nil {
		view.code.SetText("")
		view.Error(err.Error())
		return
	}

	view.requestCodeButton.SetText("Sending request")
	view.parent.Render()

	err = beeper.SendLoginEmail(resp.RequestID, view.email.GetText())
	if err != nil {
		view.code.SetText("")
		view.code.SetPlaceholder("123456")
		view.Error(err.Error())
		return
	}
	view.session = resp.RequestID

	view.requestCodeButton.SetText("Check email for code")
	view.parent.Render()
}

func (view *LoginView) emailChanged(to string) {
	var emailFilled = len(view.email.GetText()) > 0

	if emailFilled {
		view.requestCodeButton.SetText("Request code").
			SetOnClick(view.emailAuthFlow)
	} else {
		view.requestCodeButton.SetText("Enter email to request code").
			SetOnClick(nil)
	}
}

func (view *LoginView) codeChanged(to string) {
	var codeFilled = len(to) > 0

	if codeFilled {
		view.loginButton.SetText("Login").
			SetOnClick(view.Login)
	} else {
		view.loginButton.SetText("Enter code to log in").
			SetOnClick(nil)
	}
}

func (view *LoginView) Error(err string) {
	if len(err) > 0 {
		view.box.SetTitle(err)
	} else {
		view.box.SetTitle("Log in to Matrix")
	}

	view.parent.Render()
}

func (view *LoginView) actuallyLogin(session, code string) {
	debug.Printf("Logging into Beeper with code %s...", code)
	view.config.HS = "https://matrix.beeper.com"

	if err := view.matrix.InitClient(false); err != nil {
		debug.Print("Init error:", err)
		view.Error(err.Error())

	} else if err = view.matrix.BeeperLogin(session, code); err != nil {
		debug.Print("Login error:", err)

		if httpErr, ok := err.(mautrix.HTTPError); ok {

			// Known status codes
			if httpErr.IsStatus(400) {
				view.Error("Beeper: invalid email")

			} else if httpErr.IsStatus(404) {
				view.Error("Beeper: unknown email")

			} else if httpErr.IsStatus(429) {
				view.Error("Beeper: too many req's")

			// Other status codes
			} else if httpErr.RespError == nil {
				view.Error(httpErr.RespError.Err)

			// General errors
			} else if len(httpErr.Message) > 0 {
				view.Error(httpErr.Message)
			} else {
				view.Error(err.Error())
			}

		} else {
			view.Error(err.Error())
		}
	}

	view.loginButton.SetText("Login").
		SetOnClick(view.Login)
}

func (view *LoginView) Login() {
	code := view.code.GetText()

	view.loginButton.SetText("Logging in...").
		SetOnClick(nil)
	go view.actuallyLogin(view.session, code)
}
