package file_util

import (
	"github.com/NoahAmethyst/dispatch-center/constant"
	"os"
)

var fileRoot string

func GetFileRoot() string {
	if len(fileRoot) == 0 {
		fileRoot = os.Getenv(constant.FILE_ROOT)
		if len(fileRoot) == 0 {
			fileRoot = "/tmp"
		}
	}
	return fileRoot
}

func SetFileRoot(newRoot string) {
	fileRoot = newRoot
}
