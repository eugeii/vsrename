package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// removeExt removes the extension, specified in argument ext, from a
// filename. It returns the filename without the extension.
func removeExt(filename, ext string) string {
	return strings.TrimRight(filename, "."+ext)
}

// addExt adds an extension, specified in argument ext, to a filename. It
// returns the filename with the extension.
func addExt(filenameWithoutExt, ext string) string {
	return filenameWithoutExt + "." + ext
}

func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

var (
	subExt     string
	subMatcher string
	subPath    string
	vidExt     string
	vidMatcher string
	vidPath    string
	isWrite    bool
)

func init() {
	// Handle command line options.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0]))
		fmt.Printf("  vsrename (for full options, see below)\n")
		fmt.Printf("    [--subext=<subtitle extension>]\n")
		fmt.Printf("    [--vidext=<video extension>]\n")
		fmt.Printf("    [--subregex=<regex pattern to identify episode in subtitle files>]\n")
		fmt.Printf("    [--vidregex=<regex pattern to identify episode in video files>]\n\n")

		fmt.Printf("  Examples:\n")
		fmt.Printf("    (show renames without actually renaming)\n")
		fmt.Printf("    vsrename -vext='mkv' -subregex='x([0-9]+)' -vidregex='E([0-9]+)'\n\n")
		fmt.Printf("    (show renames and actually rename)\n")
		fmt.Printf("    vsrename -vext='mkv' -subregex='x([0-9]+)' -vidregex='E([0-9]+) -w'\n\n")

		flag.PrintDefaults()
		os.Exit(0)
	}

	// Subtitles.
	flag.StringVar(&subExt, "subext", "srt", "The extension of the subtitle files (e.g. 'srt')")
	flag.StringVar(&vidExt, "sext", "srt", "The extension of the subtitle files (e.g. 'srt') (shorthand)")
	flag.StringVar(&subMatcher, "subregex", "x([0-9]+)", "The regex to identify episode of each subtitle file (as a regex group)")

	// Videos.
	flag.StringVar(&vidExt, "vidext", "mkv", "The extension of the video files (e.g. 'mp4')")
	flag.StringVar(&vidExt, "vext", "mkv", "The extension of the video files (e.g. 'mp4') (shorthand)")
	flag.StringVar(&vidMatcher, "vidregex", "E([0-9]+)", "The regex to identify episode of each video file (as a regex group)")

	// Paths.
	flag.StringVar(&vidPath, "vidpath", ".", "The path to the location of the video and subtitle files")
	flag.StringVar(&subPath, "subpath", "./subs", "The path to the location of the video and subtitle files")

	// Commit the rename.
	flag.BoolVar(&isWrite, "write", false, "Actually perform the rename")
	flag.BoolVar(&isWrite, "w", false, "Actually perform the rename (shorthand)")
}

func main() {
	// Command line options.
	flag.Parse()

	// Regex (required arguments).
	if subMatcher == "" || vidMatcher == "" {
		fmt.Printf("Regex pattern for subtitle and videos required. Aborting.\n")
		return
	}
	subRegex := regexp.MustCompile(subMatcher)
	vidRegex := regexp.MustCompile(vidMatcher)

	// Find subtitle and video files.
	subFiles, err := filepath.Glob(filepath.Join(subPath, "*."+subExt))
	if err != nil {
		fmt.Printf("Find subtitle files: %v\n", err)
		return
	}

	vidFiles, err := filepath.Glob(filepath.Join(vidPath, "*."+vidExt))
	if err != nil {
		fmt.Printf("Find video files: %v\n", err)
		return
	}

	fmt.Printf("Found total %v video files (*.%v) and %v subtitle files (*.%v).\n", len(vidFiles), vidExt, len(subFiles), subExt)
	if len(vidFiles) <= 0 {
		fmt.Printf("No video files found. Aborting.\n")
		return
	}

	// Store map of all subtitles.
	subtitles := map[string]string{}
	for _, sf := range subFiles {
		// Get video episode.
		m := (subRegex).FindStringSubmatch(sf)
		if m == nil || len(m) < 2 {
			fmt.Printf("  [X] Ignoring subtitle file '%v' (does not match regex).\n", sf)
			continue
		}
		episode := m[1]
		subtitles[episode] = sf
	}
	if len(subtitles) <= 0 {
		fmt.Printf("No subtitles matching regex found. Aborting.\n")
		return
	}

	// Rename files.
	numRenamed := 0
	for _, vf := range vidFiles {
		// Get video episode.
		m := (vidRegex).FindStringSubmatch(vf)
		if m == nil || len(m) < 2 {
			fmt.Printf("  [X] '%v' -> Skipping (episode not found matching regex).\n", vf)
			continue
		}
		episode := m[1]

		// Find matching subtitle file.
		if sf, ok := subtitles[episode]; ok {
			// Video file takes name of subtitle file (retaining video extension).
			newVf := filepath.Join(
				vidPath, // Place renamed files in same directory as original video files
				filepath.Base(addExt(removeExt(sf, subExt), vidExt)), // Rename video files to subtitle files
			)
			fmt.Printf("  [*] '%v' -> '%v'\n", vf, newVf)
			if isWrite {
				numRenamed++

				// Rename video file.
				os.Rename(vf, newVf)

				// Copy subtitle file to where video file is.
				newSf := filepath.Join(vidPath, filepath.Base(sf))
				err = copyFile(sf, newSf)
				if err != nil {
					fmt.Printf("  [!] Copy subtitle '%v' to '%v': %v\n", sf, newSf, err)
					continue
				}
			}

			continue
		}

		// Could not find matching subtitle file.
		fmt.Printf("  [X] '%v' -> No subtitle file found. Skipping.\n", vf)
	}

	// Show number of files renamed.
	fmt.Println()
	switch numRenamed {
	case 0:
		fmt.Printf("No files renamed.\n")
	default:
		fmt.Printf("%v files renamed.\n", numRenamed)

	}
}
