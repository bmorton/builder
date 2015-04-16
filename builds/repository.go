package builds

import (
	"database/sql"
	"errors"

	"code.google.com/p/go-uuid/uuid"
	"github.com/bmorton/builder/streams"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = errors.New("Build not found")

type Repository struct {
	buildStreams map[string]*streams.Output
	pushStreams  map[string]*streams.Output
	db           *gorm.DB
}

func NewRepository(driver string, db *sql.DB) *Repository {
	gormDB, _ := gorm.Open(driver, db)
	return &Repository{
		buildStreams: make(map[string]*streams.Output),
		pushStreams:  make(map[string]*streams.Output),
		db:           &gormDB,
	}
}

func (r *Repository) Find(key string) (*Build, error) {
	var build Build
	if r.db.First(&build, &Build{ID: key}).RecordNotFound() {
		return &build, ErrNotFound
	}
	return &build, nil
}

func (r *Repository) Create(build *Build) {
	temp := uuid.New()
	build.ID = temp
	r.db.Create(build)
	build.ID = temp
}

func (r *Repository) Save(build *Build) {
	r.db.Save(build)
}

func (r *Repository) FindBuildLog(key string) *BuildLog {
	var log BuildLog
	r.db.First(&log, &BuildLog{ID: key})
	return &log
}

func (r *Repository) DestroyStreams(key string) {
	delete(r.buildStreams, key)
	delete(r.pushStreams, key)
}

func (r *Repository) All() []*Build {
	all := make([]*Build, 0)
	r.db.Order("created_at DESC").Find(&all)
	return all
}

func (r *Repository) Migrate() {
	r.db.AutoMigrate(&Build{}, &BuildLog{}, &PushLog{})
}
