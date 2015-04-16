package streams

import "errors"

var ErrNotFound = errors.New("Build not found")

type Repository struct {
	streams map[string]*BuildStream
}

func NewRepository() *Repository {
	return &Repository{
		streams: make(map[string]*BuildStream),
	}
}

func (r *Repository) Create(stream *BuildStream) {
	r.streams[stream.BuildID] = stream
}

func (r *Repository) Find(key string) (*BuildStream, error) {
	stream, ok := r.streams[key]
	if !ok {
		return &BuildStream{}, ErrNotFound
	}

	return stream, nil
}

func (r *Repository) Destroy(key string) {
	delete(r.streams, key)
}
