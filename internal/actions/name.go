package actions

import (
	"fmt"
	"strings"
)

type Name string

const (
	NameUpdate Name = "update"
	NameGet    Name = "get"
)

var Names = []string{string(NameUpdate), string(NameGet)}

func NameFromString(val string) (Name, error) {
	if val == "" {
		return "", fmt.Errorf(
			"invalid action string. Available actions: [%s]",
			strings.Join(Names, ", "),
		)
	}

	for _, action := range Names {
		if action == val {
			return Name(action), nil
		}
	}

	return "", fmt.Errorf(
		"unexpected action. Available actions: [%s]",
		strings.Join(Names, ", "),
	)
}
