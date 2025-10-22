const API_BASE = 'https://task-manager-production-6c61.up.railway.app';
const DASHBOARD_URL = (location.pathname.includes('/vanilla/'))
  ? '/home.html'
  : '/home.html';

/* Theme */
const root = document.documentElement;
const themeBtn = document.getElementById('themeToggle');
const setTheme = (m)=>{ m==='dark'?root.classList.add('dark'):root.classList.remove('dark');
  localStorage.setItem('tm_theme', m); themeBtn.textContent = m==='dark'?'☀️':'🌙'; };
setTheme(localStorage.getItem('tm_theme')||'light');
themeBtn?.addEventListener('click',()=>setTheme(root.classList.contains('dark')?'light':'dark'));

/* i18n */
const dict={en:{hero_title:"Manage Your Tasks Efficiently",hero_sub:"Stay organized and boost your productivity with our intuitive task manager.",get_started:"Get Started",login_title:"Login",login_hint:"Use your email and password to continue.",login_btn:"Login"},
th:{hero_title:"จัดการงานของคุณอย่างมีประสิทธิภาพ",hero_sub:"จัดระเบียบและเพิ่มประสิทธิภาพด้วยตัวจัดการงานที่ใช้งานง่าย",get_started:"เริ่มต้น",login_title:"เข้าสู่ระบบ",login_hint:"กรอกอีเมลและรหัสผ่านเพื่อใช้งานต่อ",login_btn:"เข้าสู่ระบบ"}};
const langBtn=document.getElementById('langToggle');
const applyI18n=(lang)=>{document.querySelectorAll('[data-i18n]').forEach(el=>{el.textContent=(dict[lang]&&dict[lang][el.dataset.i18n])||el.textContent;});localStorage.setItem('tm_lang',lang);langBtn.textContent=lang==='th'?'🌐 TH':'🌐 EN';};
applyI18n(localStorage.getItem('tm_lang')||'en');
langBtn?.addEventListener('click',()=>applyI18n((localStorage.getItem('tm_lang')||'en')==='en'?'th':'en'));

/* show Dashboard link if already logged in */
const dashLink=document.getElementById('dashLink');
if(localStorage.getItem('tm_access_token')){dashLink.classList.remove('hidden');dashLink.href=DASHBOARD_URL;}

/* CTA scroll */
document.getElementById('getStartedBtn')?.addEventListener('click',()=>{
  document.getElementById('email')?.focus();
  window.scrollTo({top:document.getElementById('email').getBoundingClientRect().top+window.scrollY-120,behavior:'smooth'});
});

/* helpers */
async function postJSON(path, body){
  const res = await fetch(`${API_BASE}${path}`, {
    method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify(body)
  });
  return res;
}

/* Login */
const form=document.getElementById('loginForm');
const btn=document.getElementById('loginBtn');
const errBox=document.getElementById('errorBox');
form?.addEventListener('submit', async (e)=>{
  e.preventDefault(); errBox.classList.add('hidden'); btn.disabled=true; btn.style.opacity=.7;
  try{
    const email=document.getElementById('email').value.trim();
    const password=document.getElementById('password').value;
    const res=await postJSON('/auth/login',{email,password});
    if(!res.ok){ let msg=`HTTP ${res.status}`; try{const j=await res.json(); if(j?.error) msg=j.error;}catch{}; throw new Error(msg); }
    const data=await res.json();
    localStorage.setItem('tm_access_token', data.access_token);
    window.location.href = DASHBOARD_URL;
  }catch(err){ errBox.textContent=`Login failed: ${err.message}`; errBox.classList.remove('hidden'); }
  finally{ btn.disabled=false; btn.style.opacity=1; }
});

/* Tiny register (optional) */
const toggleRegister=document.getElementById('toggleRegister');
const regForm=document.getElementById('registerForm');
const regBtn=document.getElementById('registerBtn');
const regMsg=document.getElementById('registerMsg');

toggleRegister?.addEventListener('click',()=>{regForm.classList.toggle('hidden');});

regForm?.addEventListener('submit', async (e)=>{
  e.preventDefault(); regBtn.disabled=true;
  try{
    const email=document.getElementById('regEmail').value.trim();
    const password=document.getElementById('regPassword').value;
    const res=await postJSON('/auth/register',{email,password,role:'user'});
    if(!res.ok){ let msg=`HTTP ${res.status}`; try{const j=await res.json(); if(j?.error) msg=j.error;}catch{}; throw new Error(msg); }
    regMsg.classList.remove('hidden'); regForm.reset();
  }catch(err){ regMsg.textContent=`Create failed: ${err.message}`; regMsg.classList.remove('hidden'); }
  finally{ regBtn.disabled=false; }
});
