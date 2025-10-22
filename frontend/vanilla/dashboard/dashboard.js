const API_BASE = 'https://task-manager-production-6c61.up.railway.app';

// Load auth utilities
const script = document.createElement('script');
script.src = '../utils/auth-utils.js';
script.onload = function() {
  // Check authentication and get token
  const token = localStorage.getItem('access_token') || window.AuthUtils.requireAuth('../index.html');
  window.authToken = token;
};
script.onerror = function() {
  console.error('Failed to load auth-utils.js');
  // Fallback authentication check
  const token = localStorage.getItem('access_token');
  if (!token) {
    location.replace('../index.html');
  } else {
    window.authToken = token;
  }
};
document.head.appendChild(script);

// Get token for API calls
function getToken() {
  return window.authToken || window.AuthUtils?.getAuthToken() || localStorage.getItem('tm_access_token');
}

const root = document.documentElement;
const themeBtn = document.getElementById('themeToggle');
const setTheme = m => { m==='dark'?root.classList.add('dark'):root.classList.remove('dark');
  localStorage.setItem('tm_theme', m); themeBtn.textContent = m==='dark'?'☀️':'🌙'; };
setTheme(localStorage.getItem('tm_theme')||'light');
themeBtn?.addEventListener('click',()=>setTheme(root.classList.contains('dark')?'light':'dark'));

document.getElementById('logoutBtn')?.addEventListener('click',()=>{
  if (window.AuthUtils) {
    window.AuthUtils.logout('../index.html');
  } else {
    localStorage.removeItem('access_token');
    location.replace('../index.html');
  }
});

const toast = (msg)=>{ const t = document.getElementById('toast'); t.textContent = msg; t.classList.remove('hidden');
  setTimeout(()=>t.classList.add('hidden'), 1600); };

async function api(path, opts={}) {
  const headers = { 'Content-Type': 'application/json', ...(opts.headers||{}) };
  const token = getToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  const res = await fetch(`${API_BASE}${path}`, { ...opts, headers });
  return res;
}

/* ------- state ------- */
let tasks = []; // {id,title,status,due_date,...}
const selected = new Set();

/* ------- UI helpers ------- */
const tbodyTodo = document.getElementById('tbodyTodo');
const tbodyDone = document.getElementById('tbodyDone');
const searchInput = document.getElementById('searchInput');
const statusFilter = document.getElementById('statusFilter');
const summaryText = document.getElementById('summaryText');

function statusPill(status){
  const map = { todo:'st-todo', doing:'st-doing', done:'st-done' };
  return `<span class="status-pill ${map[status]||'st-todo'}" data-status="${status}">${labelStatus(status)}</span>`;
}
function labelStatus(s){ return s==='done'?'Done':s==='doing'?'Doing':'Todo'; }

function rowTemplate(t){
  const due = t.due_date || '';
  return `<tr data-id="${t.id}">
    <td class="checkbox-col"><input type="checkbox" class="row-check"/></td>
    <td contenteditable="true" class="cell-title" spellcheck="false">${escapeHtml(t.title)}</td>
    <td class="col-status">
      <div class="status-cell">${statusPill(t.status)}</div>
    </td>
    <td class="col-date">
      <input type="date" class="input date-input" value="${due}">
    </td>
  </tr>`;
}
function cardTemplate(t){
  const due = t.due_date ? dayjs(t.due_date).format('MMM D') : '-';
  return `<div class="card-item" data-id="${t.id}">
    <div class="muted">Task</div>
    <div class="cell-title" contenteditable="true">${escapeHtml(t.title)}</div>
    <div class="card-foot">
      <div>${statusPill(t.status)}</div>
      <div>📅 ${due}</div>
    </div>
  </div>`;
}

function escapeHtml(s){ return s.replace(/[&<>"']/g, m => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[m])); }

function render(){
  const q = searchInput.value.trim().toLowerCase();
  const f = statusFilter.value;

  const view = document.getElementById('tableView').classList.contains('hidden') ? 'cards' : 'table';
  const filtered = tasks.filter(t=>{
    const byQ = !q || t.title.toLowerCase().includes(q);
    const byS = !f || t.status===f;
    return byQ && byS;
  });

  summaryText.textContent = `${filtered.length} tasks • ${tasks.filter(t=>t.status==='done').length} done`;

  if (view==='table'){
    const todo = filtered.filter(t=>t.status!=='done');
    const done = filtered.filter(t=>t.status==='done');
    tbodyTodo.innerHTML = todo.map(rowTemplate).join('');
    tbodyDone.innerHTML = done.map(rowTemplate).join('');
    bindTableEvents();
  } else {
    document.getElementById('cardsWrap').innerHTML = filtered.map(cardTemplate).join('');
    bindCardEvents();
  }

  document.getElementById('bulkDelete').disabled = selected.size===0;
}

function bindTableEvents(){
  document.querySelectorAll('tbody .row-check').forEach(cb=>{
    cb.addEventListener('change', (e)=>{
      const id = Number(e.target.closest('tr').dataset.id);
      e.target.checked ? selected.add(id) : selected.delete(id);
      document.getElementById('bulkDelete').disabled = selected.size===0;
    });
  });

  // inline title edit
  document.querySelectorAll('tbody .cell-title').forEach(cell=>{
    cell.addEventListener('blur', async (e)=>{
      const tr = e.target.closest('tr'); const id = Number(tr.dataset.id);
      const title = e.target.textContent.trim(); if(!title) return;
      await api(`/tasks/${id}`, { method:'PUT', body: JSON.stringify({ title }) });
      const t = tasks.find(x=>x.id===id); if(t){ t.title = title; }
      toast('Saved');
    });
  });

  // status change dropdown (simple)
  document.querySelectorAll('tbody .status-cell').forEach(box=>{
    box.addEventListener('click', async (e)=>{
      const tr = e.target.closest('tr'); const id = Number(tr.dataset.id);
      const cur = tasks.find(x=>x.id===id)?.status || 'todo';
      const next = cur==='todo'?'doing':cur==='doing'?'done':'todo'; // cycle
      await api(`/tasks/${id}`, { method:'PUT', body: JSON.stringify({ status: next }) });
      const t = tasks.find(x=>x.id===id); if(t){ t.status = next; }
      render();
    });
  });

  // due date
  document.querySelectorAll('tbody .date-input').forEach(inp=>{
    inp.addEventListener('change', async (e)=>{
      const tr = e.target.closest('tr'); const id = Number(tr.dataset.id);
      const due_date = e.target.value || "";
      await api(`/tasks/${id}`, { method:'PUT', body: JSON.stringify({ due_date }) });
      const t = tasks.find(x=>x.id===id); if(t){ t.due_date = due_date || null; }
      toast('Saved');
    });
  });
}

function bindCardEvents(){
  document.querySelectorAll('.card-item .cell-title').forEach(cell=>{
    cell.addEventListener('blur', async (e)=>{
      const id = Number(e.target.closest('.card-item').dataset.id);
      const title = e.target.textContent.trim(); if(!title) return;
      await api(`/tasks/${id}`, { method:'PUT', body: JSON.stringify({ title }) });
      const t = tasks.find(x=>x.id===id); if(t){ t.title = title; }
      toast('Saved');
    });
  });
  document.querySelectorAll('.card-item .status-pill').forEach(p=>{
    p.addEventListener('click', async (e)=>{
      const id = Number(e.target.closest('.card-item').dataset.id);
      const t = tasks.find(x=>x.id===id); if(!t) return;
      const next = t.status==='todo'?'doing':t.status==='doing'?'done':'todo';
      await api(`/tasks/${id}`, { method:'PUT', body: JSON.stringify({ status: next }) });
      t.status = next; render();
    });
  });
}

/* ------- load/create/delete ------- */
async function load(){
  const res = await api('/tasks?limit=200'); // ใช้ของจริง
  if (!res.ok){ console.error('Load failed'); return; }
  const data = await res.json();
  tasks = data.tasks || [];
  render();
}

async function createOne(status='todo'){
  const title = prompt('Task title');
  if (!title) return;
  const res = await api('/tasks', { method:'POST', body: JSON.stringify({ title, status }) });
  if (!res.ok){ toast('Create failed'); return; }
  const t = await res.json();
  if (t.due_date === undefined) t.due_date = null;
  tasks.unshift({ id:t.id, title:t.title, status:t.status, due_date:t.due_date });
  render(); toast('Created');
}

async function deleteSelected(){
  const ids = Array.from(selected);
  for (const id of ids){
    await api(`/tasks/${id}`, { method:'DELETE' });
    tasks = tasks.filter(t=>t.id!==id);
  }
  selected.clear(); render(); toast('Deleted');
}

/* ------- controls ------- */
const dropdown = document.querySelector('.dropdown');
document.getElementById('newBtn').addEventListener('click', ()=>{
  dropdown.classList.toggle('open');
});
document.getElementById('newMenu').addEventListener('click',(e)=>{
  if (e.target.matches('button[data-status]')){
    dropdown.classList.remove('open'); createOne(e.target.dataset.status);
  }
});
document.getElementById('addTodo').addEventListener('click',()=>createOne('todo'));
document.getElementById('addDone').addEventListener('click',()=>createOne('done'));
document.getElementById('bulkDelete').addEventListener('click', deleteSelected);

searchInput.addEventListener('input', render);
statusFilter.addEventListener('change', render);

const viewToggle = document.getElementById('viewToggle');
viewToggle.addEventListener('click', ()=>{
  const table = document.getElementById('tableView');
  const cards = document.getElementById('cardsView');
  const isTable = !table.classList.contains('hidden');
  if (isTable){ table.classList.add('hidden'); cards.classList.remove('hidden'); viewToggle.textContent = '📋 Table'; }
  else { cards.classList.add('hidden'); table.classList.remove('hidden'); viewToggle.textContent = '🗂 Cards'; }
  render();
});

/* select-all checkboxes */
document.getElementById('checkAllTodo').addEventListener('change', (e)=>{
  tbodyTodo.querySelectorAll('.row-check').forEach(cb=>{ cb.checked = e.target.checked;
    const id = Number(cb.closest('tr').dataset.id); e.target.checked?selected.add(id):selected.delete(id);});
  document.getElementById('bulkDelete').disabled = selected.size===0;
});
document.getElementById('checkAllDone').addEventListener('change', (e)=>{
  tbodyDone.querySelectorAll('.row-check').forEach(cb=>{ cb.checked = e.target.checked;
    const id = Number(cb.closest('tr').dataset.id); e.target.checked?selected.add(id):selected.delete(id);});
  document.getElementById('bulkDelete').disabled = selected.size===0;
});

window.addEventListener('click',(e)=>{ if(!dropdown.contains(e.target)) dropdown.classList.remove('open'); });

/* go! */
load();
