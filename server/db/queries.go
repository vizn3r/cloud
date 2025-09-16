package db

const (
	Q_SHARE_FIND_BY_ID       = "SELECT file_id, downloads, expires_at FROM shares WHERE id = ?"
	Q_SHARE_INSERT           = "INSERT INTO shares (id, file_id, expires_at) VALUES (?, ?, ?)"
	Q_SHARE_UPDATE_DOWNLOADS = "UPDATE shares SET downloads = downloads + 1 WHERE id = ?"
	Q_SHARE_DELETE           = "DELETE FROM shares WHERE id = ?"

	Q_FOLDER_FIND_BY_ID      = "SELECT owner_id, file_ids, created_at, updated_at FROM folders WHERE id = ?"
	Q_FOLDER_FIND_BY_OWNER   = "SELECT id, file_ids, created_at, updated_at FROM folders WHERE owner_id = ?"
	Q_FOLDER_INSERT          = "INSERT INTO folders (id, owner_id, file_ids) VALUES (?, ?, ?)"
	Q_FOLDER_UPDATE_FILES    = "UPDATE folders SET file_ids = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	Q_FOLDER_DELETE          = "DELETE FROM folders WHERE id = ?"

	Q_USER_CREATE             = "INSERT INTO users (id, email, password_hash) VALUES (?, ?, ?)"
	Q_USER_FIND_BY_EMAIL      = "SELECT id, email, password_hash, created_at FROM users WHERE email = ?"
	Q_USER_FIND_BY_ID         = "SELECT id, email, created_at FROM users WHERE id = ?"

	Q_FILE_INSERT             = "INSERT INTO files (id, owner_id) VALUES (?, ?)"
	Q_FILE_FIND_BY_ID         = "SELECT owner_id, uploaded_at, updated_at FROM files WHERE id = ?"
	Q_FILE_FIND_BY_OWNER      = "SELECT id, uploaded_at, updated_at FROM files WHERE owner_id = ?"
	Q_FILE_DELETE             = "DELETE FROM files WHERE id = ?"

	Q_SESSION_CREATE          = "INSERT INTO sessions (id, user_id, token, expires_at) VALUES (?, ?, ?, ?)"
	Q_SESSION_FIND_BY_TOKEN   = "SELECT user_id, expires_at FROM sessions WHERE token = ?"
	Q_SESSION_DELETE          = "DELETE FROM sessions WHERE token = ?"
	Q_SESSION_DELETE_EXPIRED  = "DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP"
)
