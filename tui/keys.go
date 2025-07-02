package tui

import "fmt"

type Key byte

const (
	Key0 Key = 48
	Key1 Key = 49
	Key2 Key = 50
	Key3 Key = 51
	Key4 Key = 52
	Key5 Key = 53
	Key6 Key = 54
	Key7 Key = 55
	Key8 Key = 56
	Key9 Key = 57

	KeyA Key = 97
	KeyB Key = 98
	KeyC Key = 99
	KeyD Key = 100
	KeyE Key = 101
	KeyF Key = 102
	KeyG Key = 103
	KeyH Key = 104
	KeyI Key = 105
	KeyJ Key = 106
	KeyK Key = 107
	KeyL Key = 108
	KeyM Key = 109
	KeyN Key = 110
	KeyO Key = 111
	KeyP Key = 112
	KeyQ Key = 113
	KeyR Key = 114
	KeyS Key = 115
	KeyT Key = 116
	KeyU Key = 117
	KeyV Key = 118
	KeyW Key = 119
	KeyX Key = 120
	KeyY Key = 121
	KeyZ Key = 122

	KeyEnter     Key = 13
	KeyBackspace Key = 127
	KeyHash      Key = 35
	KeyAmpersand Key = 38

	KeyCtrlN Key = 14
	KeyCtrlP Key = 16
	KeyCtrlB Key = 2
	KeyCtrlF Key = 6
	KeyCtrlU Key = 21
	KeyCtrlC Key = 3
	KeyCtrlD Key = 4
)

func (k Key) String() string {
	switch k {
	case Key0:
		return "Key0"
	case Key1:
		return "Key1"
	case Key2:
		return "Key2"
	case Key3:
		return "Key3"
	case Key4:
		return "Key4"
	case Key5:
		return "Key5"
	case Key6:
		return "Key6"
	case Key7:
		return "Key7"
	case Key8:
		return "Key8"
	case Key9:
		return "Key9"

	case KeyA:
		return "KeyA"
	case KeyB:
		return "KeyB"
	case KeyC:
		return "KeyC"
	case KeyD:
		return "KeyD"
	case KeyE:
		return "KeyE"
	case KeyF:
		return "KeyF"
	case KeyG:
		return "KeyG"
	case KeyH:
		return "KeyH"
	case KeyI:
		return "KeyI"
	case KeyJ:
		return "KeyJ"
	case KeyK:
		return "KeyK"
	case KeyL:
		return "KeyL"
	case KeyM:
		return "KeyM"
	case KeyN:
		return "KeyN"
	case KeyO:
		return "KeyO"
	case KeyP:
		return "KeyP"
	case KeyQ:
		return "KeyQ"
	case KeyR:
		return "KeyR"
	case KeyS:
		return "KeyS"
	case KeyT:
		return "KeyT"
	case KeyU:
		return "KeyU"
	case KeyV:
		return "KeyV"
	case KeyW:
		return "KeyW"
	case KeyX:
		return "KeyX"
	case KeyY:
		return "KeyY"
	case KeyZ:
		return "KeyZ"

	case KeyEnter:
		return "KeyEnter"
	case KeyBackspace:
		return "KeyBackspace"
	case KeyHash:
		return "KeyHash"
	case KeyAmpersand:
		return "KeyAmpersand"

	case KeyCtrlN:
		return "KeyCtrlN"
	case KeyCtrlP:
		return "KeyCtrlP"
	case KeyCtrlB:
		return "KeyCtrlB"
	case KeyCtrlF:
		return "KeyCtrlF"
	case KeyCtrlU:
		return "KeyCtrlU"
	case KeyCtrlC:
		return "KeyCtrlC"
	case KeyCtrlD:
		return "KeyCtrlD"

	default:
		return fmt.Sprintf("Key(%d)", k)
	}
}
