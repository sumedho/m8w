package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed static/*
var embededFiles embed.FS

type appEnv struct {
	port      string
	info      bool
	local_dir bool
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("m8w", flag.ContinueOnError)
	fl.StringVar(
		&app.port, "p", "8000", "server port number (default 8000)",
	)
	fl.BoolVar(
		&app.info, "v", false, "print information",
	)
	fl.BoolVar(
		&app.local_dir, "local", false, "Use local static directory",
	)
	if err := fl.Parse(args); err != nil {
		return err
	}
	return nil
}

func (app *appEnv) getFileSystem() http.FileSystem {
	if app.local_dir {
		log.Print("using live mode")
		return http.FS(os.DirFS("static"))
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(embededFiles, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func (app *appEnv) run() error {
	log.Print("Running M8WebDisplay on port: ", app.port)
	http.Handle("/", http.FileServer(app.getFileSystem()))
	http.ListenAndServe(":"+string(app.port), nil)
	return nil
}

func CLI(args []string) int {
	var app appEnv
	err := app.fromArgs(args)
	if err != nil {
		return 2
	}
	if err = app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(CLI(os.Args[1:]))
}
