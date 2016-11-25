package flag

import "strings"

const seperator string = ","

type List struct {
	Value []string
}

func (list *List) Set(value string) error {
	list.Value = []string{}

	preSingleValues := strings.Split(value, seperator)
	for _, preSingleValue := range preSingleValues {
		postSingleValue := strings.TrimSpace(preSingleValue)
		if postSingleValue != "" {
			list.Value = append(list.Value, postSingleValue)
		}
	}

	return nil
}

func (list *List) String() string {
	return strings.Join(list.Value, seperator)
}
