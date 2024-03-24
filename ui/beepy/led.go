package beepy

import (
	"fmt"
	"os"
	"sync"
)

const ledPath = "/sys/firmware/beepy/led"
const ledRedPath = "/sys/firmware/beepy/led_red"
const ledBluePath = "/sys/firmware/beepy/led_blue"
const ledGreenPath = "/sys/firmware/beepy/led_green"

const (
	ledSetOff = 0x00
	ledSetOn = 0x01
	ledSetFlash = 0x02
	ledSetFlashUntilKey = 0x03
)

type LED struct {
	lck sync.RWMutex
}

func NewLED() (*LED, error) {

	// Ensure LED control path writable
	file, err := os.OpenFile(ledPath, os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return &LED{}, nil
}

func writeInt(path string, value uint16) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%d\n", value)
	return nil
}

func (l *LED) setNextColor(r, g, b uint16) error {
	if err := writeInt(ledRedPath, r); err != nil {
		return err
	}

	if err := writeInt(ledGreenPath, g); err != nil {
		return err
	}

	if err := writeInt(ledBluePath, b); err != nil {
		return err
	}

	return nil
}

func (l *LED) On(r, g, b uint16) error {

	l.lck.Lock()
	defer l.lck.Unlock()

	if err := l.setNextColor(r, g, b); err != nil {
		return err
	}

	return writeInt(ledPath, ledSetOn)
}

func (l *LED) FlashUntilKey(r, g, b uint16) error {

	l.lck.Lock()
	defer l.lck.Unlock()

	if err := l.setNextColor(r, g, b); err != nil {
		return err
	}

	return writeInt(ledPath, ledSetFlashUntilKey)
}

func (l *LED) Off() error {
	return writeInt(ledPath, ledSetOff)
}
