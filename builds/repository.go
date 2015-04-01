package builds

type Repository struct {
	builds map[string]*Build
}

func NewRepository() *Repository {
	return &Repository{
		builds: make(map[string]*Build),
	}
}

func (r *Repository) Find(key string) (*Build, bool) {
	output, ok := r.builds[key]
	return output, ok
}

func (r *Repository) Save(key string, build *Build) {
	r.builds[key] = build
}

func (r *Repository) Destroy(key string) {
	r.builds[key] = nil
}

func (r *Repository) Keys() []string {
	keys := make([]string, 0, len(r.builds))
	for k := range r.builds {
		keys = append(keys, k)
	}
	return keys
}

func (r *Repository) All() []*Build {
	all := make([]*Build, 0, len(r.builds))
	for _, build := range r.builds {
		all = append(all, build)
	}

	return all
}
