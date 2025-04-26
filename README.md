
# Online Survey Platform

## Features
- **User**: Enter survey token, answer textbox/multiple-choice/Likert questions, submit responses.
- **Admin**: Authenticate with token, manage surveys, view responses.
- **Backend**: RESTful API with Go, MongoDB for surveys and responses.
- **Frontend**: Responsive HTML/JavaScript interface.

## Setup
### Prerequisites
- Go 1.24.2
- MongoDB Atlas account

### Installation
1. Clone the repository:

2. Create a `.env` file in the root directory:
    (ADMIN_TOKEN is for admin login)
    ```env
    MONGODB_URI=<your-MongoDB-URI>
    ADMIN_TOKEN=<your-admin-token>
    ```
3. Install Go dependencies:
    go mod tidy

4. Run the server:
    go run main.go

5. Access the app at `http://localhost:8080`.

## API Endpoints
- **POST /api/check-token**: Validate survey token.
- **GET /api/surveys/{token}**: Retrieve survey by token.
- **POST /api/responses**: Submit survey responses.
- **POST /api/admin/verify**: Authenticate admin.
- **GET /api/admin/surveys**: List all surveys.
- **POST /api/admin/surveys**: Create new survey.
- **PUT /api/admin/surveys/{token}**: Update survey.
- **DELETE /api/admin/surveys/{token}**: Delete survey.
- **GET /api/admin/responses/{token}**: View survey responses.

## Frontend Pages
- `/`: Enter survey token.
- `/survey.html?token=<token>`: Answer survey.
- `/admin-login.html`: Admin login.
- `/admin.html`: Manage surveys.
- `/create-survey.html`: Create new survey.
- `/edit-survey.html?token=<token>`: Edit survey.
- `/survey-detail.html?token=<token>`: View survey responses.