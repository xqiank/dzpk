package security

import (
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ---------- 限流：按 IP 令牌桶 ----------

type bucket struct {
	tokens float64
	last   time.Time
}

type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64
	burst   float64
}

func NewLimiter(perSec int) *Limiter {
	if perSec <= 0 {
		perSec = 30
	}
	return &Limiter{buckets: map[string]*bucket{}, rate: float64(perSec), burst: float64(perSec) * 3}
}

func (l *Limiter) Allow(key string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.burst, last: now}
		l.buckets[key] = b
	}
	elapsed := now.Sub(b.last).Seconds()
	b.tokens += elapsed * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func (l *Limiter) Sweep(maxIdle time.Duration) {
	for {
		time.Sleep(time.Minute)
		cut := time.Now().Add(-maxIdle)
		l.mu.Lock()
		for k, b := range l.buckets {
			if b.last.Before(cut) {
				delete(l.buckets, k)
			}
		}
		l.mu.Unlock()
	}
}

// ---------- 输入校验 ----------

var unsafeCtrl = regexp.MustCompile(`[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]`)
var spaceRun = regexp.MustCompile(`\s{2,}`)

// CleanText 统一清洗：去控制字符、折叠空白、裁剪长度。
func CleanText(s string, max int) string {
	s = unsafeCtrl.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\u200b", "")
	s = strings.ReplaceAll(s, "\u200c", "")
	s = strings.ReplaceAll(s, "\u200d", "")
	s = spaceRun.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) > max {
		s = string(r[:max])
	}
	return s
}

var nickRe = regexp.MustCompile(`^[\p{Han}A-Za-z0-9_\-\.\+\*\s]{1,16}$`)

func ValidNickname(s string) bool { return nickRe.MatchString(s) }

var userRe = regexp.MustCompile(`^[A-Za-z0-9_]{3,20}$`)

func ValidUsername(s string) bool { return userRe.MatchString(s) }

func ValidPassword(s string) bool {
	if len(s) < 6 || len(s) > 64 {
		return false
	}
	return true
}

// LooksSuspicious 识别常见注入/探测形态，用于审计而非阻断（阻断由参数化查询保证）。
var injectPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(\bunion\b\s+\bselect\b)`),
	regexp.MustCompile(`(?i)(\bor\b\s+1\s*=\s*1)`),
	regexp.MustCompile(`(?i)(\bdrop\b\s+\btable\b)`),
	regexp.MustCompile(`(?i)(\bsleep\s*\(|benchmark\s*\()`),
	regexp.MustCompile(`(?i)(<script|javascript:|onerror\s*=)`),
	regexp.MustCompile(`(?i)(\.\./\.\./|/etc/passwd|\.\.%2f)`),
	regexp.MustCompile(`(?i)(\$where|\$ne\s*:|__proto__|constructor\s*\[)`),
}

func LooksSuspicious(s string) bool {
	if len(s) > 4096 {
		return true
	}
	for _, re := range injectPatterns {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}

// ---------- HTTP 安全头 ----------

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "SAMEORIGIN")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("X-XSS-Protection", "1; mode=block")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(self), camera=()")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		h.Set("Content-Security-Policy",
			"default-src 'self'; img-src 'self' data: blob: https:; media-src 'self' blob: data:; "+
				"script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; "+
				"connect-src 'self' ws: wss:; font-src 'self' data:; frame-ancestors 'self'")
		if strings.HasPrefix(r.URL.Path, "/api/") {
			h.Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}

func ClientIP(r *http.Request) string {
	if v := r.Header.Get("CF-Connecting-IP"); v != "" {
		return v
	}
	if v := r.Header.Get("X-Real-IP"); v != "" {
		return v
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// CORS 仅放白名单 Origin，且不回通配符
func CORS(origins []string) func(http.Handler) http.Handler {
	set := map[string]bool{}
	for _, o := range origins {
		set[strings.TrimSpace(o)] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			o := r.Header.Get("Origin")
			if o != "" && set[o] {
				w.Header().Set("Access-Control-Allow-Origin", o)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With")
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
				w.Header().Set("Access-Control-Max-Age", "600")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func OriginAllowed(origins []string, o string) bool {
	for _, x := range origins {
		if strings.TrimSpace(x) == o {
			return true
		}
	}
	return false
}
