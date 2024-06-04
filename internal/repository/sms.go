package repository

import (
	"LinkMe/internal/repository/dao"
)

// SendCodeRepository 接口定义了异步 SMS 记录操作的相关方法
type SendVCodeRepository interface {
	CheckCode(mobile, smsID, vCode string)
	SendCode(mobile, smsID string)
}

// sendCodeRepository 实现了 AsyncSmsRepository 接口
type sendCodeRepository struct {
	dao dao.SendVCodeDAO
}

// NewAsyncSMSRepository 创建并返回一个新的 asyncSmsRepository 实例
func NewAsyncSMSRepository(dao dao.SendVCodeDAO) SendVCodeRepository {
	return &sendCodeRepository{
		dao: dao,
	}
}

func (s *sendCodeRepository) CheckCode(mobile, smsID, vCode string) {
	//todo
}

func (s *sendCodeRepository) SendCode(mobile, smsID string) {
	//todo
}
