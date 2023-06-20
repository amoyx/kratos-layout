package assembler

import (
	"github.com/jinzhu/copier"
)

type Assembler[T any, D any] interface {
	// DTOToDomain 领域模型转model
	DTOToDomain(*T) *D
	// DomainToDto model转领域模型
	DomainToDto(*D) *T
}

type BaseAssembler[T any, M any] struct {
}

func (b BaseAssembler[T, M]) DTOToDomain(t *T) *M {
	var m M
	err := copier.Copy(&m, t)
	if err != nil {
		return nil
	}
	return &m
}

func (b BaseAssembler[T, M]) DomainToDto(m *M) *T {
	var t T
	err := copier.Copy(&t, m)
	if err != nil {
		return nil
	}
	return &t
}
