package models

func (CategoryJasa) TableName() string {
	return "service_categories"
}

func (Service) TableName() string {
	return "services"
}
