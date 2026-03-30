package cli

import (
	"fmt"
	"strings"

	"github.com/zy84338719/upftp/internal/config"
)

func (m Model) renderStatus() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, labelStyle.Render(" "+t(lang, "server_config")+":"))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "http_port")+":"))+" "+config.AppConfig.Port[1:])
	if config.AppConfig.EnableFTP {
		lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_port")+":"))+" "+config.AppConfig.FTPPort[1:])
	}
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "root_directory")+":"))+" "+truncatePath(config.AppConfig.Root, 36))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render(" "+t(lang, "features_status")+":"))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "http_auth")+":"))+" "+renderEnabled(config.AppConfig.HTTPAuth.Enabled))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "file_upload_feat")+":"))+" "+renderEnabled(config.AppConfig.Upload.Enabled))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_server_feat")+":"))+" "+renderEnabled(config.AppConfig.EnableFTP))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "https")+":"))+" "+renderEnabled(config.AppConfig.HTTPS.Enabled))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render(" "+t(lang, "statistics")+":"))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "files_available")+":"))+" "+fmt.Sprintf("%d", len(m.fileMap)))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "upload_max_size")+":"))+" "+formatSize(config.AppConfig.Upload.MaxSize))

	return "  " + dimStyle.Render(t(lang, "server_status_box")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderVersion() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "version_label")+":"))+" "+config.AppConfig.Version)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "git_commit")+":"))+" "+config.AppConfig.LastCommit)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "build_date")+":"))+" "+config.AppConfig.BuildDate)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "go_version")+":"))+" "+config.AppConfig.GoVersion)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "platform")+":"))+" "+config.AppConfig.Platform)
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "project_homepage")+":"))+" "+greenStyle.Render(truncateString(config.AppConfig.ProjectURL, 40)))
	lines = append(lines, "")

	return "  " + dimStyle.Render(t(lang, "version_info_box")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderAbout() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "")
	lines = append(lines, "  "+greenStyle.Render(config.AppConfig.ProjectName)+" "+t(lang, "about_desc"))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "key_features")+":"))
	lines = append(lines, "    ● "+t(lang, "feat_web_ui"))
	lines = append(lines, "    ● "+t(lang, "feat_mcp"))
	lines = append(lines, "    ● "+t(lang, "feat_upload"))
	lines = append(lines, "    ● "+t(lang, "feat_http_auth"))
	lines = append(lines, "    ● "+t(lang, "feat_ftp"))
	lines = append(lines, "    ● "+t(lang, "feat_qr"))
	lines = append(lines, "    ● "+t(lang, "feat_tree"))
	lines = append(lines, "    ● "+t(lang, "feat_lang"))
	lines = append(lines, "    ● "+t(lang, "feat_yaml"))
	lines = append(lines, "    ● "+t(lang, "feat_https"))
	lines = append(lines, "    ● "+t(lang, "feat_cli"))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "links")+":"))
	lines = append(lines, "    Homepage:  "+truncateString(config.AppConfig.ProjectURL, 47))
	lines = append(lines, "    GitHub:    https://github.com/zy84338719/upftp")
	lines = append(lines, "")
	lines = append(lines, "  "+dimStyle.Render(t(lang, "license")))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "build_info")+":"))
	lines = append(lines, "    "+t(lang, "version_label")+": "+config.AppConfig.Version)
	lines = append(lines, "    "+t(lang, "git_commit")+": "+config.AppConfig.LastCommit)
	lines = append(lines, "    "+t(lang, "build_date")+": "+config.AppConfig.BuildDate)
	lines = append(lines, "    "+t(lang, "platform")+": "+config.AppConfig.Platform)
	lines = append(lines, "")

	return "  " + dimStyle.Render(t(lang, "about_box")+" v"+config.AppConfig.Version) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderFTPInfo() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_server")+":"))+" "+m.serverIP)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_port")+":"))+" "+config.AppConfig.FTPPort[1:])
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "username")+":"))+" "+config.AppConfig.Username)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "password")+":"))+" "+config.AppConfig.Password)
	lines = append(lines, "")
	lines = append(lines, "  $ ftp "+m.serverIP)
	lines = append(lines, "  Name: "+config.AppConfig.Username)
	lines = append(lines, "  Password: "+config.AppConfig.Password)
	lines = append(lines, "  ftp> ls")
	lines = append(lines, "  ftp> get filename")
	lines = append(lines, "  ftp> put localfile")
	lines = append(lines, "")

	return "  " + dimStyle.Render(t(lang, "ftp_connection_info")) + "\n" + panel(w, strings.Join(lines, "\n"))
}
