// datfile package provides you the ability to save each variable in an
// independent `.dat` file storing whatever content you want (ex: string, etc).
package datfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func NewString(variableName, variableData string) error {
	filename := fmt.Sprintf("./%v.dat", variableName)
	filename = strings.Replace(filename, " ", "_", -1)
	return ioutil.WriteFile(filename, []byte(variableData), 0644)
}

func ReadString(variableName string) (string, error) {
	filename := fmt.Sprintf("./%v.dat", variableName)
	filename = strings.Replace(filename, " ", "_", -1)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return "", fmt.Errorf("file '%s' does not exist", filename)
		case os.IsPermission(err):
			return "", fmt.Errorf("permission denied to read file '%s'", filename)
		default:
			return "", fmt.Errorf("error reading file '%s': %w", filename, err)
		}
	}
	return string(data), nil
}
