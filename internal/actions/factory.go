package actions

import "github.com/ci-space/edit-config/internal/fs"

func CreateActions(fs fs.Filesystem) map[Name]Action {
	return map[Name]Action{
		NameGet:            NewGetAction(),
		NameUpdate:         NewUpdateAction(fs),
		NameUpImageVersion: NewUpImageVersionAction(fs),
	}
}
