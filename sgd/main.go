package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	var (
		outDir string
		wg     sync.WaitGroup
	)

	flag.StringVar(&outDir, "o", "./docs", "output dir for docs")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Please specifiy a go package")
	}

	basePkg := args[0]
	basePath := filepath.Join(os.Getenv("GOPATH"), "src", basePkg)

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		wg.Add(1)
		defer wg.Done()

		// Ignore files and hidden folders
		if !info.IsDir() || strings.Contains(path, "/.") {
			return nil
		}

		// fix for base pkg docs
		if path == basePath {
			path += "/"
		}

		// package name w/o base pkg path prefix
		pkg := strings.Replace(path, basePath, "", 1)

		// output dir
		outPath, err := filepath.Abs(filepath.Join(outDir, pkg))
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}

		// create all folders for path
		if err := os.MkdirAll(outPath, 0700); err != nil {
			os.Stderr.WriteString(err.Error())
		}

		// create pkg doc file
		f, err := os.Create(filepath.Join(outDir, pkg, "/index.html"))
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
		defer f.Close()

		// write html top wrapper
		f.WriteString("<!doctype html>\n<html>\n<head>\n<title>Docs</title>\n</head>\n<body>\n")

		// generate godoc html and write to file
		cmd := exec.Command("godoc", "-html", filepath.Join(basePkg, pkg))
		cmd.Stderr = os.Stderr
		cmd.Stdout = f
		cmd.Run()

		// write html bottom
		f.WriteString("\n</body></html>")

		return nil
	})

	wg.Wait()
}
