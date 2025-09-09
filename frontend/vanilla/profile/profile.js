document.addEventListener('DOMContentLoaded', () => {
    const userName = document.getElementById('user-name');
    const userEmail = document.getElementById('user-email');

    // Fetch user information
    async function fetchUserInfo() {
        try {
            const response = await fetch('http://localhost:8080/user', {
                method: 'GET',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('accessToken')}`,
                },
            });
            if (!response.ok) {
                throw new Error('Failed to fetch user information');
            }
            const user = await response.json();
            userName.textContent = user.name;
            userEmail.textContent = user.email;
        } catch (error) {
            console.error(error);
        }
    }

    // Handle password change
    const changePasswordForm = document.getElementById('change-password-form');
    changePasswordForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        const currentPassword = document.getElementById('current-password').value;
        const newPassword = document.getElementById('new-password').value;
        const confirmPassword = document.getElementById('confirm-password').value;

        if (newPassword !== confirmPassword) {
            alert('New passwords do not match');
            return;
        }

        try {
            const response = await fetch('http://localhost:8080/user/password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${localStorage.getItem('accessToken')}`,
                },
                body: JSON.stringify({ currentPassword, newPassword }),
            });
            if (!response.ok) {
                throw new Error('Failed to change password');
            }
            alert('Password changed successfully!');
            changePasswordForm.reset();
        } catch (error) {
            console.error(error);
        }
    });

    // Initialize
    fetchUserInfo();
});
