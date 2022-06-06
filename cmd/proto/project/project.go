// Package project
// @author: xs
// @date: 2022/5/30
// @Description: project
package project

import (
	"context"
	"embed"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"time"
)

var stubsFS embed.FS

// CmdServer represents the new command.
var CmdServer = &cobra.Command{
	Use:   "new",
	Short: "Create a service template",
	Long:  "Create a service project using the repository template. Example: proto new object",
	Run:   run,
}

var (
	repoURL string
	branch  string
	timeout string
	nomod   bool
)

func init() {
	//if repoURL = os.Getenv("KRATOS_LAYOUT_REPO"); repoURL == "" {
	//	repoURL = "https://github.com/go-kratos/kratos-layout.git"
	//}
	repoURL = "https://github.com/china-xs/ghub.git"
	timeout = "6000s"
	CmdServer.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	CmdServer.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
	CmdServer.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
	CmdServer.Flags().BoolVarP(&nomod, "nomod", "", nomod, "retain go mod")
}

// NewProject is a project template.
type NewProject struct {
	Name string
	Path string
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is project name ?",
			Help:    "Created project name.",
		}
		err = survey.AskOne(prompt, &name)
		if err != nil || name == "" {
			return
		}
	} else {
		name = args[0]
	}
	p := &Project{Name: path.Base(name), Path: name}
	done := make(chan error, 1)
	go func() {
		if !nomod {
			done <- p.New(ctx, wd, repoURL, branch)
			return
		}
		if _, e := os.Stat(path.Join(wd, "go.mod")); os.IsNotExist(e) {
			done <- fmt.Errorf("🚫 go.mod don't exists in %s", wd)
			return
		}

		mod, e := ModulePath(path.Join(wd, "go.mod"))
		if e != nil {
			panic(e)
		}
		done <- p.Add(ctx, wd, repoURL, branch, mod)
	}()
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, "\033[31mERROR: project creation timed out\033[m\n")
			return
		}
		fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to create project(%s)\033[m\n", ctx.Err().Error())
	case err = <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: Failed to create project(%s)\033[m\n", err.Error())
		}
	}
}
