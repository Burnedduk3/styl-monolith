package postgres

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"styl-monolith/internal/users/adapters/postgres/models"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/ports"
	"styl-monolith/pkg/errorhandler"
)

type UserReportRepository struct {
	log  *logrus.Logger
	conn *gorm.DB
}

func NewUserReportRepository(log *logrus.Logger, conn *gorm.DB) ports.UserReportRepository {
	return &UserReportRepository{log: log, conn: conn}
}

// CreateReport creates a new report in the database.
func (r *UserReportRepository) CreateReport(report domain.Report) (domain.Report, error) {
	rReport := models.NewPostgresReportFromDomainReport(report)
	result := r.conn.Create(&rReport)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.Report{}, domainError
	}
	return rReport.ToReportDomain(), nil
}

// ListReports retrieves a paginated list of reports from the database.
func (r *UserReportRepository) ListReports(page, size int) ([]domain.Report, int, error) {
	var reports []models.Report
	var dReports []domain.Report
	offset := (page - 1) * size
	var totalCount int64
	r.conn.Model(&models.Report{}).Count(&totalCount)
	result := r.conn.Limit(size).Offset(offset).Find(&reports)
	// Calculate total pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error)
		return nil, 0, domainError
	}
	for _, report := range reports {
		dReports = append(dReports, report.ToReportDomain())
	}
	return dReports, totalPages, nil
}

// GetReportByID retrieves a report from the database using its ID.
func (r *UserReportRepository) GetReportByID(id uint) (domain.Report, error) {
	var report models.Report
	result := r.conn.First(&report, id)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrReportNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrReportNotFound), id),
			result.Error)
		return domain.Report{}, domainError
	}
	return report.ToReportDomain(), nil
}

// DeleteReport deletes a report from the database using its ID.
func (r *UserReportRepository) DeleteReport(id uint) error {
	result := r.conn.Delete(&models.Report{}, id)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrReportNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrReportNotFound), id),
			nil,
		)
	}
	return nil
}

// UpdateReport updates an existing report in the database.
func (r *UserReportRepository) UpdateReport(report domain.Report) (domain.Report, error) {
	rReport := models.NewPostgresReportFromDomainReport(report)
	result := r.conn.Model(&models.Report{}).Where("id = ?", report.ID).Save(&rReport)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.Report{}, domainError
	}
	return rReport.ToReportDomain(), nil
}

// GetReportByUserID retrieves all reports from the database by their associated User ID.
func (r *UserReportRepository) GetReportByUserID(userID uint) ([]domain.Report, error) {
	var reports []models.Report
	var dReports []domain.Report
	result := r.conn.Where("user_id = ?", userID).Find(&reports)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error)
	}
	for _, report := range reports {
		dReports = append(dReports, report.ToReportDomain())
	}
	return dReports, nil
}

// PartialReportUpdate updates specific fields of a report identified by its ID.
func (r *UserReportRepository) PartialReportUpdate(id uint, updatedReport domain.Report) (domain.Report, error) {
	rReport := models.NewPostgresReportFromDomainReport(updatedReport)
	result := r.conn.Model(&models.Report{}).Where("id = ?", id).Updates(rReport)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.Report{}, domainError
	}
	return rReport.ToReportDomain(), nil
}
