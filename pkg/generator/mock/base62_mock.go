package mock

import (
	"github.com/stretchr/testify/mock"
)

type Generator struct {
	mock.Mock
}

func (g *Generator) GenerateShortURLCode(length int) string {
	args := g.Called(length)
	return args.String(0)
}
