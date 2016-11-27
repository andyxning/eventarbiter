package flag

import (
	"reflect"
	"testing"
)

func TestMap_SetWithEmptyString(t *testing.T) {
	m := &Map{}

	value := ""
	err := m.Set(value)
	if err != nil {
		t.Error("error in processing empty string to map flag")
	}

	if !reflect.DeepEqual(map[string]string{}, m.Value) {
		t.Error("error in splitting empty string to map flag")
	}
}

func TestMap_SetWithNormalString(t *testing.T) {
	m := &Map{}

	value := "key1=value1,key2=value2"
	err := m.Set(value)
	if err != nil {
		t.Error("error in processing normal string to map flag")
	}

	if !reflect.DeepEqual(map[string]string{"key1": "value1", "key2": "value2"}, m.Value) {
		t.Error("error in normal string to map flag")
	}
}

func TestMap_SetWithNormalStringContainingSpace(t *testing.T) {
	m := &Map{}

	value := "key1=value1 , key2=value2 "
	err := m.Set(value)
	if err != nil {
		t.Error("error in processing normal string with space to map flag")
	}

	if !reflect.DeepEqual(map[string]string{"key1": "value1", "key2": "value2"}, m.Value) {
		t.Error("error in normal string with space to map flag")
	}
}

func TestMap_SetWithAbnormalString(t *testing.T) {
	m := &Map{}

	value := "key1=value1 , "
	err := m.Set(value)
	if err != nil {
		t.Error("error in processing abnormal string to map flag")
	}

	if !reflect.DeepEqual(map[string]string{"key1": "value1"}, m.Value) {
		t.Error("error in abnormal string to map flag")
	}
}
