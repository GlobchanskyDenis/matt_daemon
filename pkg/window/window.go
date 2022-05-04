package window

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2" // Это у нас fyne -- руки у них из жопы
	"time"
	"strconv"
)

/*	Тут описано создание клиентского приложения и его поведение. Работа по подключению к сокету, общение с ним
**	передается через аргументы (функции)  */

var _ Window = (*window)(nil)

type Window interface {
	ShowAndRun()
}

type window struct {
	mainWindow fyne.Window
}

func New(isConnectServerCorrect func(bool, string, uint) bool, checkAuth func(string, string) bool, sendMessage func(string)) Window {
	a := app.New()

	mainWindow := a.NewWindow("Подключение")
	mainWindow.Resize(fyne.NewSize(300, 200))

	oldMessagesLabel1 := widget.NewLabel("")
	oldMessagesLabel2 := widget.NewLabel("")
	oldMessagesLabel3 := widget.NewLabel("")
	oldMessagesLabel4 := widget.NewLabel("")
	sendingEntry := widget.NewEntry()
	sendButton := widget.NewButton("Submit", func() {
		sendMessage(sendingEntry.Text)
		oldMessagesLabel1.SetText(oldMessagesLabel2.Text)
		oldMessagesLabel2.SetText(oldMessagesLabel3.Text)
		oldMessagesLabel3.SetText(oldMessagesLabel4.Text)
		oldMessagesLabel4.SetText(sendingEntry.Text)
		sendingEntry.SetText("")
	})

	sendVBox := container.NewVBox(
		oldMessagesLabel1,
		oldMessagesLabel2,
		oldMessagesLabel3,
		oldMessagesLabel4,
		sendingEntry,
		sendButton,
	)
	sendVBox.Hide()

	successLabel := widget.NewLabel("Success")

	loginLabel := widget.NewLabel("")
	loginEntry := widget.NewEntry()
	loginEntry.SetText("Login")
	passwordEntry := widget.NewPasswordEntry()
	authVBox := container.NewVBox(
		loginLabel,
		loginEntry,
		passwordEntry,
		widget.NewButton("Submit", func() {
			if checkAuth(loginEntry.Text, passwordEntry.Text) == true {
				mainWindow.SetContent(successLabel)
				time.Sleep(1500 * time.Millisecond)
				mainWindow.SetContent(sendVBox)
				mainWindow.SetTitle("Работаем")
				sendVBox.Show()
			} else {
				loginLabel.SetText("Ошибка")
				time.Sleep(1000 * time.Millisecond)
				loginLabel.SetText("")
			}
			
		}),
	)
	authVBox.Hide()

	connLabel := widget.NewLabel("Choose connection type to continue")
	connTypeRadio := widget.NewRadioGroup([]string{"tcp", "udp"}, nil)
	connTypeRadio.SetSelected("tcp")
	ipEntry := widget.NewEntry()
	ipEntry.SetText("Ip")
	portEntry := widget.NewEntry()
	portEntry.SetText("Port")
	
	connVBox := container.NewVBox(
		connLabel,
		connTypeRadio,
		ipEntry,
		portEntry,
		widget.NewButton("Submit", func() {
			port, err := strconv.ParseUint(portEntry.Text, 10, 64)
			if err != nil {
				connLabel.SetText("Ошибка ввода")
				time.Sleep(1000 * time.Millisecond)
				connLabel.SetText("Choose connection type to continue")
				return
			}
			var isTcpConn bool
			if connTypeRadio.Selected == "tcp" {
				isTcpConn = true
			} else {
				isTcpConn = false
			}
			if isConnectServerCorrect(isTcpConn, ipEntry.Text, uint(port)) == true {
				mainWindow.SetContent(successLabel)
				time.Sleep(1500 * time.Millisecond)
				mainWindow.SetContent(authVBox)
				mainWindow.SetTitle("Авторизация")
				authVBox.Show()
			} else {
				connLabel.SetText("Ошибка")
				time.Sleep(1000 * time.Millisecond)
				connLabel.SetText("Choose connection type to continue")
			}
		}),
	)

	mainWindow.SetContent(connVBox)
	
	return &window{
		mainWindow: mainWindow,
	}
}

func (w *window) ShowAndRun() {
	w.mainWindow.ShowAndRun()
}

