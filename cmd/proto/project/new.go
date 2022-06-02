// Package project
// @author: xs
// @date: 2022/5/31
// @Description: project
package project

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"golang.org/x/mod/modfile"
	"os"
	"path"
)

var repoAddIgnores = []string{
	".git", ".github", "api", "README.md", "LICENSE", "go.mod", "go.sum", "third_party",
}

// Project is a project template.
type Project struct {
	Name string
	Path string
}

// New new a project from remote repo.
func (p *Project) New(ctx context.Context, dir string, layout string, branch string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		os.RemoveAll(to)
	}
	fmt.Printf("🚀 Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)
	repo := NewRepo(layout, branch)
	if err := repo.CopyTo(ctx, to, p.Path, []string{".git", ".github"}); err != nil {
		return err
	}
	e := os.Rename(
		path.Join(to, "cmd", "ghub"),
		path.Join(to, "cmd", p.Name),
	)
	if e != nil {
		return e
	}
	Tree(to, dir)

	fmt.Printf("\n🍺 Project creation succeeded %s\n", color.GreenString(p.Name))
	fmt.Print("💻 Use the following command to start the project 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ go generate ./..."))
	fmt.Println(color.WhiteString("$ go build -o ./bin/ ./... "))
	fmt.Println(color.WhiteString("$ ./bin/%s -conf ./configs\n", p.Name))
	fmt.Println("			🤝 Thanks for using Kratos")
	fmt.Println("	📚 Tutorial: https://go-kratos.dev/docs/getting-started/start")
	return nil
}

func (p *Project) Add(ctx context.Context, dir string, layout string, branch string, mod string) error {
	to := path.Join(dir, p.Path)

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		os.RemoveAll(to)
	}

	fmt.Printf("🚀 Add service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)

	repo := NewRepo(layout, branch)

	if err := repo.CopyToV2(ctx, to, path.Join(mod, p.Path), repoAddIgnores, []string{path.Join(p.Path, "api"), "api"}); err != nil {
		return err
	}

	e := os.Rename(
		path.Join(to, "cmd", "ghub"),
		path.Join(to, "cmd", p.Name),
	)
	if e != nil {
		return e
	}

	Tree(to, dir)

	fmt.Printf("\n🍺 Repository creation succeeded %s\n", color.GreenString(p.Name))
	fmt.Print("💻 Use the following command to add a project 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ go generate ./..."))
	fmt.Println(color.WhiteString("$ go build -o ./bin/ ./... "))
	fmt.Println(color.WhiteString("$ ./bin/%s -conf ./configs\n", p.Name))
	fmt.Println("			🤝 Thanks for using Kratos")
	fmt.Println("	📚 Tutorial: https://go-kratos.dev/docs/getting-started/start")
	return nil
}

// ModulePath returns go module path.
func ModulePath(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}
