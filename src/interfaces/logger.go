package interfaces

import "github.com/fatih/color"

var WarningColor = color.New(color.FgYellow).Add(color.Bold)
var ErrorColor = color.New(color.FgRed).Add(color.Underline)
var InfoColor = color.New(color.FgHiCyan)
var CompleteColor = color.New(color.BgGreen).Add(color.FgWhite).Add(color.Bold)

type Logger interface {
	PrintSuccessLogMessage(msg string) error
	PrintWarningLogMessage(msg string) error
	PrintErrorLogMessage(msg string) error
	PrintInfoLogMessage(msg string) error
}
