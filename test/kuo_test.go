package test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestCommandKuo(t *testing.T) {
	var out []byte
	var err error

	// .kuoconfigがないときはエラーになる
	exec.Command("rm", "~/.kuoconfig").Run()
	out, err = exec.Command("kubectl", "kuo").Output()
	assert.Error(t, err)

	// .kuoconfigがあるときは設定されたcontextの一覧が帰ってくる
	InitializeSetContext()
	out, err = exec.Command("kubectl", "kuo").Output()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "[context1 context2]")
}

func TestCommandKuoSet(t *testing.T) {
	var out []byte
	var err error

	// kuo setにはargsが2つ必要
	_, err = exec.Command("kubectl", "kuo", "set").Output()
	assert.Error(t, err)
	_, err = exec.Command("kubectl", "kuo", "set", "context1").Output()
	assert.Error(t, err)

	// kuo set context1 context2はエラーにならず、メッセージが帰ってくる
	out, err = exec.Command("kubectl", "kuo", "set", "context1", "context2").Output()
	assert.Contains(t, string(out), "set .kuoconfig: [context1 context2]")
}

func InitializeSetContext() {
	exec.Command("kubectl", "kuo", "set", "context1", "context2").Run()
}
