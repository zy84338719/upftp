package cli

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zy84338719/upftp/pkg/conf"
)

type CLI struct {
	serverIP string
	fileMap  map[string]string
	lang     string
}

func NewCLI() *CLI {
	lang := conf.AppConfig.GetLanguage()
	return &CLI{
		fileMap: make(map[string]string),
		lang:    lang,
	}
}

func (c *CLI) SetServerIP(ip string) {
	c.serverIP = ip
}

func (c *CLI) Start(ctx context.Context, s chan os.Signal) {
	c.fileMap = scanDirectory(conf.AppConfig.Root, c.serverIP)

	if !isTerminal() {
		fmt.Printf("\n  %s\n", t(c.lang, "stdin_closed"))
		fmt.Printf("  %s\n", t(c.lang, "headless_note"))
		<-ctx.Done()
		return
	}

	quitRequested := false
	model := NewModel(c.serverIP, c.fileMap, c.lang, &quitRequested)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("TUI error: %v\n", err)
	}

	if quitRequested {
		s <- syscall.SIGQUIT
	}
}

func scanDirectory(rootDir, serverIP string) map[string]string {
	files := make(map[string]string)
	walkDir(rootDir, rootDir, serverIP, files)
	return files
}

func walkDir(currentDir, rootDir, serverIP string, files map[string]string) {
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		fullPath := path.Join(currentDir, entry.Name())
		if entry.IsDir() {
			walkDir(fullPath, rootDir, serverIP, files)
		} else {
			relDir := strings.TrimPrefix(currentDir, rootDir)
			relDir = strings.TrimPrefix(relDir, "/")

			var authPrefix string
			if conf.AppConfig.HTTPAuth.Enabled {
				authPrefix = conf.AppConfig.HTTPAuth.Username + ":" + conf.AppConfig.HTTPAuth.Password + "@"
			}

			hostPart := authPrefix + serverIP + ":" + fmt.Sprintf("%d", conf.AppConfig.GetHTTPPort())
			var downloadURL string
			if relDir != "" {
				downloadURL = fmt.Sprintf("http://%s/download/%s/%s", hostPart, relDir, entry.Name())
			} else {
				downloadURL = fmt.Sprintf("http://%s/download/%s", hostPart, entry.Name())
			}
			files[entry.Name()] = downloadURL
		}
	}
}

func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
