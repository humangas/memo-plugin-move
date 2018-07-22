package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// ExitOK is 0
	ExitOK = 0
	// ExitError is 1
	ExitError = 1
)

func usage() {
	msg := `memo move <directory>
Move the selected file to the specified directory`
	fmt.Println(msg)
}

func selectFile(cfg *config) (*bytes.Buffer, error) {
	f, err := os.Open(cfg.MemoDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var files []string
	files, err = f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	files = filterMarkdown(files)

	var buf bytes.Buffer
	var cmd *exec.Cmd
	if cfg.SelectCmd == "fzf" {
		// TODO: Extract this setting to the setting file?
		option := "--multi --cycle --bind=ctrl-u:half-page-up,ctrl-d:half-page-down"
		cmd = exec.Command(cfg.SelectCmd, (strings.Split(option, " "))[0:]...)
	} else {
		cmd = exec.Command(cfg.SelectCmd)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = &buf
	cmd.Stdin = strings.NewReader(strings.Join(files, "\n"))

	if err := cmd.Run(); err != nil {
		// If the file is not selected, then it exit 0
		if len(buf.String()) == 0 {
			// os.Exit(ExitOK)
			return nil, nil
		}
		return nil, err
	}

	return &buf, nil
}

func cmd(arg string) error {
	var cfg config
	if err := cfg.load(); err != nil {
		return err
	}

	if arg == "" {
		return errors.New("Usage: memo move <directory>")
	}

	if !fileExists(filepath.Join(cfg.MemoDir, arg)) {
		return fmt.Errorf("Error: Can not find path: %s", arg)
	}

	buf, err := selectFile(&cfg)
	if err != nil {
		return err
	}
	if buf == nil {
		return nil
	}

	memoDir, err := filepath.EvalSymlinks(cfg.MemoDir)
	if err != nil {
		return err
	}
	oldPath := filepath.Join(memoDir, strings.TrimSpace(buf.String()))
	newPath := filepath.Join(memoDir, arg, strings.TrimSpace(buf.String()))

	return os.Rename(oldPath, newPath)
}

func returnCode(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return ExitError
	}
	return ExitOK
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func filterMarkdown(files []string) []string {
	var newfiles []string
	for _, file := range files {
		if strings.HasSuffix(file, ".md") {
			newfiles = append(newfiles, file)
		}
	}

	return newfiles
}

func main() {
	arg := ""
	if len(os.Args) > 1 {
		if os.Args[1] == "-usage" {
			usage()
			os.Exit(ExitOK)
		}
		arg = os.Args[1]
	}
	os.Exit(returnCode(cmd(arg)))
}
