const API_BASE = "/comments";

async function loadComments() {
  try {
    const res = await fetch(API_BASE);
    if (!res.ok) throw new Error('Ошибка загрузки комментариев');
    const data = await res.json();

    const comments = data.result?.comments || data.comments || [];
    console.log('Получили комментарии:', comments);

    const container = document.getElementById("comments");
    container.innerHTML = "<h3>Все комментарии:</h3>";
    renderComments(comments, container);
  } catch (e) {
    alert(e.message);
  }
}

function renderComments(comments, container, level = 0) {
  comments.forEach(comment => {
    const div = document.createElement('div');
    div.className = 'comment';

    div.style.marginLeft = `${level * 20}px`;

    div.innerHTML = `
      <div><strong>ID ${comment.id}</strong>: ${comment.text}</div>
      <div class="actions">
        <button onclick="replyTo(${comment.id})">Ответить</button>
        <button onclick="deleteComment(${comment.id})">Удалить</button>
      </div>
    `;

    container.appendChild(div);

    if (comment.children && comment.children.length > 0) {
      renderComments(comment.children, container, level + 1);
    }
  });
}

async function addComment() {
  const text = document.getElementById("new-comment-text").value.trim();
  const parentIdRaw = document.getElementById("parent-id").value;
  const parentId = parentIdRaw.trim() === "" ? null : Number(parentIdRaw);

  if (!text) {
    alert("Введите текст комментария.");
    return;
  }

  try {
    const res = await fetch(API_BASE, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        text,
        parent_id: parentId
      }),
    });
    if (!res.ok) throw new Error('Ошибка при добавлении комментария');

    document.getElementById("new-comment-text").value = "";
    document.getElementById("parent-id").value = "";

    loadComments();
  } catch (e) {
    alert(e.message);
  }
}

async function deleteComment(id) {
  if (!confirm("Удалить комментарий и все вложенные?")) return;

  try {
    const res = await fetch(`${API_BASE}/${id}`, { method: "DELETE" });
    if (!res.ok) throw new Error('Ошибка при удалении комментария');
    loadComments();
  } catch (e) {
    alert(e.message);
  }
}

function replyTo(id) {
  document.getElementById("parent-id").value = id;
  document.getElementById("new-comment-text").focus();
}

async function searchComments() {
  const q = document.getElementById("search-input").value.trim();
  if (!q) return loadComments();

  try {
    const res = await fetch(`${API_BASE}?search=${encodeURIComponent(q)}`);
    if (!res.ok) throw new Error('Ошибка при поиске комментариев');
    const data = await res.json();

    const comments = data.result?.comments || data.comments || [];

    const container = document.getElementById("comments");
    container.innerHTML = "<h3>Результаты поиска:</h3>";
    renderComments(comments, container);
  } catch (e) {
    alert(e.message);
  }
}

window.addEventListener("DOMContentLoaded", loadComments);
