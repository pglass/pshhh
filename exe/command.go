package exe

import (
	"log"
	"os"
	"path"
	"strings"
	"syscall"
)

type PshProc struct {
	Name         string
	Args         []string
	ProcAttr     *syscall.ProcAttr
	IsBackground bool
}

func NewPshProc(args []string, env []string) (*PshProc, error) {
	work_dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	proc := &PshProc{
		Name: args[0],
		Args: args,
		ProcAttr: &syscall.ProcAttr{
			Dir:   work_dir,
			Env:   env,
			Files: []uintptr{0, 1, 2},
		},
	}
	return proc, nil
}

func (c *PshProc) ForkExec() (int, error) {
	command_name := c.PathLookup(c.Name)

	pid, err := syscall.ForkExec(command_name, c.Args, c.ProcAttr)
	if err != nil {
		return pid, err
	}

	// logging
	if c.IsBackground {
		log.Printf("Forked pid = %v %v [background]", pid, c.Args)
	} else {
		log.Printf("Forked pid = %v %v", pid, c.Args)
	}

	if !c.IsBackground {
		waitstatus := syscall.WaitStatus(0)
		rusage := syscall.Rusage{}

		// https://linux.die.net/man/2/wait
		wait_opts := syscall.WEXITED | syscall.WSTOPPED

		// wait for the process to exit
		wpid, err := syscall.Wait4(pid, &waitstatus, wait_opts, &rusage)
		if err != nil {
			log.Printf("%v", err)
		} else if wpid != pid {
			log.Printf("Wait4 return non-matching pid %v (expected %v). Did process %v exit?", wpid, pid, pid)
		}
	}

	return pid, err
}

/* Lookup the name on the path. This works as follows
 *   - If name is a path (starts with "/" or "./")
 *   - Else, find the first directory, dir, in the PATH variable
 *	 containing a file called name. Return the path of that file.
 *   - Otherwise, return name by default.
 *
 * TODO: This doesn't belong here. Command search must find existing functions,
 *		 shell builtins, etc.
 */
func (c *PshProc) PathLookup(name string) string {
	if strings.Contains(name, "/") {
		log.Printf("PathLookup: %q is already a path", name)
		return name
	}

	path_var := c.FetchEnvVar("PATH")
	for _, dir := range strings.Split(path_var, ":") {
		if dir == "" {
			continue
		}

		check_file := path.Join(dir, name)
		if _, err := os.Stat(check_file); err == nil {
			log.Printf("PathLookup: Found %q at %q in dir %q", name, check_file, dir)
			return check_file
		} else {
			log.Printf("PathLookup: Tried %q at %q: %v", name, check_file, err)
		}
	}
	return name
}

/* c.ProcAttr.Env stores environment variables as a list of "<key>=<value>"
 * strings. This fetches the <value> portion given the <key>, or returns
 * empty string.
 */
func (c *PshProc) FetchEnvVar(key string) string {
	key = key + "="
	for _, item := range c.ProcAttr.Env {
		if strings.HasPrefix(item, key) {
			return strings.SplitN(item, "=", 2)[1]
		}
	}
	return ""
}
