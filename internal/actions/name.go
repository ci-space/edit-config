package actions

import (
	"fmt"
	"strings"
)

type Name string

const (
	NameUpdate         Name = "update"
	NameGet            Name = "get"
	NameUpImageVersion Name = "up-image-version"
)

var Names = []string{string(NameUpdate), string(NameGet), string(NameUpImageVersion)}

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
