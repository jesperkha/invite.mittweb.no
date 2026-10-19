package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"singel-sondag/config"
	"time"
)

const sessionCookie = "invite_admin"

func sessionValue(cfg *config.Config) string {
	secret := cfg.SessionSecret
	if secret == "" {
		secret = cfg.AdminPassword
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("admin"))
	return hex.EncodeToString(mac.Sum(nil))
}

func validSession(r *http.Request, cfg *config.Config) bool {
	cookie, err := r.Cookie(sessionCookie)
	return err == nil && hmac.Equal([]byte(cookie.Value), []byte(sessionValue(cfg)))
}

func AdminAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !validSession(r, cfg) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AdminLoginHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("password") != cfg.AdminPassword || cfg.AdminPassword == "" {
			http.Error(w, "invalid password", http.StatusUnauthorized)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: sessionValue(cfg), Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: r.TLS != nil, MaxAge: 86400})
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: "", Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, MaxAge: -1, Expires: time.Unix(0, 0)})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
