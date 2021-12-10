package archive

import (
	"fmt"
	"one-backup/cmd"
)

const Debug = false

func ArchiveTar(remove bool, directory, srcName, destName string) {
	cmdStr := ""
	if remove == true {
		cmdStr = fmt.Sprintf("cd %v && tar zcf %v %v --remove-files", directory, destName, srcName)
	} else {
		cmdStr = fmt.Sprintf("cd %v && tar zcf %v %v ", directory, destName, srcName)
	}
	cmd.Run(cmdStr, Debug)

}
func UnarchiveTar(filename, directory string) {
	cmdStr := fmt.Sprintf("tar xf %v -C %v", filename, directory)
	cmd.Run(cmdStr, Debug)
}
