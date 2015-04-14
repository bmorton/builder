package builds

import (
	"bytes"
	"database/sql"

	"code.google.com/p/go-uuid/uuid"
	"github.com/bmorton/builder/streams"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

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

func (r *Repository) Find(key string) *Build {
	var build Build
	r.db.First(&build, &Build{ID: key})
	build.BuildStream = r.buildStreams[key]
	build.PushStream = r.pushStreams[key]
	return &build
}

func (r *Repository) Create(build *Build) {
	temp := uuid.New()
	build.ID = temp
	r.db.Create(build)
	build.ID = temp
}

func (r *Repository) Save(build *Build) {
	r.buildStreams[build.ID] = build.BuildStream
	r.pushStreams[build.ID] = build.PushStream
	r.db.Save(build)
}

func (r *Repository) FindBuildLog(key string) *BuildLog {
	var log BuildLog
	r.db.First(&log, &BuildLog{ID: key})
	return &log
}

func (r *Repository) PersistStreams(key string) {
	build := r.Find(key)
	log := new(bytes.Buffer)
	build.BuildStream.Replay(log)
	buildLog := &BuildLog{
		ID:   build.ID,
		Data: log.String(),
	}

	log = new(bytes.Buffer)
	build.PushStream.Replay(log)
	pushLog := &PushLog{
		ID:   build.ID,
		Data: log.String(),
	}
	r.db.Create(&buildLog)
	r.db.Create(&pushLog)
}

func (r *Repository) DestroyStreams(key string) {
	delete(r.buildStreams, key)
	delete(r.pushStreams, key)
}

func (r *Repository) All() []*Build {
	all := make([]*Build, 0)
	r.db.Find(&all)
	return all
}

func (r *Repository) Migrate() {
	r.db.AutoMigrate(&Build{}, &BuildLog{}, &PushLog{})
}
