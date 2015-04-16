package builds

import (
	"bytes"
	"database/sql"
	"time"

	"github.com/bmorton/builder/streams"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type BuildLog struct {
	ID        int       `json:"id"`
	BuildID   string    `json:"build_id", sql:"index:idx_build_id,index:idx_build_id_type"`
	Type      string    `json:"type", sql:"index:idx_build_id_type"`
	Data      string    `json:"data", sql:"blob"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository(driver string, db *sql.DB) *LogRepository {
	gormDB, _ := gorm.Open(driver, db)
	return &LogRepository{
		db: &gormDB,
	}
}

func (r *LogRepository) CreateFromOutput(stream *streams.BuildStream) (buildLog *BuildLog, pushLog *BuildLog) {
	log := new(bytes.Buffer)
	stream.BuildOutput.Replay(log)
	buildLog = &BuildLog{
		BuildID: stream.BuildID,
		Type:    "build",
		Data:    log.String(),
	}

	log = new(bytes.Buffer)
	stream.PushOutput.Replay(log)
	pushLog = &BuildLog{
		BuildID: stream.BuildID,
		Type:    "push",
		Data:    log.String(),
	}

	r.db.Create(buildLog)
	r.db.Create(pushLog)
	return
}

func (r *LogRepository) FindByBuildID(key string, logType string) (*BuildLog, error) {
	var buildLog BuildLog
	if r.db.First(&buildLog, &BuildLog{BuildID: key, Type: logType}).RecordNotFound() {
		return &buildLog, ErrNotFound
	}
	return &buildLog, nil
}

func (r *LogRepository) Migrate() {
	r.db.AutoMigrate(&BuildLog{})
}
