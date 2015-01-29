package main

import (
	"os"
	"syscall"

	"github.com/gosexy/gettext"
	"github.com/lxc/lxd"
	"golang.org/x/crypto/ssh/terminal"
)

type execCmd struct{}

func (c *execCmd) usage() string {
	return gettext.Gettext(
		"exec specified command in a container.\n" +
			"\n" +
			"lxc exec container [command]\n")
}

func (c *execCmd) flags() {}

func (c *execCmd) run(config *lxd.Config, args []string) error {
	if len(args) < 2 {
		return errArgs
	}

	remote, name := config.ParseRemoteAndContainer(args[0])
	d, err := lxd.NewClient(config, remote)
	if err != nil {
		return err
	}

	cfd := syscall.Stdout
	if terminal.IsTerminal(cfd) {
		oldttystate, err := terminal.MakeRaw(cfd)
		if err != nil {
			return err
		}
		defer terminal.Restore(cfd, oldttystate)
	}

	// TODO: we should exit with the same exit code as the command in the
	// container.
	return d.Exec(name, args[1:], os.Stdin, os.Stdout, os.Stderr)
}