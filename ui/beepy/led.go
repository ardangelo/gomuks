package beepy

import (
	"os"
	"strconv"
)

type LED struct {
	level byte
	r byte
	g byte
	b byte
}

func NewLED() (*LED, error) {
	return &LED{0, 0, 0, 0}, nil
}

func (l *LED) Close() error {
	return nil
}

func (l *LED) On() error {
	l.level = 0xff
	return os.WriteFile("/sys/firmware/beepy/led", []byte(strconv.Itoa(int(l.level)) + "\n"), 0220)
}

func (l *LED) Off() error {
	l.level = 0x0
	return os.WriteFile("/sys/firmware/beepy/led", []byte(strconv.Itoa(int(l.level)) + "\n"), 0220)
}

func (l *LED) SetColor(r, g, b uint16) error {

	l.r = byte(r)
	l.g = byte(g)
	l.b = byte(b)

	if err := os.WriteFile("/sys/firmware/beepy/led_red", []byte(strconv.Itoa(int(l.r)) + "\n"), 0220); err != nil {
		return err
	}

	if err := os.WriteFile("/sys/firmware/beepy/led_green", []byte(strconv.Itoa(int(l.g)) + "\n"), 0220); err != nil {
		return err
	}

	if err := os.WriteFile("/sys/firmware/beepy/led_blue", []byte(strconv.Itoa(int(l.b)) + "\n"), 0220); err != nil {
		return err
	}

	return nil
}

func (l *LED) IsOn() (bool, error) {
	return l.level > 0, nil
}

func (l *LED) Color() ([]byte, error) {
	return []byte{l.r, l.g, l.b}, nil
}
