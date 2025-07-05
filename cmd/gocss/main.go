package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/su3h7am/gocss/pkg/core"
	"github.com/su3h7am/gocss/pkg/extractor"
	"github.com/su3h7am/gocss/pkg/preset"
)

func main() {
	// Define command-line flags
	inputPatterns := flag.String("input", "./**/*.html", "Glob pattern for input files (e.g., ./**/*.html, ./src/**/*.templ)")
	outputFile := flag.String("output", "./gocss.css", "Output CSS file path")
	watchMode := flag.Bool("watch", false, "Enable watch mode to rebuild CSS on file changes")
	flag.Parse()

	// Configure GOCSS
	cfg := &core.Config{
		Presets: []core.Preset{
			preset.NewWind(),
		},
		Layers: map[string]int{
			"base":       0,
			"components": 1,
			"utilities":  2,
		},
		Extractors: []core.Extractor{
			&extractor.ExtractorSplit{},
			&extractor.TemplExtractor{},
		},
	}

	resolvedConfig := core.NewResolvedConfig(cfg)
	generator := core.NewGenerator(resolvedConfig)

	// Build function
	build := func() {
		fmt.Println("Building CSS...")
		filesToProcess := map[string]string{
		"test.html": `<div class="m-4 p-8 block text-red-500 bg-blue-500 text-lg font-bold w-full h-screen hover:text-green-500 sm:p-16 btn btn-red flex items-center justify-center grid grid-cols-2 gap-4"></div><span class="text-white rounded"></span>`,
	}

		// Find files matching glob patterns
		matches, err := filepath.Glob(*inputPatterns)
		if err != nil {
			log.Fatalf("Error matching glob pattern: %v", err)
		}

		for _, match := range matches {
			content, err := ioutil.ReadFile(match)
			if err != nil {
				log.Printf("Error reading file %s: %v", match, err)
				continue
			}
			filesToProcess[match] = string(content)
		}

		css, err := generator.Generate(filesToProcess)
		if err != nil {
			log.Fatalf("Error generating CSS: %v", err)
		}

		err = ioutil.WriteFile(*outputFile, []byte(css), 0644)
		if err != nil {
			log.Fatalf("Error writing CSS to file: %v", err)
		}
		fmt.Printf("CSS built successfully to %s\n", *outputFile)
	}

	// Initial build
	build()

	// Watch mode
	if *watchMode {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		// Add directories to watch
		// This is a simplified approach; a real implementation might need to watch parent directories
		// or use a more sophisticated file watching library.
		matches, err := filepath.Glob(*inputPatterns)
		if err != nil {
			log.Fatalf("Error matching glob pattern for watch: %v", err)
		}
		watchedDirs := make(map[string]bool)
		for _, match := range matches {
			dir := filepath.Dir(match)
			if !watchedDirs[dir] {
				err = watcher.Add(dir)
				if err != nil {
					log.Fatal(err)
				}
				watchedDirs[dir] = true
			}
		}

		fmt.Println("Watching for file changes...")
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("Detected change in %s, rebuilding...", event.Name)
					build()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}
}