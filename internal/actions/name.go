package actions

import (
	"fmt"
	"strings"
)

type Name string

const (
	NameUpImageVersion Name = "up-image-version"
	NameAppend         Name = "append"
)

var Names = []string{string(NameUpImageVersion), string(NameAppend)}

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
