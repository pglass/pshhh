package test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var PSH_EXE = "../psh"

func exec_psh(args ...string) (string, error) {
	cmd := exec.Command(PSH_EXE, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

type exeData struct {
	Args   []string
	Output string
	Error  error
}

func run_psh_test(t *testing.T, data exeData) {
	t.Logf("Args: %q", data.Args)
	t.Logf("Error (expected): %v", data.Error)
	t.Logf("Output (expected): %v", data.Output)

	out, err := exec_psh(data.Args...)

	t.Logf("Error (received): %v", err)
	t.Logf("Output (received): %v", out)

	if data.Error != nil {
		assert.Equal(t, data.Error, err)
	} else {
		assert.Nil(t, err)
	}

	assert.Equal(t, data.Output, out)

}

func TestPsh(t *testing.T) {
	for _, data := range PSH_CASES {
		name := strings.Join(data.Args, " ")
		t.Run(name, func(t *testing.T) { run_psh_test(t, data) })
	}
}

// TODO: assumes /bin/echo exists
var PSH_CASES = []exeData{
	exeData{
		Args:   []string{"-t", "echo"},
		Output: "ERROR: failed to run [echo]: no such file or directory\n",
	},
	exeData{
		Args:   []string{"-t", "echo", "-e", "PATH=/bin"},
		Output: "\n",
	},
	exeData{
		Args:   []string{"-t", "echo $PATH", "-e", "PATH=/bin"},
		Output: "/bin\n",
	},
	exeData{
		Args:   []string{"-t", `echo $PATH $WUMBO`, "-e", "PATH=/bin", "-e", "WUMBO=mini"},
		Output: "/bin mini\n",
	},
	exeData{
		Args:   []string{"-t", `echo "$PATH $WUMBO"`, "-e", "PATH=/bin", "-e", "WUMBO=mini"},
		Output: "/bin mini\n",
	},
	exeData{
		Args:   []string{"-t", `echo ${PATH} "${WUMBO}"`, "-e", "PATH=/bin", "-e", "WUMBO=mini"},
		Output: "/bin mini\n",
	},

	/*
	 * 				P set, not null 		P set, but null 	P not set
	 *				--------------- 		--------------- 	---------------
	 * ${P:-W} 		substitute P 			substitute W 		substitute W
	 */
	exeData{
		Args:   []string{"-t", `/bin/echo ${X:-word}`},
		Output: "word\n",
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X:-word}`, "-e", "X="},
		Output: "word\n",
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X:-word}`, "-e", "X=param"},
		Output: "param\n",
	},
	/*
	 * 				P set, not null 		P set, but null 	P not set
	 * ${P-word} 	substitute P 			substitute null		substitute W
	 */
	exeData{
		Args:   []string{"-t", `/bin/echo ${X-word}`},
		Output: "word\n",
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X-word}`, "-e", "X="},
		Output: "\n",
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X-word}`, "-e", "X=param"},
		Output: "param\n",
	},
}
