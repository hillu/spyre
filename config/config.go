package config

import (
	"github.com/spf13/afero"
	"github.com/spf13/pflag"

	"github.com/dcso/spyre"
	"github.com/dcso/spyre/log"

	"os"
	"strings"
)

var (
	Paths         = simpleStringSlice(defaultPaths)
	MaxFileSize   = fileSize(32 * 1024 * 1024)
	ReportTargets = simpleStringSlice([]string{"spyre.log"})
	Hostname      string
	HighPriority  bool
	YaraFiles     simpleStringSlice
	IocFiles      simpleStringSlice
)

// Fs is the "filesystem" in which configuration and rules are found.
// This can be provided through a ZIP file appended to the binary.
var Fs afero.Fs

func Init() error {
	pflag.VarP(&Paths, "path", "p", "paths to be scanned (default: / on Unix, all fixed drives on Windows)")
	pflag.Var(&YaraFiles, "yara-rule-files",
		"yara files to be used for file scan (default: search recursively for files matching *.yr, *.yar, *.yara)")
	pflag.Var(&IocFiles, "ioc-files",
		"IOC files to be used for descriptive IOCs (default: ioc.json)")
	pflag.Var(&MaxFileSize, "max-file-size",
		"maximum size of individual files to be scanned, turn off by setting to 0 or negative value")
	pflag.StringVar(&spyre.Hostname, "set-hostname", spyre.DefaultHostname, "hostname")
	pflag.VarP(&log.GlobalLevel, "loglevel", "l", "loglevel")
	pflag.VarP(&ReportTargets, "report", "r", "report target(s)")
	pflag.BoolVar(&HighPriority, "high-priority", false,
		"run at high priority instead of giving up CPU and I/O resources to other processes")

	var args []string
	if len(os.Args) > 1 {
		log.Debug("Using user-provided command line parameters.")
		args = os.Args[1:]
	} else if buf, err := afero.ReadFile(Fs, "params.txt"); err != nil {
		log.Debug("Using default parameters.")
	} else {
		log.Debug("Using parametes form params.txt.")
		for _, line := range strings.Split(string(buf), "\n") {
			line = strings.TrimSpace(line)
			if len(line) == 0 || line[0] == '#' {
				continue
			}
			if tokens := strings.Fields(line); len(tokens) > 1 && !strings.Contains(tokens[0], "=") {
				args = append(args, tokens[0])
				args = append(args, strings.Join(tokens[1:len(tokens)], " "))
			} else {
				args = append(args, line)
			}
		}
	}
	pflag.CommandLine.Parse(args)

	pflag.VisitAll(func(f *pflag.Flag) {
		log.Debugf("config: --%s %s%s", f.Name, f.Value, map[bool]string{false: " (unchanged)"}[f.Changed])
	})

	log.Init()
	return nil
}
