package storage

type Repository struct {
	URLStorer
}

func NewRepository(db URLStorer) *Repository {
	return &Repository{
		URLStorer: db,
	}
}
