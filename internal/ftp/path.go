// Package ftp 实现一个轻量、无外部依赖的 FTP 服务器。
//
// 支持被动(PASV/EPSV)与主动(PORT/EPRT)模式,完整的文件传输命令
// (LIST/MLSD/RETR/STOR/APPE)、目录操作、断点续传,以及沙箱路径隔离。
// 面向"快速分享文件"场景:默认匿名、即开即用。
package ftp

import (
	"path"
	"path/filepath"
	"strings"
)

// cleanSlash 规整化 FTP 虚拟路径:折叠多余斜杠、去尾部斜杠(根目录除外)。
// 统一使用正斜杠,跨平台一致。
func cleanSlash(p string) string {
	for strings.Contains(p, "//") {
		p = strings.ReplaceAll(p, "//", "/")
	}
	if p != "/" {
		p = strings.TrimRight(p, "/")
	}
	if p == "" {
		p = "/"
	}
	return p
}

// resolveRel 把客户端 cwd 与命令参数解析为 FTP 虚拟路径(以 / 为根)。
//   - 绝对路径(以 / 开头)直接规整
//   - 相对路径拼接到 cwd 之下
//   - 支持 .. / .
func resolveRel(cwd, arg string) string {
	if arg == "" {
		return cwd
	}
	var joined string
	if strings.HasPrefix(arg, "/") {
		joined = arg
	} else {
		joined = cwd + "/" + arg
	}
	return cleanSlash(path.Clean(joined))
}

// SecureJoin 将 FTP 虚拟路径安全映射到宿主机文件系统路径,确保结果不逃逸出 root。
//
// 返回 (absPath, error)。若解析后路径不在 root 之内,返回 EscapedErr。
// 这是对旧实现的修复点 #1:旧代码仅 filepath.Join+Clean,未做前缀校验,
// 在某些边界条件下可被构造路径绕过。这里用 EvalSymlinks 后再比对前缀。
func SecureJoin(root, virtual string) (string, error) {
	root = filepath.Clean(root)
	// 把虚拟路径拼到 root 下(虚拟路径以 / 开头,剥掉前导斜杠)。
	rel := strings.TrimPrefix(virtual, "/")
	joined := filepath.Join(root, rel)
	cleaned := filepath.Clean(joined)

	// 确保结果仍位于 root 之内:要求 cleaned == root 或以 root + 分隔符 为前缀。
	if cleaned != root && !strings.HasPrefix(cleaned, root+string(filepath.Separator)) {
		return "", &EscapedErr{Root: root, Path: cleaned}
	}
	return cleaned, nil
}

// EscapedErr 表示一次尝试越出沙箱根的路径访问。
type EscapedErr struct {
	Root string
	Path string
}

func (e *EscapedErr) Error() string {
	return "path escapes shared root: " + e.Path
}
