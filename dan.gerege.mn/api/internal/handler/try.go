package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strings"

	danpkg "dan.gerege.mn/api/internal/dan"
)

// Try initiates a standalone DAN verification — no client registration needed.
// Redirects to sso.gov.mn, and on callback renders citizen data directly.
// Try initiates a standalone DAN verification — no client registration needed.
// Uses the same /authorized callback as the normal flow, but state mode="try"
// tells authorized handler to render data directly instead of POSTing to a callback.
func (h *Handler) Try(w http.ResponseWriter, r *http.Request) {
	statePayload, _ := json.Marshal(map[string]string{
		"mode": "try",
	})
	signedState := signState(statePayload, h.cfg.StateSecret)

	loginURL := fmt.Sprintf("https://sso.gov.mn/login?state=%s&grant_type=authorization_code&response_type=code&client_id=%s&scope=%s&redirect_uri=%s",
		url.QueryEscape(signedState),
		url.QueryEscape(h.cfg.DAN.ClientID),
		url.QueryEscape(h.cfg.DAN.Scope),
		url.QueryEscape(h.cfg.DAN.CallbackURI),
	)

	slog.Info("try: redirecting to sso.gov.mn")
	http.Redirect(w, r, loginURL, http.StatusFound)
}

// HandleTryResult fetches citizen data and renders it as HTML.
// Called from Authorized when state mode is "try".
func (h *Handler) HandleTryResult(w http.ResponseWriter, r *http.Request, code string) {
	accessToken, err := danpkg.GetAccessToken(h.cfg.DAN, code)
	if err != nil {
		slog.Error("try: token exchange failed", "error", err)
		renderTryError(w, "access_token авахад алдаа гарлаа")
		return
	}

	citizen, err := danpkg.GetCitizenData(h.cfg.DAN, accessToken)
	if err != nil {
		slog.Error("try: citizen data failed", "error", err)
		renderTryError(w, "Иргэний мэдээлэл авахад алдаа гарлаа")
		return
	}

	slog.Info("try: success")
	renderTryResult(w, citizen)
}

func renderTryError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(400)
	fmt.Fprintf(w, tryErrorHTML, html.EscapeString(msg))
}

func renderTryResult(w http.ResponseWriter, citizen map[string]string) {
	// Separate image from other fields
	image := citizen["image"]

	// Build rows sorted by key
	type row struct{ Key, Label, Value string }
	labelMap := map[string]string{
		"reg_no": "Регистрийн дугаар", "given_name": "Нэр", "family_name": "Овог",
		"surname": "Ургийн овог", "gender": "Хүйс", "birth_date": "Төрсөн огноо",
		"birth_place": "Төрсөн газар", "nationality": "Иргэншил", "civil_id": "Иргэний ID",
		"aimag_name": "Аймаг/Хот", "aimag_code": "Аймаг код", "sum_name": "Сум/Дүүрэг",
		"sum_code": "Сум код", "bag_name": "Баг/Хороо", "bag_code": "Баг код",
		"address_detail": "Хаягийн дэлгэрэнгүй", "passport_address": "Паспортын хаяг",
		"passport_expire_date": "Паспорт дуусах", "passport_issue_date": "Паспорт олгосон",
		"apartment_name": "Байр", "street_name": "Гудамж",
	}

	// Priority order for display
	priority := []string{
		"reg_no", "family_name", "given_name", "surname", "gender", "birth_date",
		"birth_place", "nationality", "civil_id",
		"aimag_name", "sum_name", "bag_name",
		"address_detail", "street_name", "apartment_name",
		"passport_address", "passport_issue_date", "passport_expire_date",
	}

	var rows []row
	seen := map[string]bool{}
	for _, k := range priority {
		if v, ok := citizen[k]; ok && v != "" && k != "image" {
			label := labelMap[k]
			if label == "" {
				label = k
			}
			rows = append(rows, row{k, label, v})
			seen[k] = true
		}
	}
	// Append remaining keys
	remaining := make([]string, 0)
	for k := range citizen {
		if !seen[k] && k != "image" {
			remaining = append(remaining, k)
		}
	}
	sort.Strings(remaining)
	for _, k := range remaining {
		if v := citizen[k]; v != "" {
			label := labelMap[k]
			if label == "" {
				label = k
			}
			rows = append(rows, row{k, label, v})
		}
	}

	// Build table HTML
	var tableRows strings.Builder
	for _, r := range rows {
		tableRows.WriteString(fmt.Sprintf(
			`<tr><td class="lbl">%s</td><td class="val">%s</td></tr>`,
			html.EscapeString(r.Label), html.EscapeString(r.Value),
		))
	}

	// Image tag
	imageHTML := ""
	if image != "" {
		imageHTML = fmt.Sprintf(`<img src="data:image/jpeg;base64,%s" alt="photo" class="photo"/>`, image)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, tryResultHTML, imageHTML, tableRows.String())
}

const tryErrorHTML = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DAN Verify — Алдаа</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#0b1120;color:#e2e8f0;min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{background:rgba(255,255,255,.03);border:1px solid rgba(239,68,68,.2);border-radius:20px;padding:40px;max-width:460px;text-align:center}
.icon{font-size:48px;margin-bottom:16px}
h1{font-size:20px;font-weight:700;color:#f87171;margin-bottom:8px}
p{font-size:14px;color:#94a3b8;line-height:1.6;margin-bottom:24px}
a{display:inline-block;padding:12px 28px;background:rgba(255,255,255,.06);border:1px solid rgba(255,255,255,.1);color:#fff;text-decoration:none;border-radius:12px;font-weight:600;font-size:14px}
a:hover{background:rgba(255,255,255,.1)}
</style>
</head>
<body>
<div class="card">
  <div class="icon">&#9888;</div>
  <h1>Алдаа</h1>
  <p>%s</p>
  <a href="/try">Дахин оролдох</a>
</div>
</body>
</html>`

const tryResultHTML = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DAN Verify — Үр дүн</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#0b1120;color:#e2e8f0;min-height:100vh;padding:32px 16px}
.wrap{max-width:640px;margin:0 auto}
.header{display:flex;align-items:center;gap:12px;margin-bottom:32px}
.logo{width:36px;height:36px;background:#2563eb;border-radius:10px;display:flex;align-items:center;justify-content:center;color:#fff;font-weight:800;font-size:10px}
.header h1{font-size:20px;font-weight:700;color:#fff}
.header .badge{margin-left:auto;padding:5px 14px;background:rgba(22,163,74,.1);border:1px solid rgba(22,163,74,.25);border-radius:20px;font-size:11px;color:#4ade80;font-weight:600}
.profile{display:flex;gap:24px;margin-bottom:24px;align-items:flex-start}
.photo{width:120px;height:150px;border-radius:14px;object-fit:cover;border:2px solid rgba(255,255,255,.1);flex-shrink:0}
.info{flex:1}
.info .name{font-size:22px;font-weight:800;color:#fff;margin-bottom:4px}
.info .reg{font-size:14px;color:#60a5fa;font-family:monospace;margin-bottom:12px}
.info .meta{font-size:12px;color:#64748b;line-height:1.8}
.info .meta span{color:#94a3b8}
.table-wrap{background:rgba(255,255,255,.03);border:1px solid rgba(255,255,255,.06);border-radius:16px;overflow:hidden;margin-bottom:24px}
table{width:100%%;border-collapse:collapse}
tr{border-bottom:1px solid rgba(255,255,255,.04)}
tr:last-child{border-bottom:none}
td{padding:10px 16px;font-size:13px}
td.lbl{color:#64748b;font-weight:600;width:40%%}
td.val{color:#e2e8f0}
.actions{display:flex;gap:12px}
.btn{padding:12px 28px;font-weight:700;font-size:14px;border-radius:12px;text-decoration:none;transition:all .2s;display:inline-flex;align-items:center;gap:8px}
.btn-primary{background:linear-gradient(135deg,#2563eb,#1d4ed8);color:#fff;box-shadow:0 4px 16px rgba(37,99,235,.3)}
.btn-primary:hover{transform:translateY(-2px);box-shadow:0 8px 24px rgba(37,99,235,.4)}
.btn-outline{border:1px solid rgba(255,255,255,.15);color:#fff;background:transparent}
.btn-outline:hover{background:rgba(255,255,255,.05)}
@media(max-width:520px){.profile{flex-direction:column;align-items:center;text-align:center}.photo{width:100px;height:125px}}
</style>
</head>
<body>
<div class="wrap">
  <div class="header">
    <div class="logo">DAN</div>
    <h1>DAN Verify</h1>
    <div class="badge">&#10003; Амжилттай</div>
  </div>

  <div class="profile">
    %s
    <div class="info"></div>
  </div>

  <div class="table-wrap">
    <table>%s</table>
  </div>

  <div class="actions">
    <a href="/try" class="btn btn-primary">Дахин шалгах</a>
    <a href="/" class="btn btn-outline">Нүүр хуудас</a>
  </div>
</div>
</body>
</html>`
