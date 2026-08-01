// Package core 提供 FTP 与 HTTP 共享的业务逻辑:会话管理、文件操作、认证。
//
// 它是整个应用的中枢——HTTP API、CLI、TUI 都通过 core 操作 FTP 实例和文件,
// 而不是各自直接碰底层。
package core

import (
	"crypto/subtle"
	"net/http"
)

// Auth 封装认证策略:匿名或单一用户名密码。
type Auth struct {
	Anonymous bool
	User      string
	Pass      string
}

// Enabled 表示是否启用凭据认证(非匿名)。
func (a *Auth) Enabled() bool { return !a.Anonymous && a.User != "" }

// Check 校验用户名密码是否匹配。匿名模式恒为 true。
func (a *Auth) Check(user, pass string) bool {
	if !a.Enabled() {
		return true
	}
	uMatch := subtle.ConstantTimeCompare([]byte(user), []byte(a.User)) == 1
	pMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(a.Pass)) == 1
	return uMatch && pMatch
}

// HTTPMiddleware 返回一个 net/http 中间件,启用时强制 HTTP Basic Auth。
func (a *Auth) HTTPMiddleware(next http.Handler) http.Handler {
	if !a.Enabled() {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || !a.Check(u, p) {
			w.Header().Set("WWW-Authenticate", `Basic realm="upftp"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
