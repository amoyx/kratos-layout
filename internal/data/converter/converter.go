package converter

import (
	"github.com/jinzhu/copier"
)

type Converter[D any, M any] interface {
	// DomainToModel 领域模型转model
	DomainToModel(*D) *M
	// ModelToDomain model转领域模型
	ModelToDomain(*M) *D
	// ApplyDomainToModel 根据领域模型更新model
	ApplyDomainToModel(*D, *M)
}

type BaseConverter[D any, M any] struct {
}

func (b BaseConverter[D, M]) ApplyDomainToModel(d *D, m *M) {
	copier.CopyWithOption(m, d, copier.Option{IgnoreEmpty: true, DeepCopy: false})
}

func (b BaseConverter[D, M]) DomainToModel(d *D) *M {
	var m M
	err := copier.Copy(&m, d)
	if err != nil {
		return nil
	}
	return &m
}

func (b BaseConverter[D, M]) ModelToDomain(m *M) *D {
	var d D
	err := copier.Copy(&d, m)
	if err != nil {
		return nil
	}
	return &d
}
