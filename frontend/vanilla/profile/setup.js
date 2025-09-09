// Profile Setup JavaScript
class ProfileSetup {
    constructor() {
        this.form = document.getElementById('profileForm');
        this.fullNameInput = document.getElementById('fullName');
        this.submitBtn = document.getElementById('submitBtn');
        this.errorMessage = document.getElementById('errorMessage');
        this.errorText = document.getElementById('errorText');
        this.loadingState = document.getElementById('loadingState');
        
        this.init();
    }

    init() {
        this.form.addEventListener('submit', this.handleSubmit.bind(this));
        this.fullNameInput.addEventListener('input', this.clearError.bind(this));
        
        // Check if user is already logged in and has a name
        this.checkUserStatus();
    }

    async checkUserStatus() {
        const token = localStorage.getItem('accessToken');
        if (!token) {
            // Redirect to login if no token
            window.location.href = '/auth/login.html';
            return;
        }

        try {
            const response = await fetch('/api/auth/me', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                const user = await response.json();
                if (user.name && user.name.trim() !== '') {
                    // User already has a name, redirect to dashboard
                    window.location.href = '/home.html';
                }
            } else if (response.status === 401) {
                // Token expired, redirect to login
                localStorage.removeItem('accessToken');
                localStorage.removeItem('refreshToken');
                window.location.href = '/auth/login.html';
            }
        } catch (error) {
            console.error('Error checking user status:', error);
        }
    }

    async handleSubmit(e) {
        e.preventDefault();
        
        const fullName = this.fullNameInput.value.trim();
        
        if (!fullName) {
            this.showError('กรุณากรอกชื่อเต็มของคุณ');
            return;
        }

        if (fullName.length < 2) {
            this.showError('ชื่อต้องมีอย่างน้อย 2 ตัวอักษร');
            return;
        }

        await this.updateProfile(fullName);
    }

    async updateProfile(fullName) {
        this.setLoading(true);
        
        const token = localStorage.getItem('accessToken');
        
        try {
            const response = await fetch('/api/users/profile', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    name: fullName
                })
            });

            if (response.ok) {
                // Success! Redirect to dashboard
                window.location.href = '/home.html';
            } else {
                const errorData = await response.json();
                this.showError(errorData.message || 'เกิดข้อผิดพลาดในการบันทึกข้อมูล');
            }
        } catch (error) {
            console.error('Error updating profile:', error);
            this.showError('เกิดข้อผิดพลาดในการเชื่อมต่อ กรุณาลองใหม่อีกครั้ง');
        } finally {
            this.setLoading(false);
        }
    }

    showError(message) {
        this.errorText.textContent = message;
        this.errorMessage.classList.remove('hidden');
        
        // Auto hide error after 5 seconds
        setTimeout(() => {
            this.clearError();
        }, 5000);
    }

    clearError() {
        this.errorMessage.classList.add('hidden');
    }

    setLoading(isLoading) {
        if (isLoading) {
            this.form.classList.add('hidden');
            this.loadingState.classList.remove('hidden');
            this.submitBtn.disabled = true;
        } else {
            this.form.classList.remove('hidden');
            this.loadingState.classList.add('hidden');
            this.submitBtn.disabled = false;
        }
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new ProfileSetup();
});