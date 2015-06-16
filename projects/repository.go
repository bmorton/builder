package projects

import (
	"database/sql"
	"errors"

	"code.google.com/p/go-uuid/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = errors.New("Project not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(driver string, db *sql.DB) *Repository {
	gormDB, _ := gorm.Open(driver, db)
	return &Repository{
		db: &gormDB,
	}
}

func (r *Repository) Find(key string) (*Project, error) {
	var record Project
	if r.db.First(&record, &Project{ID: key}).RecordNotFound() {
		return &record, ErrNotFound
	}
	return &record, nil
}

func (r *Repository) Create(record *Project) {
	temp := uuid.New()
	record.ID = temp
	r.db.Create(record)
	record.ID = temp
}

func (r *Repository) Save(record *Project) {
	r.db.Save(record)
}

func (r *Repository) All() []*Project {
	all := make([]*Project, 0)
	r.db.Order("name ASC").Find(&all)
	return all
}

func (r *Repository) Migrate() {
	r.db.AutoMigrate(&Project{})
}

func (r *Repository) FindOrCreateByCloneURL(cloneURL string) *Project {
	var record Project
	r.db.Where(Project{CloneURL: cloneURL}).FirstOrInit(&record)

	if r.db.NewRecord(record) {
		record.SetDefaultName()
		r.Create(&record)
	}

	return &record
}
