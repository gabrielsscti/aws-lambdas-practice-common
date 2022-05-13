package movie

const MOVIE_TABLE_ENV = "MOVIE_TABLE_ENV"

type Repository interface {
	GetByID(id string) ([]byte, error)
	GetAll() ([]byte, error)
	Insert(movie []byte) error
	Delete(id string) error
}
