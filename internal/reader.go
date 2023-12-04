package ami

import (
	"os"
	"strings"
)

func ReadScript(file string) ([]string, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	data := string(f)

	s := strings.Split(data, ",")

	return s, nil
}
