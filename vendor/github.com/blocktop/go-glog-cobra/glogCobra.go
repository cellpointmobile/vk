package glogcobra

import (
	"strconv"
	"github.com/spf13/cobra"
	"flag"
)

const (
	// LogToStdErr specifies that logs are written to standard error instead of to files.
	LogToStdErr = "logtostderr"
	// AlsoLogToStdErr specifes that logs are written to standard error as well as to files.
	AlsoLogToStdErr = "alsologtostderr"
	// StdErrThreshold specifies the severity threshold at or above which events are logged to standard error as well as to files.
	StdErrThreshold = "stderrthreshold"
	// LogDir specifies the directory to which log files will be written instead of the default temporary directory.""
	LogDir = "log_dir"
	// LogBacktraceAt is set to a file and line number holding a logging statement,
  // such as
  //     -log_backtrace_at=gopherflakes.go:234
  // to write a stack trace to the Info log whenever execution
  // hits that statement. (Unlike with -vmodule, the ".go" must be
  // present.)
	LogBacktraceAt = "log_backtrace_at"
	// Verbosity (shortcut -v) enables V-leveled logging at the specified level.
	Verbosity = "verbosity"
	// VModule specifies the verbosity for a given module. The syntax of the argument 
	// is a comma-separated list of pattern=N,
  // where pattern is a literal file name (minus the ".go" suffix) or
  // "glob" pattern and N is a V level. For instance,
	//     -vmodule=gopher*=3
  // sets the V level to 3 in all Go files whose names begin "gopher".
	VModule = "vmodule"
)

// Init loads the flags required by glog into the cobra command's
// PersistantFlags collection.
func Init(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.Bool(LogToStdErr, false, "Logs are written to standard error instead of to files.")
	flags.Bool(AlsoLogToStdErr, false, "Logs are written to standard error as well as to files.")
  flags.String(StdErrThreshold, "ERROR", `Log events at or above this severity are logged to standard
error as well as to files.`)
	flags.String(LogDir, "",	`Log files will be written to this directory instead of the
default temporary directory.`)
  flags.String(LogBacktraceAt, "", `When set to a file and line number holding a logging statement,
such as
	  -log_backtrace_at=gopherflakes.go:234
a stack trace will be written to the Info log whenever execution
hits that statement. (Unlike with -vmodule, the ".go" must be
present.)`)
	flags.IntP(Verbosity, "v", 0, "Enable V-leveled logging at the specified level.")
	flags.String(VModule, "", `The syntax of the argument is a comma-separated list of pattern=N,
where pattern is a literal file name (minus the ".go" suffix) or
"glob" pattern and N is a V level. For instance,
		-vmodule=gopher*=3
sets the V level to 3 in all Go files whose names begin "gopher".`)
}

// Parse reads flag valuess from the cobra command's peristant flags collection
// and sets them in the flag module for glog. flag.Parse() is called if it has not
// already been called.
func Parse(cmd *cobra.Command) error {
	flags := cmd.PersistentFlags()
	
	fb, err := flags.GetBool(LogToStdErr)
	if err != nil {
		return err
	}
	flag.Set(LogToStdErr, strconv.FormatBool(fb))

	fb, err = flags.GetBool(AlsoLogToStdErr)
	if err != nil {
		return err
	}
	flag.Set(AlsoLogToStdErr, strconv.FormatBool(fb))

	fs, err := flags.GetString(StdErrThreshold)
	if err != nil {
		return err
	}
	flag.Set(StdErrThreshold, fs)

	fs, err = flags.GetString(LogDir)
	if err != nil {
		return err
	}
	flag.Set(LogDir, fs)
	
	fs, err = flags.GetString(LogBacktraceAt)
	if err != nil {
		return err
	}
	flag.Set(LogBacktraceAt, fs)

	fi, err := flags.GetInt(Verbosity)
	if err != nil {
		return err
	}
	flag.Set("v", strconv.FormatInt(int64(fi), 10))

	fs, err = flags.GetString(VModule)
	if err != nil {
		return err
	}
	flag.Set(VModule, fs)

	if !flag.Parsed() {
		flag.Parse()
	}
	return nil
}