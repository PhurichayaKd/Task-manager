// Settings functionality
document.addEventListener('DOMContentLoaded', function() {
    loadUserSettings();
    initializeEventListeners();
});

function loadUserSettings() {
    // Load saved settings from localStorage or API
    const savedSettings = JSON.parse(localStorage.getItem('userSettings')) || {};
    
    // Apply saved language
    if (savedSettings.language) {
        document.getElementById('language-select').value = savedSettings.language;
    }
    
    // Apply saved theme
    if (savedSettings.theme === 'dark') {
        document.getElementById('theme-toggle').checked = true;
        document.body.classList.add('dark-theme');
    }
    
    // Apply notification settings
    if (savedSettings.notifications) {
        document.getElementById('email-notifications').checked = savedSettings.notifications.email || false;
        document.getElementById('push-notifications').checked = savedSettings.notifications.push || false;
        document.getElementById('task-reminders').checked = savedSettings.notifications.reminders || false;
    }
}

function initializeEventListeners() {
    // Language settings
    const languageSelect = document.getElementById('language-select');
    languageSelect.addEventListener('change', function() {
        console.log('Language changed to:', this.value);
        showNotification('Language setting will be applied after page reload.');
    });

    // Theme settings
    const themeToggle = document.getElementById('theme-toggle');
    themeToggle.addEventListener('change', function() {
        document.body.classList.toggle('dark-theme', this.checked);
        console.log('Theme toggled to:', this.checked ? 'dark' : 'light');
        showNotification('Theme changed successfully!');
    });

    // Password form
    const passwordForm = document.getElementById('password-form');
    passwordForm.addEventListener('submit', function(e) {
        e.preventDefault();
        handlePasswordChange();
    });
}

function changeProfilePicture() {
    const fileInput = document.getElementById('avatar-upload');
    fileInput.click();
    
    fileInput.addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function(e) {
                // Update avatar display
                const avatars = document.querySelectorAll('.current-avatar, .dropdown-avatar');
                const newImageSrc = e.target.result;
                avatars.forEach(avatar => {
                    avatar.style.backgroundImage = `url(${newImageSrc})`;
                    avatar.style.backgroundSize = 'cover';
                    avatar.style.backgroundPosition = 'center';
                    avatar.textContent = '';
                });
                
                // อัปเดตรูปโปรไฟล์ในหน้า home.html ทันที
                localStorage.setItem('userProfileImage', newImageSrc);
                
                // อัปเดตไอคอนในหน้า home หากเปิดอยู่
                updateHomeProfileImage(newImageSrc);
                
                showNotification('Profile picture updated successfully!');
            };
            reader.readAsDataURL(file);
        }
    });
}

function editUsername() {
    const usernameInput = document.getElementById('username');
    const editButton = usernameInput.nextElementSibling;
    
    if (usernameInput.readOnly) {
        usernameInput.readOnly = false;
        usernameInput.focus();
        editButton.textContent = 'Save';
        usernameInput.style.borderColor = '#6366f1';
    } else {
        usernameInput.readOnly = true;
        editButton.textContent = 'Edit';
        usernameInput.style.borderColor = '#d1d5db';
        
        // บันทึกชื่อผู้ใช้ใหม่ใน localStorage
        localStorage.setItem('username', usernameInput.value);
        
        // อัปเดตชื่อผู้ใช้ในหน้า home ทันที
        updateHomeUsername(usernameInput.value);
        
        showNotification('ชื่อผู้ใช้ถูกอัปเดตแล้ว!', 'success');
    }
}

function openPasswordModal() {
    document.getElementById('passwordModal').style.display = 'block';
}

function closePasswordModal() {
    document.getElementById('passwordModal').style.display = 'none';
    document.getElementById('password-form').reset();
}

function handlePasswordChange() {
    const currentPassword = document.getElementById('current-password').value;
    const newPassword = document.getElementById('new-password').value;
    const confirmPassword = document.getElementById('confirm-password').value;
    
    if (newPassword !== confirmPassword) {
        showNotification('New passwords do not match!', 'error');
        return;
    }
    
    if (newPassword.length < 8) {
        showNotification('Password must be at least 8 characters long!', 'error');
        return;
    }
    
    // Simulate API call
    setTimeout(() => {
        closePasswordModal();
        showNotification('Password changed successfully!');
    }, 1000);
}

function confirmDeleteAccount() {
    document.getElementById('deleteAccountModal').style.display = 'block';
    
    // เปิดใช้งานปุ่มลบเมื่อมีการกรอกข้อมูลครบถ้วน
    const deletePassword = document.getElementById('deletePassword');
    const confirmDelete = document.getElementById('confirmDelete');
    const deleteBtn = document.getElementById('deleteConfirmBtn');
    
    function checkDeleteForm() {
        if (deletePassword.value.trim() && confirmDelete.checked) {
            deleteBtn.disabled = false;
            deleteBtn.style.opacity = '1';
        } else {
            deleteBtn.disabled = true;
            deleteBtn.style.opacity = '0.5';
        }
    }
    
    deletePassword.addEventListener('input', checkDeleteForm);
    confirmDelete.addEventListener('change', checkDeleteForm);
}

function closeDeleteAccountModal() {
    const modal = document.getElementById('deleteAccountModal');
    modal.style.display = 'none';
    
    // รีเซ็ตฟอร์ม
    document.getElementById('deletePassword').value = '';
    document.getElementById('confirmDelete').checked = false;
    document.getElementById('deleteConfirmBtn').disabled = true;
    document.getElementById('deleteConfirmBtn').style.opacity = '0.5';
}

function finalDeleteAccount() {
    const password = document.getElementById('deletePassword').value;
    
    // ตรวจสอบรหัสผ่าน (ในระบบจริงจะต้องส่งไปตรวจสอบกับเซิร์ฟเวอร์)
    const storedPassword = localStorage.getItem('userPassword') || 'password123'; // ค่าเริ่มต้นสำหรับการทดสอบ
    
    if (password !== storedPassword) {
        showNotification('รหัสผ่านไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง', 'error');
        return;
    }
    
    // ลบข้อมูลผู้ใช้ทั้งหมด
    localStorage.removeItem('userProfileImage');
    localStorage.removeItem('username');
    localStorage.removeItem('userPassword');
    localStorage.removeItem('userEmail');
    
    closeDeleteAccountModal();
    showNotification('บัญชีถูกลบเรียบร้อยแล้ว กำลังออกจากระบบ...', 'success');
    
    // ออกจากระบบและเปลี่ยนเส้นทางไปหน้าแรก
    setTimeout(() => {
        window.location.href = '../index.html';
    }, 2000);
}

function saveSettings() {
    const settings = {
        language: document.getElementById('language-select').value,
        theme: document.getElementById('theme-toggle').checked ? 'dark' : 'light',
        notifications: {
            email: document.getElementById('email-notifications').checked,
            push: document.getElementById('push-notifications').checked,
            reminders: document.getElementById('task-reminders').checked
        },
        username: document.getElementById('username').value
    };
    
    // Save to localStorage (in real app, this would be an API call)
    localStorage.setItem('userSettings', JSON.stringify(settings));
    
    showNotification('Settings saved successfully!');
}

function resetSettings() {
    const confirmed = confirm('Are you sure you want to reset all settings to default?');
    
    if (confirmed) {
        // Reset form values
        document.getElementById('language-select').value = 'en';
        document.getElementById('theme-toggle').checked = false;
        document.getElementById('email-notifications').checked = true;
        document.getElementById('push-notifications').checked = true;
        document.getElementById('task-reminders').checked = false;
        document.getElementById('username').value = 'Admin User';
        
        // Remove dark theme
        document.body.classList.remove('dark-theme');
        
        // Clear localStorage
        localStorage.removeItem('userSettings');
        
        showNotification('Settings reset to default values!');
    }
}

function updateHomeProfileImage(imageSrc) {
    // ฟังก์ชันสำหรับอัปเดตรูปโปรไฟล์ในหน้า home
    try {
        if (window.parent && window.parent !== window) {
            // หากอยู่ใน iframe
            window.parent.postMessage({
                type: 'updateProfileImage',
                imageSrc: imageSrc
            }, '*');
        } else {
            // หากเปิดในหน้าต่างเดียวกัน
            const homeUserAvatar = document.querySelector('.user-avatar');
            if (homeUserAvatar) {
                homeUserAvatar.style.backgroundImage = `url(${imageSrc})`;
            }
        }
    } catch (error) {
        console.log('ไม่สามารถอัปเดตรูปโปรไฟล์ในหน้า home ได้');
    }
}

function updateHomeUsername(username) {
    // ฟังก์ชันสำหรับอัปเดตชื่อผู้ใช้ในหน้า home
    try {
        if (window.parent && window.parent !== window) {
            // หากอยู่ใน iframe
            window.parent.postMessage({
                type: 'updateUsername',
                username: username
            }, '*');
        } else {
            // หากเปิดในหน้าต่างเดียวกัน
            const dropdownUsername = document.querySelector('.dropdown-username');
            if (dropdownUsername) {
                dropdownUsername.textContent = username;
            }
        }
    } catch (error) {
        console.log('ไม่สามารถอัปเดตชื่อผู้ใช้ในหน้า home ได้');
    }
}

function showNotification(message, type = 'success') {
    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    
    // Style the notification
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 12px 20px;
        border-radius: 8px;
        color: white;
        font-weight: 500;
        z-index: 10000;
        transform: translateX(100%);
        transition: transform 0.3s ease;
        ${type === 'success' ? 'background: #10b981;' : ''}
        ${type === 'error' ? 'background: #ef4444;' : ''}
        ${type === 'info' ? 'background: #3b82f6;' : ''}
    `;
    
    document.body.appendChild(notification);
    
    // Animate in
    setTimeout(() => {
        notification.style.transform = 'translateX(0)';
    }, 100);
    
    // Remove after 3 seconds
    setTimeout(() => {
        notification.style.transform = 'translateX(100%)';
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
}

// Close modal when clicking outside
window.addEventListener('click', function(event) {
    const modal = document.getElementById('passwordModal');
    if (event.target === modal) {
        closePasswordModal();
    }
});

// Handle escape key for modal
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        closePasswordModal();
    }
});
