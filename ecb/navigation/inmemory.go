package navigation

import "context"

type StaticRepository struct {
	Items []Record
}

func NewStaticRepository(records []Record) *StaticRepository {
	return &StaticRepository{Items: records}
}

func (r *StaticRepository) FetchTree(ctx context.Context) ([]Record, error) {
	return r.Items, nil
}
