const page = document.body.dataset.page;
const messageBox = document.getElementById("authMessage");

const showMessage = (text, type) => {
  messageBox.className = `alert alert-${type}`;
  messageBox.textContent = text;
  messageBox.classList.remove("d-none");
};

const loginForm = document.getElementById("loginForm");
if (loginForm) {
  loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const userId = document.getElementById("userId").value.trim();
    const userPwd = document.getElementById("userPwd").value.trim();
    const res = await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ userId, userPwd })
    });
    const data = await res.json().catch(() => ({}));
    if (res.ok) {
      localStorage.setItem("currentUser", userId);
      window.location.href = "/main";
      return;
    }
    showMessage(data.message || "登录失败", "danger");
  });
}

const registerForm = document.getElementById("registerForm");
if (registerForm) {
  registerForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const userId = document.getElementById("userId").value.trim();
    const userPwd = document.getElementById("userPwd").value.trim();
    const res = await fetch("/api/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ userId, userPwd })
    });
    const data = await res.json().catch(() => ({}));
    if (res.ok) {
      showMessage("注册成功，请登录", "success");
      setTimeout(() => {
        window.location.href = "/";
      }, 800);
      return;
    }
    showMessage(data.message || "注册失败", "danger");
  });
}
