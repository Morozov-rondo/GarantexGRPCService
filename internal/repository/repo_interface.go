package repository

import (
	"context"

	"garantexGRPC/models"
)

//go:generate mockgen -source=repo_interface.go -destination=./mocks/mock_repo_interface.go -package=mocks
type Storager interface {
	SaveRates(ctx context.Context, rates models.Rates) error
}
