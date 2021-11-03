package netem

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Netem struct {
	path string
	Option
}

type Option struct {
	NetworkIface string
}

// New create netem configuration
func New(opt Option) (*Netem, error) {
	tcPath, err := exec.LookPath("tc")
	if err != nil {
		return nil, err
	}

	netem := &Netem{
		path:   tcPath,
		Option: opt,
	}

	return netem, nil
}

// AddDelay add delay for given duration to spesific interface
func (netem *Netem) AddDelay(d time.Duration) error {
	cmd := []string{"add", "dev", netem.NetworkIface, "root", "netem", "delay", getDuration(d) + "ms"}
	return netem.run(cmd...)
}

// DeleteDelay delete delay for given duration from spesific interface
func (netem *Netem) DeleteDelay(d time.Duration) error {
	cmd := []string{"delete", "dev", netem.NetworkIface, "root", "netem", "delay", getDuration(d) + "ms"}
	return netem.run(cmd...)
}

// ChangeDelay modify current delay from spesific interface
func (netem *Netem) ChangeDelay(d time.Duration) error {
	cmd := []string{"change", "dev", netem.NetworkIface, "root", "netem", "delay", getDuration(d) + "ms"}
	return netem.run(cmd...)
}

// Show display all rules attached to spesific interface
func (netem *Netem) Show() ([]string, error) {
	cmd := []string{"show", "dev", netem.NetworkIface}

	var stdout bytes.Buffer

	if err := netem.runWithOutput(cmd, &stdout); err != nil {
		return nil, err
	}

	rules := strings.Split(stdout.String(), "\n")

	return rules, nil
}

// getDuration make time duration formatted as string
func getDuration(d time.Duration) string {
	return strconv.Itoa(int(durationToMs(d)))
}

// durationToMs convert duration to milliseconds
func durationToMs(d time.Duration) int64 {
	return d.Milliseconds()
}

// run tc command with given arguments
func (netem *Netem) run(args ...string) error {
	return netem.runWithOutput(args, nil)
}

// runWithOuput tc command with given arguments with output
func (netem *Netem) runWithOutput(args []string, stdout io.Writer) error {
	var stderr bytes.Buffer

	newArgs := append([]string{netem.path, "qdisc"}, args...)

	cmd := exec.Cmd{
		Path:   netem.path,
		Args:   newArgs,
		Stderr: &stderr,
		Stdout: stdout,
	}

	if err := cmd.Run(); err != nil {
		switch err.(type) {
		case *exec.ExitError:
			return fmt.Errorf(stderr.String())
		default:
			return err
		}
	}

	return nil
}
