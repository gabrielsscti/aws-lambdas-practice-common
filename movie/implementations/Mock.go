package implementations

import (
	"encoding/json"
	"github.com/gabrielsscti/aws-lambdas-practice-common"
	"github.com/gabrielsscti/aws-lambdas-practice-common/movie"
	"net/http"
	"strconv"
)

type MockMovies struct{}

func (m MockMovies) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

var movies = []movie.Movie{
	{
		ID:   "1",
		Name: "Avengers",
	},
	{
		ID:   "2",
		Name: "Ant-Man",
	},
	{
		ID:   "3",
		Name: "Thor",
	},
	{
		ID:   "4",
		Name: "Hulk",
	},
	{
		ID:   "5",
		Name: "Doctor Strange",
	},
}

func CreateNewMockMovies() MockMovies {
	return MockMovies{}
}

func (m MockMovies) GetByID(id string) ([]byte, error) {
	_id, err := strconv.Atoi(id)
	if err != nil {
		return nil, common.ResponseError{
			StatusCode: http.StatusBadRequest,
			Body:       "ID has to be an integer",
			Info:       err.Error(),
		}
	}
	_id--
	if (_id < 0) || (_id > len(movies)-1) {
		return nil, common.ResponseError{
			StatusCode: http.StatusNotFound,
			Body:       "ID out of range",
			Info:       "outside slice range",
		}
	}

	return json.Marshal(movies[_id])
}

func (m MockMovies) GetAll() ([]byte, error) {
	return json.Marshal(movies)
}

func (m MockMovies) Insert(obj []byte) error {
	var movieObj movie.Movie
	err := json.Unmarshal(obj, &movieObj)
	if err != nil {
		return common.ResponseError{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid movie structure",
			Info:       err.Error(),
		}
	}
	movies = append(movies, movieObj)
	return nil
}
