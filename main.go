package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
)

func init() {

	//	cmdRun.Flag.Var(&mainFiles, "main", "specify main go files")
	//	cmdRun.Flag.Var(&gendoc, "gendoc", "auto generate the docs")
	//	cmdRun.Flag.Var(&downdoc, "downdoc", "auto download swagger file when not exist")
	//	cmdRun.Flag.Var(&excludedPaths, "e", "Excluded paths[].")
}

const version = "0.1.0"

type Command struct {
	Run         func(cmd *Command, args []string) int
	UsageLine   string
	Long        template.HTML
	Short       template.HTML
	Flag        flag.FlagSet
	CustomFlags bool
}

func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}
func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(string(c.Long)))
	os.Exit(2)
}
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var commands = []*Command{
	cmdRun,
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)
	args := flag.Args()
	if len(args) < 1 {
		cmd := commands[0]
		cmd.Flag.Usage = func() { cmd.Usage() }
		os.Exit(cmd.Run(cmd, nil))
		//		usage()
	}
	if args[0] == "help" {
		help(args[1:])
	}
	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			os.Exit(cmd.Run(cmd, args))
			return
		}

	}
	fmt.Fprintf(os.Stderr, "gotester: unknown subcommand %q\nRun 'gotester help' for usage.\n", args[0])
	os.Exit(2)
}

var usageTemplate = `GoTester is a tool for goagent`
var helpTemplate = `{{if .Runnable}}usage: gotester {{.UsageLine}}
{{end}}{{.Long | trim}}`

func usage() {
	tmpl(os.Stdout, usageTemplate, commands)
}
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")

	t.Funcs(template.FuncMap{"trim": func(s template.HTML) template.HTML {
		return template.HTML(strings.TrimSpace(string(s)))
	}})

	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
func help(args []string) {
	if len(args) == 0 {
		usage()
		return
	}
	if len(args) != 1 {
		fmt.Fprint(os.Stdout, "usage:gotester help command\n\nToo many arguments given.\n")
		os.Exit(2)
	}
	arg := args[0]
	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}
	fmt.Fprintf(os.Stdout, "Unknown help topic.%#q. Run 'gotester help'.\n", arg)
	os.Exit(2)
}
