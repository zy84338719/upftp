package cli

import (
	"fmt"
	"strings"

	"github.com/zy84338719/upftp/internal/config"
)

func (m Model) renderSearch() string {
	lang := m.lang
	w := m.panelWidth()

	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(dimStyle.Render(t(lang, "search_results")))
	b.WriteString("\n")

	if m.inputActive {
		b.WriteString("  ")
		b.WriteString(m.input.View())
		b.WriteString("\n")
	} else {
		term := m.input.Value()
		found := false
		count := 0

		var lines []string
		for filename, url := range m.fileMap {
			if term == "" || strings.Contains(strings.ToLower(filename), strings.ToLower(term)) {
				found = true
				count++
				line := greenStyle.Render("●") + " " + truncateString(filename, 56)
				line += "\n" + dimStyle.Render(truncateString(url, 60))
				lines = append(lines, line)
				if count >= 20 {
					lines = append(lines, dimStyle.Render(t(lang, "showing_first")))
					break
				}
			}
		}

		if !found {
			lines = append(lines, redStyle.Render(t(lang, "no_match")))
		}

		sep := "\n" + dimStyle.Render(stringsRepeat("─", w-2)) + "\n"
		b.WriteString(panel(w, strings.Join(lines, sep)))
		b.WriteString("\n  ")
		b.WriteString(dimStyle.Render(fmt.Sprintf(t(lang, "found_count"), count)))
	}

	return b.String()
}

func (m Model) renderFiles() string {
	lang := m.lang
	w := m.panelWidth()

	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(dimStyle.Render(t(lang, "all_files")))
	b.WriteString("\n")

	count := 0
	var lines []string
	for filename, url := range m.fileMap {
		count++
		line := greenStyle.Render("●") + " " + truncateString(filename, 56)
		line += "\n" + dimStyle.Render(truncateString(url, 60))
		lines = append(lines, line)
		if count >= 20 {
			lines = append(lines, dimStyle.Render(fmt.Sprintf(t(lang, "and_more"), len(m.fileMap)-20)))
			break
		}
	}

	sep := "\n" + dimStyle.Render(stringsRepeat("─", w-2)) + "\n"
	b.WriteString(panel(w, strings.Join(lines, sep)))
	b.WriteString("\n  ")
	b.WriteString(dimStyle.Render(fmt.Sprintf(t(lang, "total_files"), len(m.fileMap))))
	return b.String()
}

func (m Model) renderExamples() string {
	lang := m.lang
	w := m.panelWidth()

	if len(m.fileMap) == 0 {
		return "  " + redStyle.Render(t(lang, "no_files_example"))
	}

	var exampleFile, exampleURL string
	for filename, url := range m.fileMap {
		exampleFile = filename
		exampleURL = url
		break
	}

	var lines []string
	lines = append(lines, labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "example_file")+":"))+" "+truncateString(exampleFile, 40))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render(" "+t(lang, "browser")+":"))
	lines = append(lines, "  "+greenStyle.Render("http://"+m.serverIP+config.AppConfig.Port))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render(" "+t(lang, "cli_tools")+":"))
	lines = append(lines, "  curl -O \""+truncateString(exampleURL, 52))
	lines = append(lines, "  wget \""+truncateString(exampleURL, 54))

	if config.AppConfig.EnableFTP {
		lines = append(lines, "")
		lines = append(lines, labelStyle.Render(" "+t(lang, "ftp_client")+":"))
		lines = append(lines, "  ftp "+m.serverIP)
		lines = append(lines, "  Username: "+config.AppConfig.Username)
		lines = append(lines, "  Password: "+config.AppConfig.Password)
	}

	lines = append(lines, "")
	lines = append(lines, labelStyle.Render(" "+t(lang, "mcp_integration")+":"))
	lines = append(lines, "  \"command\": \"upftp\",")
	lines = append(lines, "  \"args\": [\"-enable-mcp\"]")

	return "  " + dimStyle.Render(t(lang, "dl_examples")) + "\n" + panel(w, strings.Join(lines, "\n"))
}
