package rest

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type (
	options struct {
		Repo        string
		ProjectName string
	}

	sourceInfo struct {
		path    string
		content string
	}
	sourceInfos []sourceInfo
)

var (
	sources sourceInfos
)

func newInitCommand() *cobra.Command {
	opts := &options{}

	cmd := &cobra.Command{
		Use:   "init /path/to/project",
		Short: "Create a basic project structure",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			opts.Valid(cmd, args)
			opts.Run(cmd, args)
		},
	}

	opts.Set(cmd.PersistentFlags())

	return cmd
}

func register(s ...sourceInfo) {
	sources = append(sources, s...)
}

func (opts *options) Set(flags *pflag.FlagSet) {
	flags.StringVar(&opts.ProjectName, "project", "boxapp", "Project Name")
	flags.StringVar(&opts.Repo, "repo", "", "Git repository, eg: github.com.cn/boxgo")
}

func (opts *options) Valid(cmd *cobra.Command, args []string) {

}

func (opts *options) Run(cmd *cobra.Command, args []string) {
	for _, s := range sources {
		opts.writeFile(s.path, s.content)
	}
}

func (opts *options) writeFile(path, content string) {
	if _, err := os.Lstat(opts.ProjectName); os.IsNotExist(err) {
		os.Mkdir(opts.ProjectName, 0775)
	}

	os.MkdirAll(filepath.Join(opts.ProjectName, filepath.Dir(path)), 0775)

	content = strings.TrimLeft(content, "\n")
	content = strings.TrimRight(content, "		")

	tmpl, err := template.New(path).Parse(content)
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, opts); err != nil {
		log.Fatalln(err)
	}

	if err := ioutil.WriteFile(opts.filePath(path), buf.Bytes(), 0664); err != nil {
		log.Fatalln(err)
	}
}

func (opts *options) filePath(file ...string) string {
	elems := []string{opts.ProjectName}
	elems = append(elems, file...)

	return filepath.Join(elems...)
}
