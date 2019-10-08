package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/brada954/restshell/shell"
)

type DirCommand struct {
	// Place getopt option value pointers here
	loadOption *bool
}

type Result struct {
	Message string
	Folder  string
	Files   []File
}

type File struct {
	Name      string
	Date      time.Time
	Size      int64
	Mode      string
	IsDir     bool
	FullPath  string
	Reference string
}

func NewDirCommand() *DirCommand {
	return &DirCommand{}
}

func (cmd *DirCommand) AddOptions(set shell.CmdSet) {
	set.SetParameters("value")
	cmd.loadOption = set.BoolLong("load", 0, "Load the files from the directory parameter into history")
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *DirCommand) ExecuteLegacy(args []string) error {

	dirArgs := osDirArgs
	for _, a := range args {
		dirArgs = append(dirArgs, a)
	}

	c := exec.Command(osDirCmd, dirArgs...)

	text, err := c.Output()
	if err != nil {
		return errors.New("Dir Error: " + err.Error())
	}
	fmt.Fprintf(shell.OutputWriter(), "%s\n", string(text))
	return nil
}

func (cmd *DirCommand) Execute(args []string) error {

	if len(args) == 0 {
		args = []string{"."}
	}

	results := cmd.GetResults(args)

	data, err := json.Marshal(results)
	if *cmd.loadOption {
		shell.PushText("application/json", string(data), err)
	}
	cmd.displayResults(results)
	return err
}

func (cmd *DirCommand) GetResults(args []string) []Result {
	var results = make([]Result, 0)

	for _, a := range args {

		var result = Result{}

		var folder = cmd.GetRootFolder(a)

		// Path Stat's successfully and is a folder Glob the contents
		{
			dirInfo, err := os.Stat(a)
			if err == nil {
				if dirInfo.IsDir() {
					a = a + "/*"
				}
			}
		}

		matches, err := filepath.Glob(a)
		if err != nil {
			result.Message = "Invalid parameter: " + a
			result.Folder = folder
			results = append(results, result)
			continue
		}

		result.Message = ""
		result.Folder = folder
		for _, path := range matches {
			file := cmd.getFileInfo(path)
			result.Files = append(result.Files, *file)
		}
		if len(result.Files) == 0 {
			result.Message = "No Files Found"
		}
		results = append(results, result)
	}
	return results
}

func (cmd *DirCommand) GetRootFolder(path string) string {
	if len(path) == 0 {
		path = "."
	}

	// A valid path to a file or directory was provided
	if info, err := os.Stat(path); err == nil {
		path, _ = filepath.Abs(path)
		if info.IsDir() {
			return path
		} else {
			return filepath.Dir(path)
		}
	}

	// A potential glob parameter was provided (or invalid file)
	matches, err := filepath.Glob(path)
	if err == nil {
		if len(matches) > 0 {
			fullPath, err := filepath.Abs(matches[0])
			if err == nil {
				path = filepath.Dir(fullPath)
			}
		}
		return path
	}
	return ""
}

func (cmd *DirCommand) displayResults(results []Result) {
	for _, result := range results {
		if len(result.Folder) > 0 {
			fmt.Fprintf(shell.ConsoleWriter(), "Folder: %s\n\n", result.Folder)
		} else {
			fmt.Fprintf(shell.ConsoleWriter(), "Invalid file or folder provided\n\n")
		}

		if len(result.Message) > 0 {
			fmt.Fprintf(shell.ConsoleWriter(), "%s\n", result.Message)
			continue
		}

		for _, f := range result.Files {
			name := f.Name
			if f.IsDir {
				name = name + "/"
			}

			fmt.Fprintf(shell.OutputWriter(),
				"%s  %-24s  %10d %s\n",
				f.Mode,
				f.Date.Format("2006-01-02 03:04:05PM"),
				f.Size,
				name)
		}
		fmt.Fprintln(shell.OutputWriter())
	}
}

func (cmd *DirCommand) getFileInfo(path string) *File {

	info, ioerr := os.Stat(path)
	if ioerr != nil {
		return nil
	}

	mode := info.Mode() & 01777
	dirChar := " "
	if info.IsDir() {
		dirChar = "D"
	}

	fullPath, _ := filepath.Abs(path)
	return &File{
		Name:      info.Name(),
		Date:      info.ModTime(),
		Mode:      fmt.Sprintf("%1s%#4o", dirChar, mode),
		Size:      info.Size(),
		IsDir:     info.IsDir(),
		FullPath:  fullPath,
		Reference: path,
	}
}
