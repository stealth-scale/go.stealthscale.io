# go.stealthscale.io

This repository is the vanity import host of the Go modules of Stealth Scale. Its pages serve the
`go-import` meta tags that resolve `go get go.stealthscale.io/<module>` to the sources at
[github.com/stealth-scale](https://github.com/stealth-scale).

GitHub Pages serves the site at <https://go.stealthscale.io/>.

## Structure

```text
.
├── index.html           generated index of the modules
├── 404.html             generated fallback page
├── <module>/index.html  generated page of each module, with its meta tags
├── assets/
│   ├── site.css         stylesheet without a build step
│   └── site.js          copies a `go get` command on click
└── _gen/                site generator, which Jekyll leaves out of the site
    ├── main.go
    ├── data.go          list of modules
    ├── go.mod
    └── templates/
```

The generator in `_gen/` renders every HTML file of the site. The repository commits the rendered
files, so GitHub Pages serves them without a build.

## Add or edit a module

1. Edit the `modules` slice in `_gen/data.go`.
2. Run `cd _gen && go run .`.
3. Commit `_gen/data.go` with the regenerated HTML.

The generator formats its output with `npx prettier`. Pass `-no-fmt` when npx is not installed.
Run the generator inside `_gen/`, because the default `-out` is `..` relative to the working
directory.

```go
{
    Name:        "cli",
    Repo:        "go-cli",          // optional; defaults to Name
    Description: "The package cli runs command-line programs ...",
    Public:      true,              // listed under Public modules
    Langs:       []string{"go"},    // optional; defaults to ["go"]
},
```

Set `Repo` when the name of the GitHub repository differs from the import path. The module
`go.stealthscale.io/cli` is in `github.com/stealth-scale/go-cli`.

## Nested modules

Each nested module with a `go.mod` of its own has a page at its import path. Set `Dir` when the
`go.mod` is in a directory other than the suffix of the import path:

```go
Subs: []Sub{
    {Name: "toml", Description: "...", Public: true},
    {Name: "lang/go", Dir: "cli-lang-go", Description: "...", Public: true},
},
```

Without `Dir`, the `go.mod` of `go.stealthscale.io/cli/toml` is in `toml/`. The page serves the
`go-import` tag of the parent:

```html
<meta name="go-import"
      content="go.stealthscale.io/cli git https://github.com/stealth-scale/go-cli" />
```

Go takes the path after the prefix of the tag as the directory of the module in the repository.
The version tags of the module are `toml/vX.Y.Z`.

With `Dir`, the page serves the import path of the module as the prefix and `Dir` as the fourth
field. Go reads the fourth field from Go 1.25 on.

```html
<meta name="go-import"
      content="go.stealthscale.io/cli/lang/go git https://github.com/stealth-scale/go-cli cli-lang-go" />
```

The version tags of the module are `cli-lang-go/vX.Y.Z`. The `go-source` tag has the same prefix.
pkg.go.dev rejects a page whose two tags have different prefixes.

> **Note:** Go ignores the `replace` directives of a dependency. A nested module that requires
> its parent at `v0.0.0-00010101000000-000000000000` with `replace ... => ../` builds only inside
> its own repository. A published nested module must require a tagged version of its parent. Put
> the local `replace` in a `go.work` instead.

## Private modules

Set `Public: false`. The page of a private module shows the `GOPRIVATE` setting that `go get`
needs:

```sh
go env -w GOPRIVATE=go.stealthscale.io/*
```

`go get` also needs git credentials with access to the `stealth-scale` organization.

## Local preview

Serve the repository and open <http://localhost:8000/>:

```sh
python3 -m http.server 8000
```

Check the meta tags of a module:

```sh
curl -sL 'http://localhost:8000/cli/toml/?go-get=1' | grep go-import
```

## No build on GitHub Pages

GitHub Pages serves the committed files of the branch `main` without CI or a deploy pipeline. The
site consists of plain HTML, CSS and about 50 lines of JavaScript.
