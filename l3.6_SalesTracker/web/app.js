const apiBase = "";

// --- Create Item ---
document.getElementById("item-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const form = e.target;
    const data = {
        type: form.type.value,
        category: form.category.value,
        amount: parseFloat(form.amount.value),
        date: form.date.value
    };

    const res = await fetch(`${apiBase}/items`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    });

    if (res.ok) {
        form.reset();
        loadItems();
        loadAnalytics();
    } else alert("Error creating item");
});


// --- Load Items ---
async function loadItems() {
    const form = document.getElementById("filter-form");
    const params = new URLSearchParams();
    if (form.from.value) params.append("from", form.from.value);
    if (form.to.value) params.append("to", form.to.value);
    if (form.category.value) params.append("category", form.category.value);
    if (form.limit.value) params.append("limit", form.limit.value);
    if (form.offset.value) params.append("offset", form.offset.value);

    const res = await fetch(`${apiBase}/items?${params.toString()}`);
    const items = await res.json();

    const tbody = document.getElementById("items-table-body");
    tbody.innerHTML = "";

    items.forEach(it => {
        const tr = document.createElement("tr");
        tr.innerHTML = `
            <td>${it.ID}</td>
            <td contenteditable="true" data-field="type">${it.Type}</td>
            <td contenteditable="true" data-field="category">${it.Category}</td>
            <td contenteditable="true" data-field="amount">${it.Amount}</td>
            <td contenteditable="true" data-field="date">${it.Date.slice(0, 10)}</td>
            <td>${it.CreatedAt.slice(0, 10)}</td>
            <td>
                <button class="update-btn" data-id="${it.ID}">Update</button>
                <button class="delete-btn" data-id="${it.ID}">Delete</button>
            </td>
        `;
        tbody.appendChild(tr);
    });

    // --- Update ---
    document.querySelectorAll(".update-btn").forEach(btn => {
        btn.onclick = async () => {
            const id = btn.dataset.id;
            const tr = btn.closest("tr");
            const updateData = {};

            tr.querySelectorAll("[contenteditable]").forEach(td => {
                const field = td.dataset.field;
                updateData[field] =
                    field === "amount" ? parseFloat(td.innerText) : td.innerText;
            });

            const res = await fetch(`${apiBase}/items/${id}`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(updateData)
            });

            if (res.ok) {
                loadItems();
                loadAnalytics();
            } else alert("Error updating item");
        };
    });

    // --- Delete ---
    document.querySelectorAll(".delete-btn").forEach(btn => {
        btn.onclick = async () => {
            const id = btn.dataset.id;
            if (confirm("Delete this item?")) {
                const res = await fetch(`${apiBase}/items/${id}`, { method: "DELETE" });
                if (res.ok) {
                    loadItems();
                    loadAnalytics();
                } else alert("Error deleting item");
            }
        };
    });
}


// --- Analytics ---
async function loadAnalytics() {
    const form = document.getElementById("filter-form");
    const params = new URLSearchParams();
    if (form.from.value) params.append("from", form.from.value);
    if (form.to.value) params.append("to", form.to.value);
    if (form.category.value) params.append("category", form.category.value);

    const res = await fetch(`${apiBase}/analytics?${params.toString()}`);
    if (!res.ok) return;

    const data = await res.json();

    document.getElementById("analytics-result").innerHTML = `
        <p><b>Sum:</b> ${data.sum}</p>
        <p><b>Avg:</b> ${data.avg}</p>
        <p><b>Median:</b> ${data.median}</p>
        <p><b>P90:</b> ${data.p90}</p>
        <p><b>Count:</b> ${data.count}</p>
    `;
}


// --- Auto-update analytics and items when filters change ---
const filterForm = document.getElementById("filter-form");

filterForm.querySelectorAll("input").forEach(input => {
    input.addEventListener("input", () => {
        loadItems();
        loadAnalytics();
    });
});


// --- Initial load ---
loadItems();
loadAnalytics();
