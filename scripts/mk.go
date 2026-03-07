package main

import (
	"cmp"
	"fmt"
	"iter"
	"log"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func chunked[T any](size int, list iter.Seq[T]) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		chunk := []T{}
		for item := range list {
			if len(chunk) == size {
				if !yield(chunk) {
					return
				}
				chunk = []T{}
			}
			chunk = append(chunk, item)
		}
		if len(chunk) != 0 {
			yield(chunk)
		}
	}
}

var images = map[string][]string{
	"shader_rgb":      {"shader", "-s", "shader_examples/rgb_coloring.glsl"},
	"zcurve":          {"zcurve"},
	"verticallines":   {"verticallines"},
	"sharpen":         {"sharpen"},
	"blur":            {"blur"},
	"quadtree":        {"quadtree"},
	"weakblur":        {"weakblur"},
	"median":          {"median"},
	"hilbertdarken":   {"hilbertdarken"},
	"hilbert":         {"hilbert"},
	"emboss":          {"emboss"},
	"horizontallines": {"horizontallines"},
	"edgedetect2":     {"edgedetect2"},
	"edgeenhance":     {"edgeenhance"},
	"edgedetect1":     {"edgedetect1"},
	"cluster":         {"cluster", "-n", "7"},
}

var imgs_dir = "img/static"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func execute_cmd(cmd string, args ...string) string {
	output, err := exec.Command(cmd, args...).Output()
	check(err)
	return string(output)
}

func ls(dir string) iter.Seq[string] {
	files, err := os.ReadDir(dir)
	check(err)
	return func(yield func(string) bool) {
		for _, file := range files {
			if !yield(file.Name()) {
				return
			}
		}
	}
}

func Map[T any, U any](list []T, f func(T) U) []U {
	res := make([]U, len(list))
	for i, item := range list {
		res[i] = f(item)
	}
	return res
}

func Sorted[T cmp.Ordered](seq iter.Seq[T]) iter.Seq[T] {
	slice := slices.Collect(seq)
	slices.Sort(slice)
	return slices.Values(slice)
}

func usage() string {
	output := execute_cmd("go", "run", "cmd/fimgs/main.go", "--help")
	// help returns commands in random order, so we sort them first
	lines := strings.Split(output, "\n")
	i := slices.Index(lines, "COMMANDS:")
	before, kekafter := lines[:i+1], lines[i+1:]
	j := slices.Index(kekafter, "")
	kek, after := kekafter[:j], kekafter[j:]
	slices.Sort(kek)
	return strings.Join(append(before, append(kek, after...)...), "\n")
}

var app = cli.App{
	Name:  "mk",
	Usage: "commands runner",
	Commands: []*cli.Command{
		{
			Name:  "imgs",
			Usage: "update example imgs from orig.png",
			Action: func(*cli.Context) error {
				origFilename := filepath.Join(imgs_dir, "orig.png")
				maxLen := slices.Max(Map(slices.Collect(maps.Keys(images)), func(s string) int { return len(s) }))
				for destination := range Sorted(maps.Keys(images)) {
					args := images[destination]
					// fmt.Println("building", destination, args)
					start := time.Now()
					imageFilename := strings.TrimSpace(execute_cmd("go", append([]string{"run", "cmd/fimgs/main.go", "-i", origFilename}, args...)...))
					end := time.Now()
					check(os.Rename(imageFilename, filepath.Join(imgs_dir, destination+".png")))
					fmt.Printf("built %s in %s\n", destination+strings.Repeat(" ", maxLen-len(destination)), end.Sub(start).String())
				}
				return nil
			},
		},
		{
			Name:  "readme",
			Usage: "compile readme file",
			Action: func(*cli.Context) error {
				w, err := os.OpenFile("README.md", os.O_WRONLY|os.O_CREATE, 0o664)
				check(err)
				defer w.Close()

				// markdown routines
				H1 := func(text string) {
					fmt.Fprintln(w, "# "+text)
					fmt.Fprintln(w)
				}
				H2 := func(text string) {
					fmt.Fprintln(w, "## "+text)
					fmt.Fprintln(w)
				}
				Code := func(lang, code string) {
					fmt.Fprintln(w, "```"+lang)
					fmt.Fprintln(w, code)
					fmt.Fprintln(w, "```")
					fmt.Fprintln(w)
				}
				Table := func(headers []string, rows iter.Seq[[]string]) {
					fmt.Fprintln(w, "|"+strings.Join(headers, "|")+"|")
					fmt.Fprintln(w, strings.Repeat("|-", len(headers))+"|")
					for row := range rows {
						fmt.Fprintln(w, "|"+strings.Join(row, "|")+"|")
					}
				}

				// content := strings.Join(func() []string {
				H1("fimgs - image filters tool")
				H2("Install")
				Code("bash", "go install github.com/rprtr258/fimgs/cmd/fimgs@latest")
				H2("Usage")
				Code("php", usage())
				H2("Examples")
				a3 := chunked(3, Sorted(func(yield func(string) bool) {
					for img := range ls(imgs_dir) {
						if !strings.HasSuffix(img, ".png") {
							continue
						}
						if !yield(img) {
							return
						}
					}
				}))
				Table([]string{"", "", ""}, func(yield func([]string) bool) {
					for imgs := range a3 {
						if !yield(Map(imgs, func(img string) string { return "![](./img/static/" + img + ")" })) || !yield(Map(imgs, func(img string) string { return img })) {
							return
						}
					}
				})

				return nil
			},
		},
	},
	HideHelpCommand: true,
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
