package ports

import "styl-monolith/internal/users/core/domain"

type UserPort interface {
	CreateUser(user domain.User) (domain.User, error)
	ListUsers(page, size int) ([]domain.User, int, error)
	GetUserByID(id uint) (domain.User, error)
	UpdateUser(user domain.User) (domain.User, error)
	DeleteUser(id uint) error
	GetUserByEmail(email string) (domain.User, error)
	GetUserByPhone(phone string) (domain.User, error)
	GetUserByUsername(username string) (domain.User, error)
	PartialUserUpdate(id uint, user domain.User) (domain.User, error)
	UpdateUserRole(user domain.User) (domain.User, error)
}

type RolePort interface {
	CreateRole(role domain.Role) (domain.Role, error)
	ListRoles(page, size int) ([]domain.Role, int, error)
	GetRoleByID(id uint) (domain.Role, error)
	UpdateRole(role domain.Role) (domain.Role, error)
	DeleteRole(id uint) error
	GetRoleByName(name string) (domain.Role, error)
	PartialRoleUpdate(id uint, role domain.Role) (domain.Role, error)
}

type UserReportRepository interface {
	CreateReport(report domain.Report) (domain.Report, error)
	ListReports(page, size int) ([]domain.Report, int, error)
	GetReportByID(id uint) (domain.Report, error)
	DeleteReport(id uint) error
	UpdateReport(report domain.Report) (domain.Report, error)
	GetReportByUserID(id uint) ([]domain.Report, error)
	PartialReportUpdate(id uint, report domain.Report) (domain.Report, error)
}
