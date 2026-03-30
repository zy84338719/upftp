package cli

import (
	"fmt"
	"strings"

	"github.com/zy84338719/upftp/pkg/conf"
)

func (m Model) renderBanner() string {
	lang := m.lang
	w := m.contentWidth()

	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render("~ upftp"))
	b.WriteString(" ")
	b.WriteString(valueStyle.Render("v" + conf.AppConfig.Version))
	b.WriteString("  ")
	b.WriteString(dimStyle.Render(t(lang, "tagline")))
	b.WriteString("\n")
	b.WriteString("  ")
	b.WriteString(dimStyle.Render(t(lang, "build") + ": " + conf.AppConfig.LastCommit + " @ " + conf.AppConfig.BuildDate))
	b.WriteString("\n\n")
	b.WriteString("  ")
	b.WriteString(separator(w))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(dimStyle.Render(fmt.Sprintf(" %-14s", t(lang, "http_server")+":")))
	b.WriteString(" ")
	b.WriteString(greenStyle.Render("http://" + m.serverIP + ":" + fmt.Sprintf("%d", conf.AppConfig.GetHTTPPort())))
	b.WriteString("\n")

	if conf.AppConfig.EnableFTP {
		b.WriteString("  ")
		b.WriteString(dimStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_server")+":")))
		b.WriteString(" ")
		b.WriteString(greenStyle.Render("ftp://" + m.serverIP + ":" + fmt.Sprintf("%d", conf.AppConfig.GetFTPPort())))
		b.WriteString("\n")
	}

	b.WriteString("  ")
	b.WriteString(dimStyle.Render(fmt.Sprintf(" %-14s", t(lang, "shared_path")+":")))
	b.WriteString(" ")
	b.WriteString(truncatePath(conf.AppConfig.Root, 47))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(dimStyle.Render(fmt.Sprintf(" %-14s", t(lang, "files_found")+":")))
	b.WriteString(" ")
	b.WriteString(whiteStyle.Render(fmt.Sprintf("%d", len(m.fileMap))))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(separator(w))
	b.WriteString("\n\n")

	b.WriteString("  ")
	b.WriteString(dimStyle.Render(t(lang, "features")))
	b.WriteString("\n")
	b.WriteString("    ")
	b.WriteString(greenStyle.Render("●"))
	b.WriteString(" ")
	b.WriteString(t(lang, "web_interface"))
	b.WriteString("  ")
	b.WriteString(greenStyle.Render("●"))
	b.WriteString(" ")
	b.WriteString(t(lang, "mcp_support"))
	b.WriteString("  ")
	b.WriteString(greenStyle.Render("●"))
	b.WriteString(" ")
	b.WriteString(t(lang, "file_upload"))
	b.WriteString("  ")
	b.WriteString(greenStyle.Render("●"))
	b.WriteString(" ")
	b.WriteString(t(lang, "qr_access"))
	b.WriteString("\n")

	if conf.AppConfig.Upload.Enabled {
		b.WriteString("    ")
		b.WriteString(yellowStyle.Render("●"))
		b.WriteString(" ")
		b.WriteString(t(lang, "upload_enabled"))
		b.WriteString("  ")
		b.WriteString(yellowStyle.Render("●"))
		b.WriteString(" ")
		b.WriteString(t(lang, "auth_enabled"))
		b.WriteString("\n")
	}

	b.WriteString("  ")
	b.WriteString(separator(w))
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderCompactHeader() string {
	w := m.contentWidth()
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render("~ upftp"))
	b.WriteString(" ")
	b.WriteString(dimStyle.Render("v" + conf.AppConfig.Version))
	b.WriteString("    ")
	b.WriteString(greenStyle.Render("http://" + m.serverIP + ":" + fmt.Sprintf("%d", conf.AppConfig.GetHTTPPort())))
	b.WriteString(" ")
	b.WriteString(dimStyle.Render("|"))
	b.WriteString(" ")
	b.WriteString(truncatePath(conf.AppConfig.Root, 30))
	b.WriteString(" ")
	b.WriteString(dimStyle.Render("|"))
	b.WriteString(" ")
	b.WriteString(fmt.Sprintf("%d files", len(m.fileMap)))
	b.WriteString("\n  ")
	b.WriteString(separator(w))
	return b.String()
}

func (m Model) renderMenu() string {
	lang := m.lang
	w := m.panelWidth()

	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(dimStyle.Render(t(lang, "commands")))
	b.WriteString("\n")

	var rows []string
	row1 := "  " + keyStyle.Render("[1]") + " " + t(lang, "search_files")
	row1 += "  " + keyStyle.Render("[2]") + " " + t(lang, "list_all_files")
	row1 += "  " + keyStyle.Render("[3]") + " " + t(lang, "download_examples")
	rows = append(rows, row1)

	row2 := "  " + keyStyle.Render("[4]") + " " + t(lang, "refresh_list")
	row2 += "  " + keyStyle.Render("[5]") + " " + t(lang, "server_status")
	row2 += "  " + keyStyle.Render("[6]") + " " + t(lang, "configuration")
	rows = append(rows, row2)

	row3 := "  " + keyStyle.Render("[7]") + " " + t(lang, "ftp_info")
	row3 += "  " + keyStyle.Render("[8]") + " " + t(lang, "about")
	rows = append(rows, row3)

	row4 := "  " + blueStyle.Render("[v]") + " " + t(lang, "version_info")
	row4 += "  " + redStyle.Render("[q]") + " " + t(lang, "quit_server")
	rows = append(rows, row4)

	row5 := "  " + blueStyle.Render("[l]") + " " + t(lang, "lang_switch")
	row5 += "  " + blueStyle.Render("[r]") + " " + "Reload Config"
	rows = append(rows, row5)

	b.WriteString(panel(w, strings.Join(rows, "\n")))
	return b.String()
}
