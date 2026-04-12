package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) renderLoginPage(w http.ResponseWriter, sessionID, clientName string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, loginPageHTML, clientName, sessionID)
}

const loginPageHTML = `<!DOCTYPE html>
<html lang="mn">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SSO нэвтрэлт — sso.gerege.mn</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:#f0f2f5;min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{background:#fff;border-radius:16px;box-shadow:0 4px 24px rgba(0,0,0,.08);max-width:420px;width:100%%;margin:20px;overflow:hidden}
.header{background:linear-gradient(135deg,#0047AB,#0033cc);padding:32px 24px;text-align:center;color:#fff}
.header svg{width:40px;height:40px;margin-bottom:12px}
.header h1{font-size:18px;font-weight:700;margin-bottom:4px}
.header p{font-size:13px;opacity:.8}
.body{padding:28px 24px}
.label{font-size:13px;font-weight:600;color:#374151;margin-bottom:6px}
.input{width:100%%;padding:14px 16px;border:2px solid #e5e7eb;border-radius:12px;font-size:16px;font-weight:500;outline:none;transition:border .2s;text-transform:uppercase;letter-spacing:1px}
.input:focus{border-color:#0047AB}
.btn{width:100%%;padding:16px;background:#0047AB;color:#fff;border:none;border-radius:12px;font-size:16px;font-weight:700;cursor:pointer;margin-top:16px;transition:background .2s}
.btn:hover{background:#003399}
.btn:disabled{background:#94a3b8;cursor:not-allowed}
.status{text-align:center;margin-top:16px;display:none}
.status.show{display:block}
.code-box{background:#f0f9ff;border:2px solid #0047AB;border-radius:12px;padding:16px;text-align:center;margin-top:16px;display:none}
.code-box.show{display:block}
.code-box .code{font-size:32px;font-weight:800;letter-spacing:4px;color:#0047AB}
.code-box .hint{font-size:12px;color:#6b7280;margin-top:8px}
.spinner{display:inline-block;width:20px;height:20px;border:3px solid #e5e7eb;border-top:3px solid #0047AB;border-radius:50%%;animation:spin 1s linear infinite;margin-right:8px;vertical-align:middle}
@keyframes spin{to{transform:rotate(360deg)}}
.error{color:#dc2626;font-size:13px;text-align:center;margin-top:12px;display:none}
.error.show{display:block}
.footer{text-align:center;padding:16px 24px 24px;font-size:12px;color:#9ca3af}
.footer a{color:#0047AB;text-decoration:none}
</style>
</head>
<body>
<div class="card">
  <div class="header">
    <svg viewBox="0 0 40 40" fill="none"><rect width="40" height="40" rx="10" fill="rgba(255,255,255,.2)"/><text x="50%%" y="55%%" dominant-baseline="middle" text-anchor="middle" fill="#fff" font-family="sans-serif" font-weight="800" font-size="16">G</text></svg>
    <h1>GeregeID</h1>
    <p>%s рүү нэвтрэх</p>
  </div>
  <div class="body">
    <div id="form-section">
      <div class="label">Регистрийн дугаар</div>
      <input id="national-id" class="input" type="text" placeholder="УБ12345678" maxlength="10" autocomplete="off">
      <button id="submit-btn" class="btn" onclick="initAuth()">Нэвтрэх</button>
    </div>
    <div id="code-section" class="code-box">
      <div class="code" id="verify-code">----</div>
      <div class="hint">GeregeID апп дээрх кодтой тулгана уу</div>
    </div>
    <div id="status-section" class="status">
      <span class="spinner"></span>
      <span id="status-text">Утсан дээрээ баталгаажуулна уу...</span>
    </div>
    <div id="error-section" class="error"></div>
  </div>
  <div class="footer">
    Powered by <a href="https://gerege.mn">GeregeID</a>
  </div>
</div>

<script>
const SESSION_ID = "%s";
let pollTimer = null;
let pollTimeout = null;

async function initAuth() {
  const nid = document.getElementById('national-id').value.trim().toUpperCase();
  if (!/^[\u0410-\u042F\u04E8\u04AE]{2}\d{8}$/.test(nid)) { showError('Регистрийн дугаар буруу формат (жнь: УБ12345678)'); return; }

  document.getElementById('submit-btn').disabled = true;
  document.getElementById('submit-btn').textContent = 'Холбогдож байна...';
  hideError();

  try {
    const resp = await fetch('/api/auth/initiate', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({session_id: SESSION_ID, national_id: nid})
    });
    const data = await resp.json();
    if (!resp.ok) { showError(data.error_description || 'Алдаа гарлаа'); resetBtn(); return; }

    // Show verification code
    document.getElementById('verify-code').textContent = data.verification_code || '----';
    document.getElementById('code-section').classList.add('show');
    document.getElementById('status-section').classList.add('show');
    document.getElementById('form-section').style.display = 'none';

    // Start polling
    startPolling();
  } catch(e) {
    showError('Сүлжээний алдаа'); resetBtn();
  }
}

function startPolling() {
  pollTimer = setInterval(async () => {
    try {
      const resp = await fetch('/api/auth/poll?session_id=' + SESSION_ID);
      const data = await resp.json();
      if (data.status === 'complete' && data.redirect_url) {
        stopPolling();
        document.getElementById('status-text').textContent = 'Амжилттай! Буцааж чиглүүлж байна...';
        window.location.href = data.redirect_url;
      } else if (data.status === 'EXPIRED' || data.status === 'CANCELLED') {
        stopPolling();
        showError('Хугацаа дууссан эсвэл цуцлагдсан');
        showForm();
      }
    } catch(e) {}
  }, 2000);
  pollTimeout = setTimeout(() => {
    stopPolling();
    showError('Хугацаа дууслаа (3 минут). Дахин оролдоно уу.');
    showForm();
  }, 180000);
}

function stopPolling() {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null; }
  if (pollTimeout) { clearTimeout(pollTimeout); pollTimeout = null; }
}

function showForm() {
  document.getElementById('form-section').style.display = '';
  document.getElementById('code-section').classList.remove('show');
  document.getElementById('status-section').classList.remove('show');
  resetBtn();
}

function showError(msg) { const el = document.getElementById('error-section'); el.textContent = msg; el.classList.add('show'); }
function hideError() { document.getElementById('error-section').classList.remove('show'); }
function resetBtn() { const btn = document.getElementById('submit-btn'); btn.disabled = false; btn.textContent = 'Нэвтрэх'; }

document.getElementById('national-id').addEventListener('keyup', function(e) { if (e.key === 'Enter') initAuth(); });
</script>
</div>
</body>
</html>`
