package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"

	// "net/url"
	"os"
	"path/filepath"
	"strings"
	// tea "github.com/charmbracelet/bubbletea"
)

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func downloadAndExtract(owner, repo, branch, targetDir, newProjectName string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", owner, repo, branch)

	resp, err := http.Get(url)
	// Handle errors
	if err != nil {
		return fmt.Errorf("error downloading from Github: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	// Handle errors
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	// Handle errors
	if err != nil {
		return fmt.Errorf("error unzipping package: %w", err)
	}

	for _, file := range zipReader.File {

		// Remove the top-level directory from paths since Github zips include a folder named after the commit SH
		extractedPath := strings.SplitN(file.Name, "/", 2)[1]
		fpath := filepath.Join(targetDir, extractedPath)

		if file.FileInfo().IsDir() {
			// Make directory
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				// Log the error, but continue processing other files
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", fpath, err)
				continue
			}
		} else {

			// Create a new file at the extracted path
			outfile, err := os.Create(fpath)
			// Handle errors
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", fpath, err)
				continue
			}
			// fmt.Fprintln("writing file %s", fpath)
			defer outfile.Close()

			rc, err := file.Open()
			// Handle errors
			if err != nil {
				return fmt.Errorf("error opening file %s: %w", file.Name, err)
			}
			defer rc.Close()

			fmt.Printf("Writing file: %s\n", file.Name)

			// buf := new(bytes.Buffer)
			reader := io.TeeReader(rc, outfile)

			fileBytes, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("error reading from TeeReader: %w", err)

			}

			if strings.Contains(string(fileBytes), repo) {
				fmt.Printf("Found '%s' in file: %s\n", repo, file.Name)
				re := regexp.MustCompile(regexp.QuoteMeta(repo))

				modifiedContent := re.ReplaceAllLiteralString(string(fileBytes), newProjectName)
				_, err = outfile.WriteAt([]byte(modifiedContent), 0)
				if err != nil {
					return fmt.Errorf("error writing modified content to file %s: %w", file.Name, err)
				}
			} else {
				_, err = io.Copy(outfile, rc)
				if err != nil {
					return fmt.Errorf("error copying file content %s: %w", file.Name, err)
				}
			}
			rc.Close()
		}
	}
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("What would you like to call your project: ")

	// Read until newline
	newDir, _ := reader.ReadString('\n')

	// Remove the trailing newline
	newDir = strings.TrimSpace(newDir)

	fmt.Println("Your new project shall be named", newDir)

	// Create the directory for new project
	err := os.MkdirAll(newDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1) //Indicate an error with a non-zero exit code
	}

	// Change directory to the new directory
	err = os.Chdir(newDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error changing directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Current working directory:", getCurrentDirectory())

	// Go get project and write it to this directory
	err = downloadAndExtract("smokeyblues", "aws-sstv4-notes", "main", ".", newDir)
	if err != nil {
		fmt.Println("error in the download and extraction process.")
		handleError(err)
	}

	fmt.Println("Project created successfully.")

}

func getCurrentDirectory() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v", err)
		os.Exit(1)
	}
	return currentDir
}
