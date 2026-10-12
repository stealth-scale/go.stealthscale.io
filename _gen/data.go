package main

import (
	"path"
	"strings"
)

const repoOrg = "stealth-scale"

// Module describes a top-level module at go.stealthscale.io/<Name>.
type Module struct {
	Name string
	// Repo names the GitHub repository when its name differs from Name, as
	// go-cli does for cli.
	Repo        string
	Description string
	Public      bool
	// Langs lists the language tags that the site shows for the module. An
	// empty list shows go.
	Langs []string
	Subs  []Sub
}

// Sub describes a nested module with its own go.mod in the repo of its parent,
// at go.stealthscale.io/<Parent>/<Name>.
type Sub struct {
	Name string
	// Dir names the directory of the go.mod inside the parent's repo when it
	// differs from Name.
	Dir         string
	Description string
	Public      bool
	Langs       []string

	parent *Module // main sets parent before it renders a page.
}

func (m Module) ImportPath() string { return "go.stealthscale.io/" + m.Name }
func (m Module) GoGet() string      { return "go get " + m.ImportPath() }

func (m Module) RepoName() string {
	if m.Repo == "" {
		return m.Name
	}
	return m.Repo
}

func (m Module) RepoURL() string { return "https://github.com/" + repoOrg + "/" + m.RepoName() }

func (m Module) Languages() []string {
	if len(m.Langs) == 0 {
		return []string{"go"}
	}
	return m.Langs
}

func (m Module) GoSource() string { return goSource(m.ImportPath(), m.RepoURL(), "") }

// goSource returns the content of a go-source meta tag for the import path
// prefix. The tag maps the path after prefix to the same path below dir, in the
// repo at repoURL.
func goSource(prefix, repoURL, dir string) string {
	base := path.Join("main", dir)
	return strings.Join([]string{
		prefix,
		"    " + repoURL,
		"    " + repoURL + "/tree/" + base + "{/dir}",
		"    " + repoURL + "/blob/" + base + "{/dir}/{file}#L{line}",
	}, "\n")
}

func (s Sub) Path() string       { return s.parent.Name + "/" + s.Name }
func (s Sub) ImportPath() string { return "go.stealthscale.io/" + s.Path() }
func (s Sub) GoGet() string      { return "go get " + s.ImportPath() }
func (s Sub) ParentName() string { return s.parent.Name }
func (s Sub) ParentImport() string {
	return s.parent.ImportPath()
}

// repoDir returns the directory of the module's go.mod inside the parent's
// repo.
func (s Sub) repoDir() string {
	if s.Dir == "" {
		return s.Name
	}
	return s.Dir
}

// RepoURL returns the URL of the module's directory on GitHub.
func (s Sub) RepoURL() string { return s.parent.RepoURL() + "/tree/main/" + s.repoDir() }

// GoImport returns the content of the go-import meta tag at the module's path.
//
// A module whose directory is its path suffix resolves through the parent's
// tag, because Go takes the path after the tag's prefix as the module's
// directory inside the repo. The tag of a module in another directory states
// the module's own path, and its fourth field states the directory. Go reads
// the fourth field from Go 1.25 on, and the module's version tags then start
// with the directory.
func (s Sub) GoImport() string {
	if s.repoDir() == s.Name {
		return s.parent.ImportPath() + " git " + s.parent.RepoURL()
	}
	return s.ImportPath() + " git " + s.parent.RepoURL() + " " + s.repoDir()
}

// GoSource returns the content of the go-source meta tag at the module's path.
// Its prefix is the prefix of [Sub.GoImport], because pkg.go.dev rejects a page
// whose two tags state different prefixes.
func (s Sub) GoSource() string {
	if s.repoDir() == s.Name {
		return s.parent.GoSource()
	}
	return goSource(s.ImportPath(), s.parent.RepoURL(), s.repoDir())
}

func (s Sub) Languages() []string {
	if len(s.Langs) == 0 {
		return []string{"go"}
	}
	return s.Langs
}

// modules is the source of truth for the site. Edit this slice and re-run
// `cd _gen && go run .`.
var modules = []Module{
	{
		Name:        "cli",
		Repo:        "go-cli",
		Description: "The package cli runs command-line programs whose options are the fields of Go structs. It derives the flag, the environment variable and the config key of each field. It converts every source to the type of the field and applies the sources in one order. A handler receives its options as typed values.",
		Public:      true,
		Subs: []Sub{
			{
				Name:        "toml",
				Description: "The package toml reads the TOML config files of the apps of go.stealthscale.io/cli. The package cli does not import a package outside the standard library, so the TOML parser is in a module of its own.",
				Public:      true,
			},
			{
				Name:        "yaml",
				Description: "The package yaml reads the YAML config files of the apps of go.stealthscale.io/cli. The package cli does not import a package outside the standard library, so the YAML parser is in a module of its own.",
				Public:      true,
			},
		},
	},
}
