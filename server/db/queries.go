package db

const (
	// Share queries
	Q_SHARE_FIND_BY_ID       = "SELECT file_id, downloads, expires_at FROM shares WHERE id = ?"
	Q_SHARE_INSERT           = "INSERT INTO shares (id, file_id, expires_at) VALUES (?, ?, ?)"
	Q_SHARE_UPDATE_DOWNLOADS = "UPDATE shares SET downloads = downloads + 1 WHERE id = ?"
	Q_SHARE_DELETE           = "DELETE FROM shares WHERE id = ?"

	// Folder queries
	Q_FOLDER_FIND_BY_ID    = "SELECT owner_id, file_ids, created_at, updated_at FROM folders WHERE id = ?"
	Q_FOLDER_FIND_BY_OWNER = "SELECT id, file_ids, created_at, updated_at FROM folders WHERE owner_id = ?"
	Q_FOLDER_INSERT        = "INSERT INTO folders (id, owner_id, file_ids) VALUES (?, ?, ?)"
	Q_FOLDER_UPDATE_FILES  = "UPDATE folders SET file_ids = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	Q_FOLDER_DELETE        = "DELETE FROM folders WHERE id = ?"

	// User queries
	Q_USER_CREATE        = "INSERT INTO users (id, email, password_hash, salt) VALUES (?, ?, ?, ?)"
	Q_USER_FIND_BY_EMAIL = "SELECT id, email, password_hash, salt, created_at FROM users WHERE email = ?"
	Q_USER_FIND_BY_ID    = "SELECT id, email, created_at FROM users WHERE id = ?"

	// File queries
	Q_FILE_INSERT        = "INSERT INTO files (id, owner_id) VALUES (?, ?)"
	Q_FILE_FIND_BY_ID    = "SELECT owner_id, uploaded_at, updated_at FROM files WHERE id = ?"
	Q_FILE_FIND_BY_OWNER = "SELECT id, uploaded_at, updated_at FROM files WHERE owner_id = ?"
	Q_FILE_DELETE        = "DELETE FROM files WHERE id = ?"

	// Upload session queries
	Q_UPLOAD_INSERT        = "INSERT INTO upload_sessions (id, user_id, file_id, expires_at) VALUES (?, ?, ?, ?)"
	Q_UPLOAD_CHUNKS_BY_ID  = "SELECT n_chunks, chunk_size, chunk_map FROM upload_sessions WHERE id = ?"
	Q_UPLOAD_UPDATE_CHUNKS = "UPDATE upload_sessions SET (n_chunks, chunk_size, chunk_map) = (?, ?, ?) WHERE id = ?"
	Q_UPLOAD_FIND_BY_ID    = "SELECT user_id, file_id, expires_at FROM upload_sessions WHERE id = ?"
	Q_UPLOAD_DELETE        = "DELETE FROM upload_sessions WHERE id = ?"

	// User session queries
	Q_SESSION_CREATE         = "INSERT INTO sessions (id, user_id, token, expires_at) VALUES (?, ?, ?, ?)"
	Q_SESSION_FIND_BY_TOKEN  = "SELECT user_id, expires_at FROM sessions WHERE token = ?"
	Q_SESSION_DELETE         = "DELETE FROM sessions WHERE token = ?"
	Q_SESSION_DELETE_EXPIRED = "DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP"
)
