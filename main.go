package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2).Foreground(lipgloss.Color("222"))

type Script struct {
	name, command string
}

type CommandExecuted struct{}

func (s Script) Title() string       { return s.name }
func (s Script) Description() string { return s.command }
func (s Script) FilterValue() string { return s.name }

var lockfilesToPackageManagers = map[string]string{
	"pnpm-lock.yaml":    "pnpm",
	"package-lock.json": "npm",
	"bun.lockb":         "bun",
	"yarn.lock":         "yarn",
}
var packageManager string

func detectPackageManager() string {
	var lockfiles []string

	cwd, cwdErr := os.Getwd()

	if cwdErr != nil {
		fmt.Println(docStyle.Margin(0, 2).Render("Unable to get current directory"))
		os.Exit(1)
	}

	files, readDirErr := os.ReadDir(cwd)

	if readDirErr != nil {
		fmt.Println(docStyle.Margin(0, 2).Render("Unable to read contents of current directory"))
		os.Exit(1)
	}

	for _, file := range files {
		if _, setInMap := lockfilesToPackageManagers[file.Name()]; setInMap {
			lockfiles = append(lockfiles, file.Name())
		}

		if len(lockfiles) > 1 {
			fmt.Println(docStyle.UnsetMargins().Render("Multiple lockfiles found in", cwd))
			os.Exit(1)
			break
		}

		if !file.IsDir() {
			switch file.Name() {
			case "pnpm-lock.yaml":
				packageManager = "pnpm"
			case "package-lock.json":
				packageManager = "npm"
			case "bun.lockb":
				packageManager = "bun"
			case "yarn.lock":
				packageManager = "yarn"
			}
		}
	}

	if packageManager == "" {
		packageManager = "npm"
	}
	return packageManager
}

func installDependencies(packageManager string) {
	_, err := os.Stat("node_modules")

	if err != nil {
		command := exec.Command(packageManager, "install")

		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		fmt.Printf("Installing %v packages\n", packageManager)
		err := command.Run()
		if err != nil {
			fmt.Printf("Error running command: %v\n", err)
		}
	}
}

func runScript(packageManager, scriptName string) tea.Cmd {
	command := exec.Command(packageManager, "run", scriptName)
	return tea.ExecProcess(command, func(err error) tea.Msg {
		if err != nil {
			return tea.Quit()
		}
		return CommandExecuted{}
	})
}

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("np-run")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "enter" {
			script, _ := m.list.SelectedItem().(Script)
			return m, runScript(detectPackageManager(), script.name)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case CommandExecuted:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	packageJsonContent, readFileErr := os.ReadFile("package.json")

	if readFileErr != nil {
		fmt.Println(docStyle.UnsetMargins().Render("package.json not found"))
		os.Exit(1)
	}

	detectPackageManager()
	installDependencies(packageManager)

	var parsedJson map[string]interface{}

	parseErr := json.Unmarshal(packageJsonContent, &parsedJson)

	if parseErr != nil {
		fmt.Println(docStyle.UnsetMargins().Render(parseErr.Error()))
		os.Exit(1)
	}

	var items []list.Item

	for scriptKey, scriptValue := range parsedJson["scripts"].(map[string]interface{}) {
		items = append(items, Script{name: scriptKey, command: scriptValue.(string)})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Scripts to Run"
	m.list.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fff")).
		Background(lipgloss.Color("#bc54c4")).
		Padding(0, 1)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
