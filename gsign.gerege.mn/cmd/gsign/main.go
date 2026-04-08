package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.Info("starting gsign.gerege.mn")

	port := envOrDefault("PORT", "8445")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"gsign.gerege.mn"}`))
	})
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40"><rect width="40" height="40" rx="10" fill="#7c3aed"/><text x="50%" y="54%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="700" font-size="11">G</text></svg>`))
	})

	addr := ":" + port
	srv := &http.Server{
		Addr:         addr,
		Handler:      logMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		slog.Info("listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexPage)
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "ip", r.RemoteAddr, "latency_ms", time.Since(start).Milliseconds())
	})
}

const indexPage = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>G-Sign Gateway — gsign.gerege.mn</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#0b1120;color:#e2e8f0;min-height:100vh}
nav{display:flex;align-items:center;justify-content:space-between;padding:16px 32px;border-bottom:1px solid rgba(255,255,255,.06)}
.nav-logo{display:flex;align-items:center;gap:10px;font-weight:700;font-size:16px;color:#fff}
.nav-logo svg{width:32px;height:32px}
.nav-links{display:flex;gap:24px}
.nav-links a{color:#94a3b8;font-size:13px;text-decoration:none;font-weight:500}
.nav-links a:hover{color:#fff}
.hero{text-align:center;padding:80px 24px 48px}
.badge{display:inline-flex;align-items:center;gap:6px;padding:6px 16px;background:rgba(124,58,237,.1);border:1px solid rgba(124,58,237,.2);border-radius:24px;font-size:12px;color:#a78bfa;font-weight:500;margin-bottom:32px}
.hero h1{font-size:48px;font-weight:800;line-height:1.1;margin-bottom:20px;color:#fff}
.hero h1 span{background:linear-gradient(135deg,#8b5cf6,#7c3aed);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.hero p{max-width:640px;margin:0 auto 16px;color:#94a3b8;font-size:16px;line-height:1.7}
.hero .sub{font-size:13px;color:#64748b;margin-bottom:40px}
.cta-group{display:flex;flex-wrap:wrap;justify-content:center;gap:16px;margin-bottom:16px}
.sign-btn{display:inline-flex;align-items:center;gap:10px;padding:18px 40px;background:linear-gradient(135deg,#7c3aed,#6d28d9);color:#fff;font-weight:700;font-size:17px;border-radius:14px;text-decoration:none;transition:all .2s;box-shadow:0 4px 20px rgba(124,58,237,.3)}
.sign-btn:hover{transform:translateY(-2px);box-shadow:0 8px 30px rgba(124,58,237,.4)}
.sign-btn svg{width:22px;height:22px}
.sign-hint{margin-top:16px;font-size:12px;color:#475569}
.features{display:grid;grid-template-columns:repeat(auto-fit,minmax(260px,1fr));gap:20px;max-width:900px;margin:60px auto 0;padding:0 24px}
.feature{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:16px;padding:28px}
.feature-icon{width:44px;height:44px;border-radius:12px;display:flex;align-items:center;justify-content:center;margin-bottom:16px;font-size:20px}
.feature-icon.purple{background:rgba(124,58,237,.1);color:#a78bfa}
.feature-icon.green{background:rgba(22,163,74,.1);color:#4ade80}
.feature-icon.amber{background:rgba(245,158,11,.1);color:#fbbf24}
.feature h3{font-size:15px;font-weight:700;color:#fff;margin-bottom:8px}
.feature p{font-size:13px;color:#94a3b8;line-height:1.6;margin:0}
.info{max-width:900px;margin:60px auto 0;padding:0 24px}
.info-card{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:16px;padding:32px}
.info-card h2{font-size:22px;font-weight:700;color:#fff;margin-bottom:20px}
.info-card h3{font-size:15px;font-weight:600;color:#e2e8f0;margin:20px 0 8px}
.info-card p{color:#94a3b8;font-size:14px;line-height:1.8;margin-bottom:8px}
.info-card code{background:rgba(255,255,255,.06);padding:2px 6px;border-radius:4px;font-size:13px;color:#a78bfa}
.info-card ul{color:#94a3b8;font-size:14px;line-height:2;padding-left:24px}
.info-card a{color:#a78bfa;text-decoration:none}
.info-card a:hover{text-decoration:underline}
.steps{display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:16px;margin:24px 0}
.step{background:rgba(124,58,237,.06);border:1px solid rgba(124,58,237,.12);border-radius:12px;padding:20px;text-align:center}
.step-num{width:32px;height:32px;background:#7c3aed;border-radius:8px;display:flex;align-items:center;justify-content:center;font-weight:700;color:#fff;font-size:14px;margin:0 auto 12px}
.step h4{font-size:13px;font-weight:600;color:#fff;margin-bottom:4px}
.step p{font-size:12px;color:#94a3b8;margin:0;line-height:1.5}
.footer{text-align:center;padding:48px 24px 32px;font-size:12px;color:#475569}
.footer a{color:#a78bfa;text-decoration:none}
@media(max-width:640px){.hero h1{font-size:32px}.features{grid-template-columns:1fr}.steps{grid-template-columns:1fr}}
</style>
</head>
<body>
<nav>
  <div class="nav-logo">
    <svg viewBox="0 0 32 32" fill="none"><rect width="32" height="32" rx="8" fill="#7c3aed"/><text x="50%" y="55%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="800" font-size="13">G</text></svg>
    G-Sign Gateway
  </div>
  <div class="nav-links">
    <a href="/">Overview</a>
    <a href="https://dan.gerege.mn/docs">DAN Docs</a>
  </div>
</nav>

<div class="hero">
  <div class="badge">Gerege Systems LLC</div>
  <h1>G-Sign <span>Mongolia</span></h1>
  <p>Монгол Улсын клауд тоон гарын үсгийн (Cloud Digital Signature) gateway. УБЕГ-ийн GSign системтэй холбогдон баримт бичигт цахим гарын үсэг зурна.</p>
  <p class="sub">sso.gov.mn / burtgel.gov.mn-р дамжуулан тоон гарын үсэг баталгаажуулна.</p>
  <div class="cta-group">
    <a href="https://sso.gov.mn/login_signature" class="sign-btn">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 19l-7-7 1.41-1.41L11 15.17V2h2v13.17l4.59-4.58L19 12l-7 7z"/><path d="M5 20h14v2H5z"/></svg>
      G-Sign Test
    </a>
  </div>
  <p class="sign-hint">sso.gov.mn/login_signature руу чиглүүлэн тоон гарын үсгээр нэвтэрнэ</p>
</div>

<div class="features">
  <div class="feature">
    <div class="feature-icon purple">&#128274;</div>
    <h3>Клауд тоон гарын үсэг</h3>
    <p>GSign апп-аар 4 оронтой PIN кодоор баталгаажуулна. Сертификат клауд дээр хадгалагдана — төхөөрөмжөөс хамааралгүй.</p>
  </div>
  <div class="feature">
    <div class="feature-icon green">&#9989;</div>
    <h3>eIDAS / LoA 4</h3>
    <p>ETSI TS 102 204 стандарт, Methics Kiuru MSSP дээр суурилсан. Олон улсын eIDAS нийцтэй, хамгийн өндөр баталгаажуулалтын түвшин.</p>
  </div>
  <div class="feature">
    <div class="feature-icon amber">&#9889;</div>
    <h3>OAuth2 нэгдмэл</h3>
    <p>sso.gov.mn-ийн нэг authentication method. Одоо байгаа DAN OAuth2 flow-д нэмэлт интеграци шаардлагагүй.</p>
  </div>
</div>

<div class="info">
  <div class="info-card">
    <h2>G-Sign гэж юу вэ?</h2>
    <p>G-Sign (GSign) нь Монгол Улсын <strong>Улсын Бүртгэлийн Ерөнхий Газар (УБЕГ)</strong> болон <strong>Tridum e-Security</strong> компанийн хамтарсан клауд тоон гарын үсгийн үйлчилгээ юм. 2022 оны 5-р сард нээгдсэн бөгөөд нэг жилийн дотор 100,000+ хэрэглэгчтэй болсон.</p>

    <h3>Хэрхэн ажилладаг вэ?</h3>
    <div class="steps">
      <div class="step">
        <div class="step-num">1</div>
        <h4>GSign апп суулгах</h4>
        <p>Google Play / App Store-оос GSign Mongolia апп татаж суулгана.</p>
      </div>
      <div class="step">
        <div class="step-num">2</div>
        <h4>Бүртгүүлэх</h4>
        <p>УБЕГ-ийн киоск дээр QR код уншуулан 4 оронтой PIN үүсгэнэ.</p>
      </div>
      <div class="step">
        <div class="step-num">3</div>
        <h4>Сертификат авах</h4>
        <p>Клауд дээр 5 жилийн хүчинтэй тоон гарын үсгийн сертификат хадгалагдана.</p>
      </div>
      <div class="step">
        <div class="step-num">4</div>
        <h4>Гарын үсэг зурах</h4>
        <p>Баримт бичигт гарын үсэг зурахдаа апп дээр PIN оруулан баталгаажуулна.</p>
      </div>
    </div>

    <h3>Техникийн мэдээлэл</h3>
    <ul>
      <li><strong>Backend:</strong> Methics Kiuru MSSP (Mobile Signature Service Provider)</li>
      <li><strong>Протокол:</strong> ETSI TS 102 204 — SOAP & REST API</li>
      <li><strong>Клиент апп:</strong> Alauda PBY SDK суурилсан</li>
      <li><strong>Нийцэл:</strong> eIDAS compliant, Level of Assurance 4 (хамгийн өндөр)</li>
      <li><strong>PIN:</strong> 3 удаа буруу оруулбал түгжигдэнэ, PUK кодоор сэргээнэ</li>
      <li><strong>PUK:</strong> 10 удаа буруу оруулбал сертификат хүчингүй болно</li>
    </ul>

    <h3>Холбоосууд</h3>
    <ul>
      <li><a href="https://burtgel.gov.mn/g-sign">burtgel.gov.mn/g-sign</a> — УБЕГ G-Sign мэдээлэл</li>
      <li><a href="https://sso.gov.mn/login_signature">sso.gov.mn/login_signature</a> — G-Sign нэвтрэх</li>
      <li><a href="https://developer.sso.gov.mn">developer.sso.gov.mn</a> — Хөгжүүлэгчийн портал</li>
      <li><a href="https://play.google.com/store/apps/details?id=mn.tridumkey.sign_app_gov">Google Play — GSign Mongolia</a></li>
      <li><a href="https://apps.apple.com/mn/app/gsign-mongolia/id1632833615">App Store — GSign Mongolia</a></li>
    </ul>

    <h3>Gerege SSO-тай холбох</h3>
    <p>G-Sign нь sso.gov.mn-ийн нэг authentication method тул DAN OAuth2 flow-оор дамжуулан ашиглах боломжтой. <code>sso.gov.mn/login_signature</code> руу шууд чиглүүлэхэд хэрэглэгч GSign апп-аараа нэвтэрнэ.</p>
  </div>
</div>

<div class="footer">
  G-Sign Gateway &middot; <a href="https://gerege.mn">gerege.mn</a> &middot; Powered by УБЕГ GSign &amp; sso.gov.mn
</div>
</body>
</html>`
