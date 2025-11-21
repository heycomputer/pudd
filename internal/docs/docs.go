package docs

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"github.com/heycomputer/pudding/internal/parser"
)
// BrowserOpener is a function type for opening URLs in a browser
type BrowserOpener func(url string) error

// CommandRunner is a function type for running external commands
type CommandRunner func(name string, args ...string) (output []byte, err error)

// Default implementations
var (
	defaultBrowserOpener BrowserOpener = openBrowser
	defaultCommandRunner CommandRunner = runCommand
)

// FetchAndOpen fetches documentation for a dependency and opens it in the browser
func FetchAndOpen(dep *parser.Dependency, projectType parser.ProjectType, keywords string) error {
	return fetchAndOpenWithFuncs(dep, projectType,  keywords, defaultCommandRunner, defaultBrowserOpener)
}

// fetchAndOpenWithFuncs allows dependency injection for testing
func fetchAndOpenWithFuncs(dep *parser.Dependency, projectType parser.ProjectType, keywords string, cmdRunner CommandRunner, browserOpener BrowserOpener) error {
	switch projectType {
	case parser.ProjectTypeElixir:
		return fetchAndOpenHexDocs(dep, keywords, cmdRunner, browserOpener)
	case parser.ProjectTypeRuby:
		return fetchAndOpenGemDocs(dep, keywords, cmdRunner, browserOpener)
	default:
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
}

func fetchAndOpenHexDocs(dep *parser.Dependency, keywords string, cmdRunner CommandRunner, browserOpener BrowserOpener) error {
	// Use mix hex.docs offline to fetch and open docs
	cmdParams := []string{"mix", "hex.docs", "fetch", dep.Name}
	if dep.Version != "" {
		cmdParams = append(cmdParams, dep.Version)
	}

	fetchDocsOutput, err := cmdRunner(cmdParams[0], cmdParams[1:]...)
	
	if err != nil {
		return fmt.Errorf("failed to fetch docs for %s: %w", dep.Name, err)
	}

	// extract the path from fetchDocsOutput
	// assuming the output contains the path in a known format
	// e.g., "Docs fetched to /path/to/docs"
	
	docPath := extractDocPath(string(fetchDocsOutput))
	if docPath == "" {
		return fmt.Errorf("failed to extract docs path for %s from %s", dep.Name, string(fetchDocsOutput))
	}

	// Construct the local URL to the documentation
	hexDocsURL := fmt.Sprintf("file://%s/", docPath)

	// Append search query if provided
	if keywords != "" {
		hexDocsURL = fmt.Sprintf("%ssearch.html?q=%s", hexDocsURL, url.QueryEscape(keywords))
	} else {
		hexDocsURL = fmt.Sprintf("%sindex.html", hexDocsURL)
	}

	// Open the documentation in browser using shell expansion
	if err := browserOpener(hexDocsURL); err != nil {
		return fmt.Errorf("failed to open docs for %s: %w", dep.Name, err)
	}

	return nil
}

func extractDocPath(output string) string {
	// e.g. "Docs fetched:  /path/to/docs\n"
	// from the first forward slash to the end of the line
	idx := strings.Index(output, "/")
	if idx == -1 {
		return ""
	}
	start := idx
	end := strings.Index(output[start:], "\n")
	if end == -1 {
		end = len(output)
	} else {
		end += start
	}
	return strings.TrimSpace(output[start:end])
}

func fetchAndOpenGemDocs(dep *parser.Dependency, keywords string, cmdRunner CommandRunner, browserOpener BrowserOpener) error {
	_, err := cmdRunner("rdoc", dep.Name, "--rdoc", "--version", dep.Version)
	if err != nil {
		return fmt.Errorf("failed to generate rdoc for %s: %w", dep.Name, err)
	}

	// run command to get gem env home and assign to variable
	gemEnvOutput, err := cmdRunner("sh", "-c", "gem env home")
	if err != nil {
		return fmt.Errorf("failed to get gem env home: %w", err)
	}
	// convert gemEnvOutput to string and strip whitespace/newlines
	gemEnvHome := string(gemEnvOutput)
	gemEnvHome = strings.TrimSpace(gemEnvHome)

	// Get the path to the generated documentation
	// open $(gem env home)/doc/GEM_NAME-GEM_VERSION/rdoc/table_of_contents.html
	gemDocTocUrl := fmt.Sprintf("file://%s/doc/%s-%s/rdoc/", gemEnvHome, dep.Name, dep.Version)

	// Append search query if provided
	if keywords != "" {
		gemDocTocUrl = fmt.Sprintf("%sindex.html?q=%s", gemDocTocUrl, url.QueryEscape(keywords))
	} else {
		gemDocTocUrl = fmt.Sprintf("%s%s", gemDocTocUrl, "table_of_contents.html")
	}

	// Open the documentation in browser using shell expansion
	if err := browserOpener(gemDocTocUrl); err != nil {
		return fmt.Errorf("failed to open rdoc for %s: %w", dep.Name, err)
	}

	return nil
}

// runCommand executes an external command
func runCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	
	if err != nil {
		return nil, fmt.Errorf("failed to run %s: %w (output: %s)", name, err, string(output))
	}
	return output, nil
}

func openBrowser(url string) error {
	// Try to open URL in default browser
	// macOS: open, Linux: xdg-open, Windows: start
	cmd := exec.Command("open", url)
	err := cmd.Start()
	if err != nil {
		// Try xdg-open for Linux
		cmd = exec.Command("xdg-open", url)
		err = cmd.Start()
	}
	return err
}
