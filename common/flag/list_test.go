package flag

import (
	"reflect"
	"testing"
)

func TestList_SetWithEmptyString(t *testing.T) {
	list := &List{}

	value := ""
	err := list.Set(value)
	if err != nil {
		t.Error("error in processing empty string to list flag")
	}

	if !reflect.DeepEqual([]string{}, list.Value) {
		t.Error("error in splitting empty string to list flag")
	}
}

func TestList_SetWithNormalString(t *testing.T) {
	list := &List{}

	value := "node_systemoom,node_oom"
	err := list.Set(value)
	if err != nil {
		t.Error("error in processing normal string to list flag")
	}

	if !reflect.DeepEqual([]string{"node_systemoom", "node_oom"}, list.Value) {
		t.Error("error in normal string to list flag")
	}
}

func TestList_SetWithNormalStringContainingSpace(t *testing.T) {
	list := &List{}

	value := "node_systemoom , node_oom "
	err := list.Set(value)
	if err != nil {
		t.Error("error in processing normal string with space to list flag")
	}

	if !reflect.DeepEqual([]string{"node_systemoom", "node_oom"}, list.Value) {
		t.Error("error in normal string with space to list flag")
	}
}

func TestList_SetWithAbnormalString(t *testing.T) {
	list := &List{}

	value := "node_systemoom , "
	err := list.Set(value)
	if err != nil {
		t.Error("error in processing abnormal string to list flag")
	}

	if !reflect.DeepEqual([]string{"node_systemoom"}, list.Value) {
		t.Error("error in abnormal string to list flag")
	}
}
