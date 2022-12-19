package test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

const (
	CONTEXT1 string = "kind"
	CONTEXT2 string = "minikube"
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
	assert.Contains(t, string(out), "["+CONTEXT1+" "+CONTEXT2+"]")
}

func TestCommandKuoSet(t *testing.T) {
	var out []byte
	var err error

	// kuo setにはargsが2つ必要
	_, err = exec.Command("kubectl", "kuo", "set").Output()
	assert.Error(t, err)
	_, err = exec.Command("kubectl", "kuo", "set", "CONTEXT1").Output()
	assert.Error(t, err)

	// kuo set CONTEXT1 CONTEXT2で正常に処理されメッセージが帰ってくる
	out, err = exec.Command("kubectl", "kuo", "set", CONTEXT1, CONTEXT2).Output()
	assert.Contains(t, string(out), "set .kuoconfig: ["+CONTEXT1+" "+CONTEXT2+"]")
}

func TestCommandKuoExecKubectl(t *testing.T) {
	var out []byte
	var err error

	// kubectl getなどのbypassしたコマンドがエラーだったらkuoもエラーを返す
	_, err = exec.Command("kubectl", "kuo", "get").Output()
	assert.Error(t, err)

	// kubectl applyが実行されたら警告を出す
	out, err = exec.Command("kubectl", "kuo", "get").Output()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "It will be applied to multiple contexts.\nDo you want to continue? (y/n)\n")
}

func InitializeSetContext() {
	exec.Command("kubectl", "kuo", "set", CONTEXT1, CONTEXT2).Run()
}
