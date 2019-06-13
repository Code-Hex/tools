package main

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("unexpected arguments more than 0")
	}
	filename := filepath.Base(args[0])
	toolsDir := filepath.Dir(args[0])
	if err := os.Chdir(toolsDir); err != nil {
		return err
	}
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "tools.go", src, 0)
	if err != nil {
		return err
	}

	var eg errgroup.Group
	for _, im := range f.Imports {
		path := strings.Trim(im.Path.Value, `"`)
		eg.Go(func() error {
			fmt.Println("+ go install", path)
			cmd := exec.Command("go", "install", path)
			cmd.Env = os.Environ()
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		})
	}

	return eg.Wait()
}
