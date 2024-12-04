package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	selected int
	// err         error
	projectName string
	typing      bool
	showHelp    bool
}

func initialModel() model {
	return model{
		choices:  []string{"No-Fee", "Usage-based pricing"},
		selected: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Bubble Tea update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case msg.String() == "q", msg.String() == "esc":
			if m.showHelp {
				m.showHelp = false
				return m, tea.Quit
			} else if !m.typing {
				return m, tea.Quit
			}

		case msg.String() == "up":
			if !m.typing && !m.showHelp {
				if m.selected > 0 {
					m.selected--
				}
			}

		case msg.String() == "down":
			if !m.typing && !m.showHelp {
				if m.selected < len(m.choices)-1 {
					m.selected++
				}
			}

		case msg.String() == "enter":
			if !m.typing {
				m.typing = true
				return m, nil
			} else {
				m.projectName = strings.TrimSpace(m.projectName)
				if m.projectName != "" {
					return m, tea.Quit
				}
				m.typing = false
				return m, nil
			}
		case msg.String() == "backspace":
			if m.typing && len(m.projectName) > 0 {
				m.projectName = m.projectName[:len(m.projectName)-1]
				return m, nil
			}
			return m, nil

		default:
			if m.typing {
				m.projectName += msg.String()
				return m, nil
			} else if msg.String() == "k" && m.selected > 0 {
				m.selected--
				return m, nil
			} else if msg.String() == "j" && m.selected < len(m.choices)-1 {
				m.selected++
				return m, nil
			} else if msg.String() == "h" {
				m.showHelp = !m.showHelp
				return m, nil
			}
		}
	}
	return m, nil
}

// Bubble Tea view function
func (m model) View() string {
	if m.showHelp {
		return "Help Screen\n\n" +
			"Up/Down/j/k: Navigate\n" +
			"Enter: Select option or start typing project name\n" +
			"Backspace: Delete character (when typing)\n" +
			"q or esc: Quit or Exit Help\n" +
			"h: Toggle Help\n\n" +
			"Press q or esc to exit help."
	}
	if m.typing {
		return fmt.Sprintf("Enter project name: %s\n", m.projectName)
	}
	s := "Choose a project type:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.selected == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress up/down keys to navigate,\n" + "enter to select, h for help,\n" + "or q or esc to quit."
	return s
}

func downloadProject(projectName string, projectType int) error {
	owner := "smokeyblues"
	repo := ""
	branch := "main"

	if projectType == 0 {
		repo = "base-aws-no-fee-notes-app"
	} else if projectType == 1 {
		repo = "base-aws-notes-app"
	} else {
		return fmt.Errorf("invalid project type selected")
	}

	err := downloadAndExtract(owner, repo, branch, ".", projectName)

	if err != nil {
		return fmt.Errorf("error downloading and extracting project: %w", err)
	}
	return nil
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running Bubble Tea program: %v\n", err)
		os.Exit(1)
	}

	// Check the model to determine if the program should quit or if the TUI quit so the program logic can begin.
	if len(m.(model).projectName) > 0 {
		selectedModel := m.(model)
		fmt.Println("Creating project", selectedModel.projectName)

		// Create new directory
		err := os.MkdirAll(selectedModel.projectName, 0755)
		if err != nil {
			handleError(fmt.Errorf("error creating directory: %w", err))
		}

		// change directory to new directory
		err = os.Chdir(selectedModel.projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error changing directory: %v\n", err)
			os.Exit(1)
		}

		if err := downloadProject(selectedModel.projectName, selectedModel.selected); err != nil {
			handleError(err)
		}

		if err != nil {
			handleError(err)
		}
		fmt.Println("Project created successfully!")
	} else {
		fmt.Println("Exiting program. No Project created.")
	}

}

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
			// reader := io.TeeReader(rc, outfile)

			fileBytes, err := io.ReadAll(rc)
			if err != nil {
				return fmt.Errorf("error reading content: %w", err)

			}

			var modifiedContent string
			if strings.Contains(string(fileBytes), repo) {
				re := regexp.MustCompile(regexp.QuoteMeta(repo))

				modifiedContent = re.ReplaceAllLiteralString(string(fileBytes), newProjectName)
			} else {
				modifiedContent = string(fileBytes)
				// _, err = io.Copy(outfile, rc)
				// if err != nil {
				// 	return fmt.Errorf("error copying file content %s: %w", file.Name, err)
				// }
			}
			_, err = outfile.Write([]byte(modifiedContent))
			if err != nil {
				return fmt.Errorf("error writing modified content to file %s: %w", file.Name, err)
			}
			rc.Close()
		}
	}
	return nil
}
