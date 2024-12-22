package game

import "errors"

type ACTION int

const (
	MOVE_NONE ACTION = iota
	MOVE_LEFT
	MOVE_RIGHT
	MOVE_UP
	MOVE_DOWN
)

func (action ACTION) GetActionString() string {
	actionMapping := map[ACTION]string{
		MOVE_UP:    "u",
		MOVE_DOWN:  "d",
		MOVE_LEFT:  "l",
		MOVE_RIGHT: "r",
	}
	return actionMapping[action]
}

func LoadActionFromString(actionName string) (ACTION, error) {
	mapping := map[string]ACTION{
		"u": MOVE_UP,
		"d": MOVE_DOWN,
		"l": MOVE_LEFT,
		"r": MOVE_RIGHT,
	}

	if val, ok := mapping[actionName]; ok {
		return val, nil
	}

	return ACTION(0), errors.New("invalid action name")
}
