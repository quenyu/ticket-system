package domain

type Department struct {
	ID   int16  `gorm:"column:department_id;primaryKey"`
	Name string `gorm:"column:name"`
}
