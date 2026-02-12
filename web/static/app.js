const goodsBody = document.getElementById("goodsBody");
const totalCostEl = document.getElementById("totalCost");
const toastEl = document.getElementById("toast");
const toastBody = document.getElementById("toastBody");
const toast = toastEl ? new bootstrap.Toast(toastEl) : null;

const showToast = (text) => {
  if (!toast) {
    return;
  }
  toastBody.textContent = text;
  toast.show();
};

const loadGoods = async () => {
  const res = await fetch("/api/goods");
  const data = await res.json().catch(() => ({ items: [], totalCost: 0 }));
  if (!res.ok) {
    showToast("加载失败");
    return;
  }
  renderGoods(data.items || []);
  totalCostEl.textContent = Number(data.totalCost || 0).toFixed(2);
};

const renderGoods = (items) => {
  goodsBody.innerHTML = "";
  items.forEach((item) => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${item.gid}</td>
      <td>${item.gName}</td>
      <td>${item.gShelf}</td>
      <td>${item.gCnt}</td>
      <td>${Number(item.gPrice).toFixed(2)}</td>
      <td>${Number(item.cost).toFixed(2)}</td>
      <td>
        <button class="btn btn-sm btn-outline-danger" data-id="${item.gid}">删除</button>
      </td>
    `;
    tr.querySelector("button").addEventListener("click", () => deleteGoods(item.gid));
    goodsBody.appendChild(tr);
  });
};

const deleteGoods = async (gid) => {
  const res = await fetch(`/api/goods/${gid}`, { method: "DELETE" });
  if (!res.ok) {
    showToast("删除失败");
    return;
  }
  showToast("已删除");
  loadGoods();
};

const addForm = document.getElementById("addForm");
if (addForm) {
  addForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const payload = {
      items: [
        {
          gName: document.getElementById("addName").value.trim(),
          gShelf: document.getElementById("addShelf").value.trim(),
          gCnt: Number(document.getElementById("addCnt").value),
          gPrice: Number(document.getElementById("addPrice").value)
        }
      ]
    };
    const res = await fetch("/api/goods", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (!res.ok) {
      showToast("新增失败");
      return;
    }
    showToast("新增成功");
    addForm.reset();
    bootstrap.Modal.getInstance(document.getElementById("addModal")).hide();
    loadGoods();
  });
}

const stockForm = document.getElementById("stockForm");
if (stockForm) {
  stockForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const payload = {
      gid: Number(document.getElementById("stockId").value),
      det: Number(document.getElementById("stockCnt").value)
    };
    const res = await fetch("/api/goods/stock", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (!res.ok) {
      showToast("出库失败");
      return;
    }
    showToast("出库成功");
    stockForm.reset();
    bootstrap.Modal.getInstance(document.getElementById("stockModal")).hide();
    loadGoods();
  });
}

const logoutBtn = document.getElementById("logoutBtn");
if (logoutBtn) {
  logoutBtn.addEventListener("click", () => {
    localStorage.removeItem("currentUser");
    window.location.href = "/";
  });
}

const currentUser = localStorage.getItem("currentUser");
const userEl = document.getElementById("currentUser");
if (userEl) {
  userEl.textContent = currentUser ? `你好，${currentUser}` : "未登录";
}

loadGoods();
