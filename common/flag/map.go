package flag

import (
	"fmt"
	"strings"
)

const (
	pairseperator     string = ","
	keyValueSeperator string = "="
)

type Map struct {
	Value map[string]string
}

func (m *Map) Set(value string) error {
	m.Value = map[string]string{}

	preSingleValues := strings.Split(value, pairseperator)
	for _, preSingleValue := range preSingleValues {
		if strings.TrimSpace(preSingleValue) != "" {
			pair := strings.SplitN(preSingleValue, "=", 2)
			if len(pair) != 2 {
				return fmt.Errorf("key value pair should be seperated by %s", keyValueSeperator)
			}

			key := strings.TrimSpace(pair[0])
			value := strings.TrimSpace(pair[1])

			m.Value[key] = value
		}
	}

	return nil
}

func (m *Map) String() string {
	return fmt.Sprintf("%s", m.Value)
}
