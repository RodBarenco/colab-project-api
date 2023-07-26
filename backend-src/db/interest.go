package db

//yet to change it ro uint32,  gorm create many-to-may relation with big int id,
type Interest struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null"`
}
