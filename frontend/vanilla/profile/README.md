# Profile Page

## Overview
The Profile page allows users to view their personal information and change their password.

## Features
- Display user information (name and email)
- Change password functionality

## API Endpoints
- **GET /user**: Fetch user information
- **POST /user/password**: Change user password

## Folder Structure
```
profile/
├── index.html
├── profile.js
├── styles.css
└── README.md
```

## Setup Instructions
1. Ensure the backend server is running.
2. Open `index.html` in your browser.

## Notes
- The page requires a valid access token stored in `localStorage` for API calls.
- Password change requires the current password and confirmation of the new password.
