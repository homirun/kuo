package test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestCommandKuo(t *testing.T) {
	InitializeTest()

	out, err := exec.Command("kubectl", "kuo").Output()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "[context1 context2]")
}

func InitializeTest() {
	exec.Command("kubectl", "kuo", "set", "context1", "context2").Run()
}
