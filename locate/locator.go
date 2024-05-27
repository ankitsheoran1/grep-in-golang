package locate

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"strings"

	"github.com/fatih/color"
)

type Locator struct {
	BaseDir string
	Options OptionConfig
}

type OptionConfig struct {
	Verbose, Hidden bool
}

type AnalyzedFile struct {
	FilePath  string
	Content   *os.File
	Locations []Location
	Ok        bool
}

type Location struct {
	LineNo   int
	Contents string
}

func NewLocator(dir string) *Locator {
	return &Locator{
		BaseDir: dir,
	}

}

func (f *AnalyzedFile) GetInfo() {
	filepath := color.YellowString(f.FilePath)
	fmt.Printf("\n%s\n", filepath)
	for _, loc := range f.Locations {
		lineNo := color.GreenString(fmt.Sprintf("%v", loc.LineNo))
		fmt.Printf("%s:%s\n", lineNo, loc.Contents)
	}

}

func findFiles(l *Locator, fs []fs.DirEntry) ([]string, []string) {
	var dirs []string
	var files []string
	for _, file := range fs {
		switch file.IsDir() {
		case true:
			if !l.Options.Hidden && file.Name()[0] == '.' {
				continue 
			}
			dirs = append(dirs, file.Name())
		case false:
			files = append(files, file.Name())
		}
	}
	return dirs, files

}

// Dig recursively searches through a directory and it's files to find the given string "text"
func (l *Locator) Dig(text string) {
	fs, err := os.ReadDir(l.BaseDir)
	if err != nil {
		fmt.Println(err)
	}

	dirs, files := findFiles(l, fs)
	wg := sync.WaitGroup{}
	for _, fileName := range files {
		wg.Add(1)
		go runAnalyze(l, fileName, text, &wg)
	}
	for _, dir := range dirs {
		wg.Add(1)

		loc := NewLocator(l.BaseDir + "/" + dir)

		go runDig(loc, text, &wg)
	}

	wg.Wait()

}

func runAnalyze(l *Locator, fileName, text string, wg *sync.WaitGroup) {
	file := l.Analyze(fileName, text)
	defer wg.Done()

	if !file.Ok {
		return
	}
	file.GetInfo()
}

func runDig(locator *Locator, text string, wg *sync.WaitGroup) {
	defer wg.Done()

	locator.Dig(text)
}

// Use Analyze to open and scan the given file, fileName, and assert whether it contains the given string, text.
// If it does it is accumulated to a slice of AnalyzedFile and later returned
func (l *Locator) Analyze(fileName string, text string) AnalyzedFile {
	file, err := os.Open(l.BaseDir + "/" + fileName)
	if err != nil {
		fmt.Printf("Path \"%s\" is not valid\n", fileName)
		return AnalyzedFile{}
	}
	scanner := bufio.NewScanner(file)
	lineNo := 1
	var locations []Location
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		line := scanner.Text()
		exists := strings.Contains(line, text)
		if exists {
			resString := strings.Replace(line, text, color.RedString(text), -1)
			loc := Location{LineNo: lineNo, Contents: resString}
			locations = append(locations, loc)
		}
		lineNo++


	}

	analyzedFile := AnalyzedFile{FilePath: l.BaseDir + "/" + fileName, Content: file, Locations: locations}
	if len(analyzedFile.Locations) > 0 {
		analyzedFile.Ok = true
	}

	return analyzedFile


}
