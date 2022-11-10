package tea

import (
	"fmt"
	"testing"
)

func TestMouseEvent_String(t *testing.T) {
	tt := []struct {
		name     string
		event    MouseEvent
		expected string
	}{
		{
			name: "unknown",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseUnknown,
			},
			expected: "unknown",
		},
		{
			name: "left",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseLeft,
			},
			expected: "left",
		},
		{
			name: "right",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseRight,
			},
			expected: "right",
		},
		{
			name: "middle",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseMiddle,
			},
			expected: "middle",
		},
		{
			name: "release",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseRelease,
			},
			expected: "release",
		},
		{
			name: "wheel up",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseWheelUp,
			},
			expected: "wheel up",
		},
		{
			name: "wheel down",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseWheelDown,
			},
			expected: "wheel down",
		},
		{
			name: "wheel left",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseWheelLeft,
			},
			expected: "wheel left",
		},
		{
			name: "wheel right",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseWheelRight,
			},
			expected: "wheel right",
		},
		{
			name: "motion",
			event: MouseEvent{
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Type:   MouseMotion,
			},
			expected: "motion",
		},
		{
			name: "shift+left",
			event: MouseEvent{
				Type:   MouseLeft,
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Shift:  true,
			},
			expected: "shift+left",
		},
		{
			name: "ctrl+shift+left",
			event: MouseEvent{
				Type:   MouseLeft,
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Shift:  true,
				Ctrl:   true,
			},
			expected: "ctrl+shift+left",
		},
		{
			name: "alt+left",
			event: MouseEvent{
				Type:   MouseLeft,
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Alt:    true,
			},
			expected: "alt+left",
		},
		{
			name: "ctrl+left",
			event: MouseEvent{
				Type:   MouseLeft,
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Ctrl:   true,
			},
			expected: "ctrl+left",
		},
		{
			name: "ctrl+alt+left",
			event: MouseEvent{
				Type:   MouseLeft,
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Alt:    true,
				Ctrl:   true,
			},
			expected: "ctrl+alt+left",
		},
		{
			name: "ctrl+alt+shift+left",
			event: MouseEvent{
				Type:   MouseLeft,
				Action: MouseActionPress,
				Button: MouseButtonNone,
				Alt:    true,
				Ctrl:   true,
				Shift:  true,
			},
			expected: "ctrl+alt+shift+left",
		},
		{
			name: "ignore coordinates",
			event: MouseEvent{
				X:      100,
				Y:      200,
				Type:   MouseLeft,
				Action: MouseActionPress,
				Button: MouseButtonNone,
			},
			expected: "left",
		},
		{
			name: "broken type",
			event: MouseEvent{
				Type:   MouseEventType(-100),
				Action: MouseActionPress,
				Button: MouseButtonNone,
			},
			expected: "",
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			actual := tc.event.String()

			if tc.expected != actual {
				t.Fatalf("expected %q but got %q",
					tc.expected,
					actual,
				)
			}
		})
	}
}

func TestParseX10MouseEvent(t *testing.T) {
	encode := func(b byte, x, y int) []byte {
		return []byte{
			'\x1b',
			'[',
			'M',
			byte(32) + b,
			byte(x + 32 + 1),
			byte(y + 32 + 1),
		}
	}

	tt := []struct {
		name     string
		buf      []byte
		expected []MouseEvent
	}{
		// Position.
		{
			name: "zero position",
			buf:  encode(0b0000_0000, 0, 0),
			expected: []MouseEvent{
				{
					X:      0,
					Y:      0,
					Type:   MouseLeft,
					Action: MouseActionPress,
					Button: MouseButtonLeft,
				},
			},
		},
		{
			name: "max position",
			buf:  encode(0b0000_0000, 222, 222), // Because 255 (max int8) - 32 - 1.
			expected: []MouseEvent{
				{
					X:      222,
					Y:      222,
					Type:   MouseLeft,
					Action: MouseActionPress,
					Button: MouseButtonLeft,
				},
			},
		},
		// Simple.
		{
			name: "left",
			buf:  encode(0b0000_0000, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseLeft,
					Action: MouseActionPress,
					Button: MouseButtonLeft,
				},
			},
		},
		{
			name: "left in motion",
			buf:  encode(0b0010_0000, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseLeft,
					Action: MouseActionMotion,
					Button: MouseButtonLeft,
				},
			},
		},
		{
			name: "middle",
			buf:  encode(0b0000_0001, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseMiddle,
					Action: MouseActionPress,
					Button: MouseButtonMiddle,
				},
			},
		},
		{
			name: "middle in motion",
			buf:  encode(0b0010_0001, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseMiddle,
					Action: MouseActionMotion,
					Button: MouseButtonMiddle,
				},
			},
		},
		{
			name: "right",
			buf:  encode(0b0000_0010, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseRight,
					Action: MouseActionPress,
					Button: MouseButtonRight,
				},
			},
		},
		{
			name: "right in motion",
			buf:  encode(0b0010_0010, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseRight,
					Action: MouseActionMotion,
					Button: MouseButtonRight,
				},
			},
		},
		{
			name: "motion",
			buf:  encode(0b0010_0011, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseMotion,
					Action: MouseActionMotion,
					Button: MouseButtonNone,
				},
			},
		},
		{
			name: "wheel up",
			buf:  encode(0b0100_0000, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseWheelUp,
					Action: MouseActionPress,
					Button: MouseButtonWheelUp,
				},
			},
		},
		{
			name: "wheel down",
			buf:  encode(0b0100_0001, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseWheelDown,
					Action: MouseActionPress,
					Button: MouseButtonWheelDown,
				},
			},
		},
		{
			name: "wheel left",
			buf:  encode(0b0100_0010, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseWheelLeft,
					Action: MouseActionPress,
					Button: MouseButtonWheelLeft,
				},
			},
		},
		{
			name: "wheel right",
			buf:  encode(0b0100_0011, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseWheelRight,
					Action: MouseActionPress,
					Button: MouseButtonWheelRight,
				},
			},
		},
		{
			name: "release",
			buf:  encode(0b0000_0011, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseRelease,
					Action: MouseActionRelease,
					Button: MouseButtonNone,
				},
			},
		},
		{
			name: "backward",
			buf:  encode(0b1000_0000, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseBackward,
					Action: MouseActionPress,
					Button: MouseButtonBackward,
				},
			},
		},
		{
			name: "forward",
			buf:  encode(0b1000_0001, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseForward,
					Action: MouseActionPress,
					Button: MouseButtonForward,
				},
			},
		},
		{
			name: "button 10",
			buf:  encode(0b1000_0010, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseUnknown,
					Action: MouseActionPress,
					Button: MouseButton10,
				},
			},
		},
		{
			name: "button 11",
			buf:  encode(0b1000_0011, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseUnknown,
					Action: MouseActionPress,
					Button: MouseButton11,
				},
			},
		},
		// Combinations.
		{
			name: "alt+right",
			buf:  encode(0b0000_1010, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Alt:    true,
					Type:   MouseRight,
					Action: MouseActionPress,
					Button: MouseButtonRight,
				},
			},
		},
		{
			name: "ctrl+right",
			buf:  encode(0b0001_0010, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Ctrl:   true,
					Type:   MouseRight,
					Action: MouseActionPress,
					Button: MouseButtonRight,
				},
			},
		},
		{
			name: "alt+right in motion",
			buf:  encode(0b0010_1010, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Alt:    true,
					Type:   MouseRight,
					Action: MouseActionMotion,
					Button: MouseButtonRight,
				},
			},
		},
		{
			name: "ctrl+right in motion",
			buf:  encode(0b0011_0010, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Ctrl:   true,
					Type:   MouseRight,
					Action: MouseActionMotion,
					Button: MouseButtonRight,
				},
			},
		},
		{
			name: "ctrl+alt+right",
			buf:  encode(0b0001_1010, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Alt:    true,
					Ctrl:   true,
					Type:   MouseRight,
					Action: MouseActionPress,
					Button: MouseButtonRight,
				},
			},
		},
		{
			name: "ctrl+wheel up",
			buf:  encode(0b0101_0000, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Ctrl:   true,
					Type:   MouseWheelUp,
					Action: MouseActionPress,
					Button: MouseButtonWheelUp,
				},
			},
		},
		{
			name: "alt+wheel down",
			buf:  encode(0b0100_1001, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Alt:    true,
					Type:   MouseWheelDown,
					Action: MouseActionPress,
					Button: MouseButtonWheelDown,
				},
			},
		},
		{
			name: "ctrl+alt+wheel down",
			buf:  encode(0b0101_1001, 32, 16),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Alt:    true,
					Ctrl:   true,
					Type:   MouseWheelDown,
					Action: MouseActionPress,
					Button: MouseButtonWheelDown,
				},
			},
		},
		// Overflow position.
		{
			name: "overflow position",
			buf:  encode(0b0010_0000, 250, 223), // Because 255 (max int8) - 32 - 1.
			expected: []MouseEvent{
				{
					X:      -6,
					Y:      -33,
					Type:   MouseLeft,
					Action: MouseActionMotion,
					Button: MouseButtonLeft,
				},
			},
		},
		// Batched events.
		{
			name: "batched events",
			buf:  append(encode(0b0010_0000, 32, 16), encode(0b0000_0011, 64, 32)...),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseLeft,
					Action: MouseActionMotion,
					Button: MouseButtonLeft,
				},
				{
					X:      64,
					Y:      32,
					Type:   MouseRelease,
					Action: MouseActionRelease,
					Button: MouseButtonNone,
				},
			},
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			actual, err := parseX10MouseEvents(tc.buf)
			if err != nil {
				t.Fatalf("unexpected error for test: %v",
					err,
				)
			}

			for i := range tc.expected {
				if tc.expected[i] != actual[i] {
					t.Fatalf("expected %#v but got %#v",
						tc.expected[i],
						actual[i],
					)
				}
			}
		})
	}
}

func TestParseX10MouseEvent_error(t *testing.T) {
	tt := []struct {
		name string
		buf  []byte
	}{
		{
			name: "empty buf",
			buf:  nil,
		},
		{
			name: "wrong high bit",
			buf:  []byte("\x1a[M@A1"),
		},
		{
			name: "short buf",
			buf:  []byte("\x1b[M@A"),
		},
		{
			name: "long buf",
			buf:  []byte("\x1b[M@A11"),
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			_, err := parseX10MouseEvents(tc.buf)

			if err == nil {
				t.Fatalf("expected error but got nil")
			}
		})
	}
}

func TestParseSGRMouseEvent(t *testing.T) {
	encode := func(b, x, y int, r bool) string {
		re := 'M'
		if r {
			re = 'm'
		}
		return fmt.Sprintf("\x1b[<%d;%d;%d%c", b, x+1, y+1, re)
	}

	tt := []struct {
		name     string
		buf      string
		expected []MouseEvent
	}{
		// Position.
		{
			name: "zero position",
			buf:  encode(0, 0, 0, false),
			expected: []MouseEvent{
				{
					X:      0,
					Y:      0,
					Type:   MouseLeft,
					Action: MouseActionPress,
					Button: MouseButtonLeft,
					isSGR:  true,
				},
			},
		},
		{
			name: "225 position",
			buf:  encode(0, 225, 225, false),
			expected: []MouseEvent{
				{
					X:      225,
					Y:      225,
					Type:   MouseLeft,
					Action: MouseActionPress,
					Button: MouseButtonLeft,
					isSGR:  true,
				},
			},
		},
		// Simple.
		{
			name: "left",
			buf:  encode(0, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseLeft,
					Action: MouseActionPress,
					Button: MouseButtonLeft,
					isSGR:  true,
				},
			},
		},
		{
			name: "left in motion",
			buf:  encode(32, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseLeft,
					Action: MouseActionMotion,
					Button: MouseButtonLeft,
					isSGR:  true,
				},
			},
		},
		{
			name: "left release",
			buf:  encode(0, 32, 16, true),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseRelease,
					Action: MouseActionRelease,
					Button: MouseButtonLeft,
					isSGR:  true,
				},
			},
		},
		{
			name: "middle",
			buf:  encode(1, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseMiddle,
					Action: MouseActionPress,
					Button: MouseButtonMiddle,
					isSGR:  true,
				},
			},
		},
		{
			name: "middle in motion",
			buf:  encode(33, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseMiddle,
					Action: MouseActionMotion,
					Button: MouseButtonMiddle,
					isSGR:  true,
				},
			},
		},
		{
			name: "middle release",
			buf:  encode(1, 32, 16, true),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseRelease,
					Action: MouseActionRelease,
					Button: MouseButtonMiddle,
					isSGR:  true,
				},
			},
		},
		{
			name: "right",
			buf:  encode(2, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseRight,
					Action: MouseActionPress,
					Button: MouseButtonRight,
					isSGR:  true,
				},
			},
		},
		{
			name: "right release",
			buf:  encode(2, 32, 16, true),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseRelease,
					Action: MouseActionRelease,
					Button: MouseButtonRight,
					isSGR:  true,
				},
			},
		},
		{
			name: "motion",
			buf:  encode(35, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseMotion,
					Action: MouseActionMotion,
					Button: MouseButtonNone,
					isSGR:  true,
				},
			},
		},
		{
			name: "wheel up",
			buf:  encode(64, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseWheelUp,
					Action: MouseActionPress,
					Button: MouseButtonWheelUp,
					isSGR:  true,
				},
			},
		},
		{
			name: "wheel down",
			buf:  encode(65, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseWheelDown,
					Action: MouseActionPress,
					Button: MouseButtonWheelDown,
					isSGR:  true,
				},
			},
		},
		{
			name: "wheel left",
			buf:  encode(66, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseWheelLeft,
					Action: MouseActionPress,
					Button: MouseButtonWheelLeft,
					isSGR:  true,
				},
			},
		},
		{
			name: "wheel right",
			buf:  encode(67, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseWheelRight,
					Action: MouseActionPress,
					Button: MouseButtonWheelRight,
					isSGR:  true,
				},
			},
		},
		{
			name: "backward",
			buf:  encode(128, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseBackward,
					Action: MouseActionPress,
					Button: MouseButtonBackward,
					isSGR:  true,
				},
			},
		},
		{
			name: "backward in motion",
			buf:  encode(160, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseBackward,
					Action: MouseActionMotion,
					Button: MouseButtonBackward,
					isSGR:  true,
				},
			},
		},
		{
			name: "forward",
			buf:  encode(129, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseForward,
					Action: MouseActionPress,
					Button: MouseButtonForward,
					isSGR:  true,
				},
			},
		},
		{
			name: "forward in motion",
			buf:  encode(161, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseForward,
					Action: MouseActionMotion,
					Button: MouseButtonForward,
					isSGR:  true,
				},
			},
		},
		// Combinations.
		{
			name: "alt+right",
			buf:  encode(10, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Alt:    true,
					Type:   MouseRight,
					Action: MouseActionPress,
					Button: MouseButtonRight,
					isSGR:  true,
				},
			},
		},
		{
			name: "ctrl+right",
			buf:  encode(18, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Ctrl:   true,
					Type:   MouseRight,
					Action: MouseActionPress,
					Button: MouseButtonRight,
					isSGR:  true,
				},
			},
		},
		{
			name: "ctrl+alt+right",
			buf:  encode(26, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Alt:    true,
					Ctrl:   true,
					Type:   MouseRight,
					Action: MouseActionPress,
					Button: MouseButtonRight,
					isSGR:  true,
				},
			},
		},
		{
			name: "alt+wheel press",
			buf:  encode(73, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Alt:    true,
					Type:   MouseWheelDown,
					Action: MouseActionPress,
					Button: MouseButtonWheelDown,
					isSGR:  true,
				},
			},
		},
		{
			name: "ctrl+wheel press",
			buf:  encode(81, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Ctrl:   true,
					Type:   MouseWheelDown,
					Action: MouseActionPress,
					Button: MouseButtonWheelDown,
					isSGR:  true,
				},
			},
		},
		{
			name: "ctrl+alt+wheel press",
			buf:  encode(89, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Alt:    true,
					Ctrl:   true,
					Type:   MouseWheelDown,
					Action: MouseActionPress,
					Button: MouseButtonWheelDown,
					isSGR:  true,
				},
			},
		},
		{
			name: "ctrl+alt+shift+wheel press",
			buf:  encode(93, 32, 16, false),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Shift:  true,
					Alt:    true,
					Ctrl:   true,
					Type:   MouseWheelDown,
					Action: MouseActionPress,
					Button: MouseButtonWheelDown,
					isSGR:  true,
				},
			},
		},
		// Batched events.
		{
			name: "batched events",
			buf:  encode(0, 32, 16, false) + encode(35, 40, 30, false) + encode(0, 64, 32, true),
			expected: []MouseEvent{
				{
					X:      32,
					Y:      16,
					Type:   MouseLeft,
					Action: MouseActionPress,
					Button: MouseButtonLeft,
					isSGR:  true,
				},
				{
					X:      40,
					Y:      30,
					Type:   MouseMotion,
					Action: MouseActionMotion,
					Button: MouseButtonNone,
					isSGR:  true,
				},
				{
					X:      64,
					Y:      32,
					Type:   MouseRelease,
					Action: MouseActionRelease,
					Button: MouseButtonLeft,
					isSGR:  true,
				},
			},
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			actual, err := parseSGRMouseEvents(tc.buf)
			if err != nil {
				t.Fatalf("unexpected error for test: %v",
					err,
				)
			}

			for i := range tc.expected {
				if tc.expected[i] != actual[i] {
					t.Fatalf("expected %#v but got %#v",
						tc.expected[i],
						actual[i],
					)
				}
			}
		})
	}
}
