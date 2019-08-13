package controller

type CreateTable struct {
	Order string `form:"Order" json:"Order" required:"false"`
}

type UpdateTable struct {
	Order string `form:"Order" json:"Order" required:"false"`
}
