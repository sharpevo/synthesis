package commandparser

import (
//"fmt"
)

type CommandInterface interface {
	Execute(args ...string) (interface{}, error)
	// TODO: rb
}

type CommandType struct {
	title string
	path  string
}

func (c CommandType) Title() string {
	return c.title
}

func (c *CommandType) SetTitle(title string) {
	c.title = title
}

func (c *CommandType) Execute(args ...string) (interface{}, error) {
	return "", nil
}

type CommandImportType struct {
	CommandType
}

var CmdImport CommandImportType

type CommandAsyncType struct {
	CommandType
}

var CmdAsync CommandAsyncType

type CommandRetryType struct {
	CommandType
}

var CmdRetry CommandRetryType

func Init() {
	CmdImport.SetTitle("IMPORT")
	CmdAsync.SetTitle("ASYNC")
	CmdRetry.SetTitle("RETRY")
}
