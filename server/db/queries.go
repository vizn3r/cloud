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
)
