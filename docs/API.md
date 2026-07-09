# Chirpy API Documentation

This document describes the API endpoints, request formats, authentication schemes, and responses for the Chirpy backend application.

## Base URL
When running locally, the API is available at:
`http://localhost:8080`

---

## Authentication Schemes

The API uses three types of authentication depending on the endpoint:

1. **Bearer Token (JWT Access Token)**
   - Used for authenticating user actions (e.g., creating/deleting chirps, updating user details).
   - Format: `Authorization: Bearer <jwt_access_token>`
   - Lifespan: 1 hour

2. **Bearer Token (Refresh Token)**
   - Used specifically to refresh access tokens or revoke active sessions.
   - Format: `Authorization: Bearer <refresh_token>`
   - Token format: 64-character hexadecimal string

3. **API Key**
   - Used by external webhooks/integrations (e.g., Polka upgrade webhook).
   - Format: `Authorization: ApiKey <api_key>`
   - Authenticated against the `POLKA_KEY` environment variable.

---

## Error Response Format

When a request fails, the API responds with a JSON error payload and a matching HTTP status code:

```json
{
  "error": "Error message explanation here"
}
```

---

## Endpoints

### 1. User & Session Endpoints

#### **Create User**
Creates a new user account.
* **Method**: `POST`
* **Path**: `/api/users`
* **Auth**: None
* **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```
* **Success Response**: `201 Created`
  ```json
  {
    "id": "e3e8ef6e-927b-4015-abfe-0a4465e94b29",
    "created_at": "2026-07-09T16:50:00Z",
    "updated_at": "2026-07-09T16:50:00Z",
    "email": "user@example.com",
    "is_chirpy_red": false
  }
  ```

#### **Login User**
Authenticates a user, returning a JWT access token and a refresh token.
* **Method**: `POST`
* **Path**: `/api/login`
* **Auth**: None
* **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```
* **Success Response**: `200 OK`
  ```json
  {
    "id": "e3e8ef6e-927b-4015-abfe-0a4465e94b29",
    "created_at": "2026-07-09T16:50:00Z",
    "updated_at": "2026-07-09T16:50:00Z",
    "email": "user@example.com",
    "is_chirpy_red": false,
    "token": "eyJhbGciOi...",
    "refresh_token": "a1b2c3d4..."
  }
  ```
* **Error Response**: `401 Unauthorized` (Incorrect email or password)

#### **Update User**
Updates the authenticated user's email and password.
* **Method**: `PUT`
* **Path**: `/api/users`
* **Auth**: Bearer Token (JWT Access Token)
  * **Header**: `Authorization: Bearer <jwt_access_token>`
* **Request Body**:
  ```json
  {
    "email": "new_email@example.com",
    "password": "newpassword"
  }
  ```
* **Success Response**: `200 OK`
  ```json
  {
    "id": "e3e8ef6e-927b-4015-abfe-0a4465e94b29",
    "created_at": "2026-07-09T16:50:00Z",
    "updated_at": "2026-07-09T16:55:00Z",
    "email": "new_email@example.com",
    "is_chirpy_red": false
  }
  ```
* **Error Response**: `401 Unauthorized` (Invalid/Missing token)

#### **Refresh Access Token**
Generates a new short-lived JWT access token using a valid refresh token.
* **Method**: `POST`
* **Path**: `/api/refresh`
* **Auth**: Bearer Token (Refresh Token)
  * **Header**: `Authorization: Bearer <refresh_token>`
* **Request Body**: None
* **Success Response**: `200 OK`
  ```json
  {
    "token": "new_jwt_access_token_here..."
  }
  ```
* **Error Response**: `401 Unauthorized` (Invalid/Expired token) or `400 Bad Request` (Missing token header)

#### **Revoke Refresh Token**
Logs out the user and revokes the active refresh token session.
* **Method**: `POST`
* **Path**: `/api/revoke`
* **Auth**: Bearer Token (Refresh Token)
  * **Header**: `Authorization: Bearer <refresh_token>`
* **Request Body**: None
* **Success Response**: `204 No Content`
* **Error Response**: `400 Bad Request` (Invalid/Missing token)

---

### 2. Chirps (Tweets) Endpoints

#### **Create Chirp**
Creates a new short message (Chirp).
* **Method**: `POST`
* **Path**: `/api/chirps`
* **Auth**: Bearer Token (JWT Access Token)
  * **Header**: `Authorization: Bearer <jwt_access_token>`
* **Request Body**:
  ```json
  {
    "body": "This is a chirp!"
  }
  ```
* **Success Response**: `201 Created`
  ```json
  {
    "id": "f5d7cf1a-42cd-41e9-913a-a1b2c3d4e5f6",
    "created_at": "2026-07-09T16:51:00Z",
    "updated_at": "2026-07-09T16:51:00Z",
    "body": "This is a chirp!",
    "user_id": "e3e8ef6e-927b-4015-abfe-0a4465e94b29"
  }
  ```
* **Error Response**: `401 Unauthorized` (Invalid/Expired token)

#### **Get All Chirps**
Retrieves a list of all chirps. Supports filtering by user and sorting by date.
* **Method**: `GET`
* **Path**: `/api/chirps`
* **Auth**: None
* **Query Parameters**:
  * `author_id` (optional): Filter chirps to a specific User UUID.
  * `sort` (optional): If set to `desc`, sorts the chirps from newest to oldest. By default, they are sorted ascending (oldest to newest).
* **Success Response**: `200 OK`
  ```json
  [
    {
      "id": "f5d7cf1a-42cd-41e9-913a-a1b2c3d4e5f6",
      "created_at": "2026-07-09T16:51:00Z",
      "updated_at": "2026-07-09T16:51:00Z",
      "body": "This is a chirp!",
      "user_id": "e3e8ef6e-927b-4015-abfe-0a4465e94b29"
    }
  ]
  ```

#### **Get Chirp by ID**
Retrieves details of a single chirp by its unique ID.
* **Method**: `GET`
* **Path**: `/api/chirps/{chirpID}`
* **Auth**: None
* **Path Parameters**:
  * `chirpID` (UUID)
* **Success Response**: `200 OK`
  ```json
  {
    "id": "f5d7cf1a-42cd-41e9-913a-a1b2c3d4e5f6",
    "created_at": "2026-07-09T16:51:00Z",
    "updated_at": "2026-07-09T16:51:00Z",
    "body": "This is a chirp!",
    "user_id": "e3e8ef6e-927b-4015-abfe-0a4465e94b29"
  }
  ```
* **Error Response**: `404 Not Found` (Couldn't get chirp)

#### **Delete Chirp**
Deletes a chirp. A user can only delete chirps they created.
* **Method**: `DELETE`
* **Path**: `/api/chirps/{chirpID}`
* **Auth**: Bearer Token (JWT Access Token)
  * **Header**: `Authorization: Bearer <jwt_access_token>`
* **Path Parameters**:
  * `chirpID` (UUID)
* **Success Response**: `204 No Content`
* **Error Response**:
  * `401 Unauthorized` (Invalid/Missing token)
  * `403 Forbidden` (Chirp does not belong to the authenticated user)
  * `404 Not Found` (Chirp does not exist)

---

### 3. Webhook Endpoints

#### **Polka Webhook**
Receives user subscription upgrade notifications from the Polka payment system. If the event is `user.upgraded`, the user's `is_chirpy_red` status is updated to `true`.
* **Method**: `POST`
* **Path**: `/api/polka/webhooks`
* **Auth**: API Key
  * **Header**: `Authorization: ApiKey <polka_api_key>`
* **Request Body**:
  ```json
  {
    "event": "user.upgraded",
    "data": {
      "user_id": "e3e8ef6e-927b-4015-abfe-0a4465e94b29"
    }
  }
  ```
* **Success Response**: `204 No Content` (Returned immediately for success, or for unhandled events like anything other than `user.upgraded`)
* **Error Response**:
  * `401 Unauthorized` (Invalid api key)
  * `404 Not Found` (Couldn't find user with the provided ID)

---

### 4. Admin & Health Endpoints

#### **Health Check**
Checks the status of the server.
* **Method**: `GET`
* **Path**: `/api/healthz`
* **Auth**: None
* **Success Response**: `200 OK`
  * **Content-Type**: `text/plain`
  * **Body**: `OK`

#### **Request Count Metrics**
Displays how many times the static assets file server (`/app/`) has been requested.
* **Method**: `GET`
* **Path**: `/admin/metrics`
* **Auth**: None
* **Success Response**: `200 OK`
  * **Content-Type**: `text/html`
  * **Body**:
    ```html
    <html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited 42 times!</p>
    </body>
    </html>
    ```

#### **Reset Server Metrics and DB**
Resets the file server request counter to 0 and clears all rows from the database. **Only enabled when the environment variable `PLATFORM` is set to `dev`.**
* **Method**: `POST`
* **Path**: `/admin/reset`
* **Auth**: None (Strictly checks `PLATFORM="dev"`)
* **Success Response**: `200 OK`
  * **Content-Type**: `text/plain`
  * **Body**: `Reset count`
* **Error Response**: `403 Forbidden` (If `PLATFORM` is not `dev`)
