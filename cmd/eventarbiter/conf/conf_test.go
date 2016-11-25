package conf

import (
	"k8s.io/heapster/common/flags"
	"net/url"
	"runtime"
	"testing"
)

func TestConfig_validateMaxProcsWithNumCPU(t *testing.T) {
	testConf := Config{
		MaxProcs: uint(runtime.NumCPU()),
	}

	if err := testConf.validateMaxProcs(); err != nil {
		t.Error("set max procs to max cpu number error")
	}
}

func TestConfig_validateMaxProcsWithNumCPUPlus1(t *testing.T) {
	testConf := Config{
		MaxProcs: uint(runtime.NumCPU() + 1),
	}

	if err := testConf.validateMaxProcs(); err == nil {
		t.Error("set max procs to max cpu number plus 1 correctly")
	}
}

func TestConfig_validateSourceWithKubernetes(t *testing.T) {
	testConf := Config{
		Source: flags.Uri{
			Key: "kubernetes",
		},
	}

	if err := testConf.validateSource(); err != nil {
		t.Error("set kubernetes source error")
	}
}

func TestConfig_validateSourceWithNoneKubernetes(t *testing.T) {
	testConf := Config{
		Source: flags.Uri{
			Key: "kubernetes1",
		},
	}

	if err := testConf.validateSource(); err == nil {
		t.Error("set none kubernetes source success")
	}
}

func TestConfig_validateSinkWithNoneSink(t *testing.T) {
	testConf := Config{
		Sinks: flags.Uris([]flags.Uri{}),
	}

	if err := testConf.validateSource(); err == nil {
		t.Error("empty source is allowed")
	}
}

func TestConfig_validateSinkWithInvalidSink(t *testing.T) {
	testConf := Config{
		Sinks: flags.Uris([]flags.Uri{{Key: "invalid"}}),
	}

	if err := testConf.validateSource(); err == nil {
		t.Error("invliad sink passed")
	}
}

func TestConfig_validateSinkWithEmptyCallbackSink(t *testing.T) {
	emptyURL, _ := url.Parse("")

	testConf := Config{
		Sinks: flags.Uris([]flags.Uri{{Key: "callback", Val: *emptyURL}}),
	}

	if err := testConf.validateSource(); err == nil {
		t.Error("empty callback sink passed")
	}
}
