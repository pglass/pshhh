package test

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: eventually it would be nice to be able to substitute "/bin/bash" here
var PSH_EXE = "../psh"

func exec_psh(args ...string) (string, error) {
	cmd := exec.Command(PSH_EXE, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

type exeData struct {
	Args     []string
	Output   string
	ExitCode int
}

func run_psh_test(t *testing.T, data exeData) {
	t.Logf("Args: %q", data.Args)
	t.Logf("ExitCode (expected): %v", data.ExitCode)
	t.Logf("Output (expected): %v", data.Output)

	out, err := exec_psh(data.Args...)

	t.Logf("Error (received): %v", err)
	t.Logf("Output (received): %v", out)

	check_exit_code(t, data.ExitCode, err)
	assert.Equal(t, data.Output, out)

}

func check_exit_code(t *testing.T, expected int, err error) {
	if expected == 0 {
		assert.Nil(t, err, "Expected exit status %v, but got error %v", expected, err)
	} else {
		switch e := err.(type) {
		case *exec.ExitError:
			// ugh
			if status, ok := e.Sys().(syscall.WaitStatus); ok {
				assert.Equal(t, expected, status.ExitStatus())
				return
			}
		}
		t.Errorf("Failed to get exit status from error: %v", err)
	}
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

	/* Use Default Values (:-)
	 *
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

	/* Use Default Values (-)
	 *
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

	/* Use Alternative Value (+)
	 *
	 * 				P set, not null 		P set, but null 	P not set
	 * ${P+word} 	substitute W 			substitute W		substitute null
	 */
	exeData{
		Args:   []string{"-t", `/bin/echo ${X+word}`},
		Output: "\n",
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X+word}`, "-e", "X="},
		Output: "word\n",
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X+word}`, "-e", "X=param"},
		Output: "word\n",
	},

	/* Use Alternative Value (:+)
	 *
	 * 				P set, not null 		P set, but null 	P not set
	 * ${P:+word} 	substitute W 			substitute null		substitute null
	 */
	exeData{
		Args:   []string{"-t", `/bin/echo ${X:+word}`},
		Output: "\n",
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X:+word}`, "-e", "X="},
		Output: "\n",
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X:+word}`, "-e", "X=param"},
		Output: "word\n",
	},

	/* Indicate Error if Null or Unset (?)
	 *
	 * 				P set, not null 		P set, but null 	P not set
	 * ${P?word} 	substitute P 			substitute null		error, exit
	 */
	exeData{
		Args:     []string{"-t", `/bin/echo ${X?word}`},
		Output:   "error: X: word\n",
		ExitCode: 1,
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X?word}`, "-e", "X="},
		Output: "\n",
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X?word}`, "-e", "X=param"},
		Output: "param\n",
	},

	/* Indicate Error if Null or Unset (:?)
	 *
	 * 				P set, not null 		P set, but null 	P not set
	 * ${P:?word} 	substitute P 			error, exit			error, exit
	 */
	exeData{
		Args:     []string{"-t", `/bin/echo ${X:?word}`},
		Output:   "error: X: word\n",
		ExitCode: 1,
	},
	exeData{
		Args:     []string{"-t", `/bin/echo ${X:?word}`, "-e", "X="},
		Output:   "error: X: word\n",
		ExitCode: 1,
	},
	exeData{
		Args:   []string{"-t", `/bin/echo ${X:?word}`, "-e", "X=param"},
		Output: "param\n",
	},
}
