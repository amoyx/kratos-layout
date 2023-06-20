package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/go-kratos/kratos-layout/internal/data/mapper"
)

type Greeter struct {
	Id         int64     `gorm:"primaryKey" json:"id"`                  // 数据库表主键
	Name       string    `gorm:"column:name" json:"name"`               // 姓名
	Sex        string    `gorm:"column:sex" json:"sex"`                 // 性别
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"` // 更新时间
}

func (g *Greeter) SetCreateTime() {
	g.CreateTime = time.Now()
}

func (g *Greeter) SetUpdateTime() {
	g.UpdateTime = time.Now()
}

func (g Greeter) TableName() string {
	return "greeter"
}

type GreeterMapper struct {
	mapper.DataMapper[Greeter]
}

func NewGreeterMapper(db *gorm.DB) *GreeterMapper {
	return &GreeterMapper{
		mapper.NewGormMapper[Greeter](db),
	}
}

// QueryBySex add user query method
func (g *GreeterMapper) QueryBySex(sex string) ([]*Greeter, error) {
	return nil, nil
}
