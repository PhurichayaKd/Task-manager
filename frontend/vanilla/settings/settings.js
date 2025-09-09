document.addEventListener('DOMContentLoaded', () => {
    const languageSelector = document.getElementById('language-selector');
    const themeToggle = document.getElementById('theme-toggle');
    const notificationForm = document.getElementById('notification-form');

    // Handle language change
    languageSelector.addEventListener('change', (event) => {
        const selectedLanguage = event.target.value;
        localStorage.setItem('language', selectedLanguage);
        alert(`Language changed to ${selectedLanguage}`);
    });

    // Handle theme toggle
    themeToggle.addEventListener('click', () => {
        document.body.classList.toggle('dark-theme');
        const currentTheme = document.body.classList.contains('dark-theme') ? 'dark' : 'light';
        localStorage.setItem('theme', currentTheme);
        alert(`Theme changed to ${currentTheme}`);
    });

    // Handle notification settings
    notificationForm.addEventListener('submit', (event) => {
        event.preventDefault();
        const emailNotifications = document.getElementById('email-notifications').checked;
        const smsNotifications = document.getElementById('sms-notifications').checked;
        localStorage.setItem('emailNotifications', emailNotifications);
        localStorage.setItem('smsNotifications', smsNotifications);
        alert('Notification settings saved!');
    });

    // Initialize settings
    const savedLanguage = localStorage.getItem('language') || 'en';
    languageSelector.value = savedLanguage;

    const savedTheme = localStorage.getItem('theme') || 'light';
    if (savedTheme === 'dark') {
        document.body.classList.add('dark-theme');
    }

    document.getElementById('email-notifications').checked = localStorage.getItem('emailNotifications') === 'true';
    document.getElementById('sms-notifications').checked = localStorage.getItem('smsNotifications') === 'true';
});
