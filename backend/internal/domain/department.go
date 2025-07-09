package domain

type Department struct {
	ID   int16  `gorm:"column:department_id;primaryKey"`
	Name string `gorm:"column:name"`
}

func (d *Department) DepartmentID() int16 {
	return d.ID
}

func (d *Department) DepartmentName() string {
	return d.Name
}
