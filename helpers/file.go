package helpers

import (
	"io/ioutil"
)

func ReadFile(fName string) ([]byte, error) {
	return ioutil.ReadFile(fName)
}
