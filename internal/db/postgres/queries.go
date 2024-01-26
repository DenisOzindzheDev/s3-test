package postgres

const (
	createRentSession   = "insert into sessions (user_id, rent_id) values ($1, $2)"
	startRentSession    = "update sessions set started_at = CURRENT_TIMESTAMP where rent_id = $1"
	completeRentSession = "update sessions set completed_at = CURRENT_TIMESTAMP where rent_id = $1"
	getRentSessions     = "select * from sessions where rent_id = $1"
	getStartedAtStatus  = "select * from sessions where rent_id = $1"
	updateImagesBefore  = "update sessions set images_before = $1 where rent_id = $2"
	updateImagesAfter   = "update sessions set images_after = $1 where rent_id = $2"
)
