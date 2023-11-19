package domain

// import (
// 	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
// 	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
// )

// // type TicketSearchCondition interface {
// // 	GetPageNo() int
// // 	GetPageSize() int
// // }

// // type ticketSearchCondition struct {
// // 	PageNo   int
// // 	PageSize int
// // }

// // func NewTicketSearchCondition(pageNo, pageSize int) (TicketSearchCondition, error) {
// // 	m := &ticketSearchCondition{
// // 		PageNo:   pageNo,
// // 		PageSize: pageSize,
// // 	}

// // 	if err := libdomain.Validator.Struct(m); err != nil {
// // 		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
// // 	}

// // 	return m, nil
// // }

// // func (p *ticketSearchCondition) GetPageNo() int {
// // 	return p.PageNo
// // }

// // func (p *ticketSearchCondition) GetPageSize() int {
// // 	return p.PageSize
// // }

// // type TicketSearchResult interface {
// // 	GetTotalCount() int
// // 	GetResults() []TicketModel
// // }

// // type ticketSearchResult struct {
// // 	TotalCount int
// // 	Results    []TicketModel
// // }

// // func NewTicketSearchResult(totalCount int, results []TicketModel) (TicketSearchResult, error) {
// // 	m := &ticketSearchResult{
// // 		TotalCount: totalCount,
// // 		Results:    results,
// // 	}

// // 	if err := libdomain.Validator.Struct(m); err != nil {
// // 		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
// // 	}

// // 	return m, nil
// // }
// // func (m *ticketSearchResult) GetTotalCount() int {
// // 	return m.TotalCount
// // }

// // func (m *ticketSearchResult) GetResults() []TicketModel {
// // 	return m.Results
// // }
