package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var imageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
	".bmp":  {},
	".tiff": {},
	".webp": {},
}

func main() {
	sourceDir := flag.String("source", `C:\Users\dell\Downloads\Maps-20260203T100443Z-3-001\Maps`, "directory where photos live")
	destDir := flag.String("dest", `C:\Users\dell\Downloads\locators images`, "root directory where the sorted folders are created")
	locationKeyword := flag.String("keyword", "locator", "keyword used to identify location-map images")
	locationFolder := flag.String("location-folder", "location-map", "subfolder name for files that match the keyword")
	mapFolder := flag.String("map-folder", "map", "subfolder name for files that do not match the keyword")
	dryRun := flag.Bool("dry-run", false, "just log what would be moved without touching files")
	verbose := flag.Bool("verbose", true, "print each move action")
	flag.Parse()

	absSource, err := filepath.Abs(*sourceDir)
	if err != nil {
		log.Fatalf("unable to resolve source directory %q: %v", *sourceDir, err)
	}

	absDest, err := filepath.Abs(*destDir)
	if err != nil {
		log.Fatalf("unable to resolve destination directory %q: %v", *destDir, err)
	}

	if err := os.MkdirAll(absDest, 0o755); err != nil {
		log.Fatalf("unable to create destination root %q: %v", absDest, err)
	}

	totalFiles := 0
	movedFiles := 0

	err = filepath.WalkDir(absSource, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if shouldSkipDir(path, absDest) {
				return fs.SkipDir
			}
			return nil
		}

		if !isImageFile(d.Name()) {
			return nil
		}

		totalFiles++

		targetSub := *mapFolder
		if strings.Contains(strings.ToLower(d.Name()), strings.ToLower(*locationKeyword)) {
			targetSub = *locationFolder
		}

		targetDir := filepath.Join(absDest, targetSub)
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			return fmt.Errorf("creating %q: %w", targetDir, err)
		}

		destPath := uniquePath(targetDir, d.Name())

		if *dryRun {
			log.Printf("[dry run] %s -> %s", path, destPath)
			return nil
		}

		if *verbose {
			log.Printf("moving %s -> %s", path, destPath)
		}

		if err := os.Rename(path, destPath); err != nil {
			return fmt.Errorf("moving %q to %q: %w", path, destPath, err)
		}

		movedFiles++
		return nil
	})

	if err != nil {
		log.Fatalf("failed to walk %q: %v", absSource, err)
	}

	log.Printf("scanned %d image(s), relocated %d", totalFiles, movedFiles)
}

func shouldSkipDir(path, absDest string) bool {
	rel, err := filepath.Rel(absDest, path)
	if err == nil && (rel == "." || !strings.HasPrefix(rel, "..")) {
		return true
	}
	return false
}

func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	_, ok := imageExtensions[ext]
	return ok
}

func uniquePath(dir, name string) string {
	target := filepath.Join(dir, name)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return target
	}

	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	for i := 1; ; i++ {
		alt := filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, i, ext))
		if _, err := os.Stat(alt); os.IsNotExist(err) {
			return alt
		}
	}
}
