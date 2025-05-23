package services

type ItemsService interface {
	GetItemsForDashboard(flag string, page int, pageSize int) []string
}
