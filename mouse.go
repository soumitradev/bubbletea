package tea

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const mouseX10ByteOffset = 32

// MouseMsg contains information about a mouse event and is sent to a program's
// update function when mouse activity occurs. Note that the mouse must first
// be enabled in order for the mouse events to be received.
type MouseMsg MouseEvent

// MouseEvent represents a mouse event, which could be a click, a scroll wheel
// movement, a cursor movement, or a combination.
type MouseEvent struct {
	X      int
	Y      int
	Shift  bool
	Alt    bool
	Ctrl   bool
	Action MouseAction
	Button MouseButton

	// Deprecated: Use MouseAction & MouseButton instead.
	Type MouseEventType

	isSGR bool
}

// IsWheel returns true if the mouse event is a wheel event.
func (m MouseEvent) IsWheel() bool {
	return m.Button == MouseButtonWheelUp || m.Button == MouseButtonWheelDown ||
		m.Button == MouseButtonWheelLeft || m.Button == MouseButtonWheelRight
}

// String returns a string representation of a mouse event.
func (m MouseEvent) String() (s string) {
	if m.Ctrl {
		s += "ctrl+"
	}
	if m.Alt {
		s += "alt+"
	}
	if m.Shift {
		s += "shift+"
	}

	if m.isSGR {
		if m.Button == MouseButtonNone && m.Action == MouseActionMotion {
			s += mouseActions[m.Action]
		} else if m.IsWheel() {
			s += mouseButtons[m.Button]
		} else {
			s += mouseButtons[m.Button]
			s += " "
			s += mouseActions[m.Action]
		}
	} else {
		s += mouseEventTypes[m.Type]
	}

	return s
}

// MouseAction represents the action that occurred during a mouse event.
type MouseAction int

// Mouse event actions.
const (
	MouseActionPress MouseAction = iota
	MouseActionRelease
	MouseActionMotion
)

var mouseActions = map[MouseAction]string{
	MouseActionPress:   "press",
	MouseActionRelease: "release",
	MouseActionMotion:  "motion",
}

// MouseButton represents the button that was pressed during a mouse event.
type MouseButton int

// Mouse event buttons
//
// This is based on X11 mouse button codes.
//
//	1 = left button
//	2 = middle button (pressing the scroll wheel)
//	3 = right button
//	4 = turn scroll wheel up
//	5 = turn scroll wheel down
//	6 = push scroll wheel left
//	7 = push scroll wheel right
//	8 = 4th button (aka browser backward button)
//	9 = 5th button (aka browser forward button)
//	10
//	11
//
// Other buttons are not supported.
const (
	MouseButtonNone MouseButton = iota
	MouseButtonLeft
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
	MouseButtonWheelLeft
	MouseButtonWheelRight
	MouseButtonBackward
	MouseButtonForward
	MouseButton10
	MouseButton11

	MouseButtonUnknown
)

var mouseButtons = map[MouseButton]string{
	MouseButtonNone:       "none",
	MouseButtonLeft:       "left",
	MouseButtonMiddle:     "middle",
	MouseButtonRight:      "right",
	MouseButtonWheelUp:    "wheel up",
	MouseButtonWheelDown:  "wheel down",
	MouseButtonWheelLeft:  "wheel left",
	MouseButtonWheelRight: "wheel right",
	MouseButtonBackward:   "backward",
	MouseButtonForward:    "forward",
	MouseButton10:         "button 10",
	MouseButton11:         "button 11",
	MouseButtonUnknown:    "unknown",
}

// MouseEventType indicates the type of mouse event occurring.
//
// Deprecated: Use MouseAction & MouseButton instead.
type MouseEventType int

// Mouse event types.
//
// Deprecated: Use MouseAction & MouseButton instead.
const (
	MouseUnknown MouseEventType = iota
	MouseLeft
	MouseRight
	MouseMiddle
	MouseRelease // mouse button release (X10 only)
	MouseWheelUp
	MouseWheelDown
	MouseWheelLeft
	MouseWheelRight
	MouseBackward
	MouseForward
	MouseMotion
)

var mouseEventTypes = map[MouseEventType]string{
	MouseUnknown:    "unknown",
	MouseLeft:       "left",
	MouseRight:      "right",
	MouseMiddle:     "middle",
	MouseRelease:    "release",
	MouseWheelUp:    "wheel up",
	MouseWheelDown:  "wheel down",
	MouseWheelLeft:  "wheel left",
	MouseWheelRight: "wheel right",
	MouseBackward:   "backward",
	MouseForward:    "forward",
	MouseMotion:     "motion",
}

var (
	mouseX10Seq = []byte("\x1b[M")
	mouseSGRSeq = []byte("\x1b[<")
)

func parseMouseEvents(buf []byte) ([]MouseEvent, error) {
	if len(buf) == 0 {
		return nil, errors.New("empty buffer")
	}

	switch {
	case bytes.Contains(buf, mouseSGRSeq):
		return parseSGRMouseEvents(string(buf))
	case bytes.Contains(buf, mouseX10Seq):
		return parseX10MouseEvents(buf)
	}

	return nil, errors.New("not a mouse event")
}

var mouseSGRRegex = regexp.MustCompile(`(\d+);(\d+);(\d+)([Mm])`)

// parseSGRMouseEvents parses SGR extended mouse events. SGR mouse events look
// like:
//
//	ESC [ < Cb ; Cx ; Cy (M or m)
//
// where:
//
//	Cb is the encoded button code
//	Cx is the x-coordinate of the mouse
//	Cy is the y-coordinate of the mouse
//	M is for button press, m is for button release
//
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Extended-coordinates
func parseSGRMouseEvents(buf string) ([]MouseEvent, error) {
	var ev []MouseEvent

	seq := string(mouseSGRSeq)
	if !strings.Contains(buf, seq) {
		return nil, errors.New("not a SGR mouse event")
	}

	for _, v := range strings.Split(buf, seq) {
		if len(v) == 0 {
			continue
		}

		matches := mouseSGRRegex.FindStringSubmatch(v)
		if len(matches) != 5 {
			return nil, errors.New("not a SGR mouse event")
		}

		b, _ := strconv.Atoi(matches[1])
		px := matches[2]
		py := matches[3]
		release := matches[4] == "m"
		m := parseMouseButton(b, true)
		// Wheel buttons don't have  release events
		// Motion can be reported as a release event in some terminals (Windows Terminal)
		if m.Action != MouseActionMotion && !m.IsWheel() && release {
			m.Action = MouseActionRelease
			m.Type = MouseRelease
		}

		x, _ := strconv.Atoi(px)
		y, _ := strconv.Atoi(py)

		// (1,1) is the upper left. We subtract 1 to normalize it to (0,0).
		m.X = x - 1
		m.Y = y - 1

		ev = append(ev, m)
	}

	return ev, nil
}

// Parse X10-encoded mouse events; the simplest kind. The last release of X10
// was December 1986, by the way. The original X10 mouse protocol limits the Cx
// and Cy ordinates to 223 (=255 - 32).
//
// X10 mouse events look like:
//
//	ESC [M Cb Cx Cy
//
// See: http://www.xfree86.org/current/ctlseqs.html#Mouse%20Tracking
func parseX10MouseEvents(buf []byte) ([]MouseEvent, error) {
	var r []MouseEvent

	seq := mouseX10Seq
	if !bytes.Contains(buf, seq) {
		return r, errors.New("not an X10 mouse event")
	}

	for _, v := range bytes.Split(buf, seq) {
		if len(v) == 0 {
			continue
		}
		if len(v) != 3 {
			return r, errors.New("not an X10 mouse event")
		}

		m := parseMouseButton(int(v[0]), false)

		// (1,1) is the upper left. We subtract 1 to normalize it to (0,0).
		m.X = int(v[1]) - mouseX10ByteOffset - 1
		m.Y = int(v[2]) - mouseX10ByteOffset - 1

		r = append(r, m)
	}

	return r, nil
}

// See: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Extended-coordinates
func parseMouseButton(b int, isSGR bool) MouseEvent {
	var m MouseEvent
	m.isSGR = isSGR
	e := b
	if !isSGR {
		e -= mouseX10ByteOffset
	}

	const (
		bitShift  = 0b0000_0100
		bitAlt    = 0b0000_1000
		bitCtrl   = 0b0001_0000
		bitMotion = 0b0010_0000
		bitWheel  = 0b0100_0000
		bitAdd    = 0b1000_0000 // additional buttons 8-11

		bitsMask = 0b0000_0011
	)

	if e&bitAdd != 0 {
		m.Button = MouseButtonBackward + MouseButton(e&bitsMask)
	} else if e&bitWheel != 0 {
		m.Button = MouseButtonWheelUp + MouseButton(e&bitsMask)
	} else {
		m.Button = MouseButtonLeft + MouseButton(e&bitsMask)
		// X10 reports a button release as 0b0000_0011 (3)
		if e&bitsMask == bitsMask {
			m.Action = MouseActionRelease
			m.Button = MouseButtonNone
		}
	}
	// Motion bit doesn't get reported for wheel events but we check for it
	// anyway in case of a faulty terminal emulator.
	if e&bitMotion != 0 && !m.IsWheel() {
		m.Action = MouseActionMotion
	}

	// backward compatibility
	switch {
	case m.Button == MouseButtonLeft && m.Action == MouseActionPress:
		m.Type = MouseLeft
	case m.Button == MouseButtonMiddle && m.Action == MouseActionPress:
		m.Type = MouseMiddle
	case m.Button == MouseButtonRight && m.Action == MouseActionPress:
		m.Type = MouseRight
	case m.Button == MouseButtonNone && m.Action == MouseActionRelease:
		m.Type = MouseRelease
	case m.Button == MouseButtonWheelUp && m.Action == MouseActionPress:
		m.Type = MouseWheelUp
	case m.Button == MouseButtonWheelDown && m.Action == MouseActionPress:
		m.Type = MouseWheelDown
	case m.Button == MouseButtonWheelLeft && m.Action == MouseActionPress:
		m.Type = MouseWheelLeft
	case m.Button == MouseButtonWheelRight && m.Action == MouseActionPress:
		m.Type = MouseWheelRight
	case m.Button == MouseButtonBackward && m.Action == MouseActionPress:
		m.Type = MouseBackward
	case m.Button == MouseButtonForward && m.Action == MouseActionPress:
		m.Type = MouseForward
	case m.Action == MouseActionMotion:
		m.Type = MouseMotion
		switch m.Button {
		case MouseButtonLeft:
			m.Type = MouseLeft
		case MouseButtonMiddle:
			m.Type = MouseMiddle
		case MouseButtonRight:
			m.Type = MouseRight
		case MouseButtonBackward:
			m.Type = MouseBackward
		case MouseButtonForward:
			m.Type = MouseForward
		}
	default:
		m.Type = MouseUnknown
	}

	m.Alt = e&bitAlt != 0
	m.Ctrl = e&bitCtrl != 0
	m.Shift = e&bitShift != 0

	return m
}
