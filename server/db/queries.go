package db

const (
	Q_SHARE_FIND_BY_ID       = "SELECT file_id, downloads FROM shares WHERE id = ?"
	Q_SHARE_INSERT           = "INSERT INTO shares (id, file_id) VALUES (?, ?)"
	Q_SHARE_UPDATE_DOWNLOADS = "UPDATE shares SET downloads = downloads + 1 WHERE id = ?"
	Q_SHARE_DELETE           = "DELETE FROM shares WHERE id = ?"
)
