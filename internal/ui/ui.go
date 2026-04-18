package ui

import "github.com/fatih/color"

var (
	Success = color.New(color.FgGreen, color.Bold).PrintfFunc()
	Info    = color.New(color.FgCyan).PrintfFunc()
	Warn    = color.New(color.FgYellow).PrintfFunc()
	Error   = color.New(color.FgRed, color.Bold).PrintfFunc()
	Bold    = color.New(color.Bold).PrintfFunc()
)
