package db

// Tabler - ...
type Tabler interface {
	TableName() string
}

// Order - Mined order info holder table model
type Order struct {
	Number              uint64       `gorm:"column:number;type:bigint;not null;unique;index:,sort:asc"`
	Timestamp           uint64       `gorm:"column:timestamp;type:bigint;not null;index:,sort:asc"`
	Aggressor 	        string       `gorm:"column:aggressor;type:char(8);not null"`
	Price   			      float64      `gorm:"column:price;type:float(8);not null"`
	Quantity            uint64       `gorm:"column:quantity;type:bigint;not null"`
}

// TableName - Overriding default table name
func (Order) TableName() string {
	return "orders"
}

// SimilarTo - Checking whether two orders are exactly similar or not
func (b *Order) SimilarTo(_b *Order) bool {
	return b.Number == _b.Number &&
		b.Number == _b.Number &&
		b.Timestamp == _b.Timestamp &&
		b.Aggressor == _b.Aggressor &&
		b.Price == _b.Price &&
		b.Quantity == _b.Quantity
}
