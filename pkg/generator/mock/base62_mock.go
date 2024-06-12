package mock

import (
	"github.com/stretchr/testify/mock"
)

type Generator struct {
	mock.Mock
}

func (g *Generator) GenerateShortURLCode() string {
	args := g.Called()
	return args.String(0)
}
