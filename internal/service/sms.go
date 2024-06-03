package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"errors"
	"github.com/ecodeclub/ekit"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	statsWindowSize       = 60                     // 统计最近 60 秒的数据
	responseTimeThreshold = 500 * time.Millisecond // 响应时间阈值
	errorRateThreshold    = 0.05                   // 错误率阈值
	checkDuration         = time.Minute            // 检查统计数据的时间范围
	syncTrafficPercentage = 0.01                   // 保留的同步流量比例
)

// SendCodeService 定义了发送验证码的服务接口
type SendCodeService interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}

// sendCodeService 实现了 SendCodeService 接口
type sendCodeService struct {
	repo   repository.SendCodeRepository
	l      *zap.Logger
	client *sms.Client
	// 统计数据
	mu            sync.Mutex
	stats         []statsEntry
	currentSecond int
	appId         *string
	signName      *string
}

type statsEntry struct {
	timestamp       time.Time
	requests        int
	errorRequests   int
	responseTimeSum time.Duration
}

// NewSendCodeService 创建并返回一个新的 sendCodeService 实例
func NewSendCodeService(repo repository.SendCodeRepository, l *zap.Logger, client *sms.Client, appId *string, signName *string) SendCodeService {
	s := &sendCodeService{
		repo:     repo,
		l:        l,
		stats:    make([]statsEntry, statsWindowSize),
		client:   client,
		appId:    appId,
		signName: signName,
	}
	go s.StartAsyncCycle()
	return s
}

// StartAsyncCycle 启动异步发送短信的循环
func (s *sendCodeService) StartAsyncCycle() {
	time.Sleep(time.Second * 3) // 初始化延迟
	for {
		s.AsyncSend()
	}
}

// AsyncSend 异步发送短信
func (s *sendCodeService) AsyncSend() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	as, err := s.repo.PreemptWaitingSMS(ctx)
	if err != nil {
		s.handleAsyncSendError(err)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.Send(ctx, as.TplId, as.Args, as.Numbers...)
	if err != nil {
		s.l.Error("Failed to send asynchronous short messages", zap.Error(err))
	}
	res := err == nil
	if er := s.repo.ReportScheduleResult(ctx, as.Id, res); err != nil {
		s.l.Error("Asynchronous sending of SMS messages succeeded, but marking the database failed", zap.Error(er))
	}
}

// handleAsyncSendError 处理异步发送中的错误
func (s *sendCodeService) handleAsyncSendError(err error) {
	if errors.Is(err, repository.ErrWaitingSMSNotFound) {
		time.Sleep(time.Second)
	} else {
		s.l.Error("Failed to preempt the asynchronous SMS sending task", zap.Error(err))
		time.Sleep(time.Second)
	}
}

// Send 发送短信并记录统计数据
func (s *sendCodeService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	start := time.Now()
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = ekit.ToPtr(tplId)
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)
	response, err := s.client.SendSms(request)
	if err != nil {
		s.l.Error("Failed to send SMS", zap.Error(err))
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			continue
		}
		status := *statusPtr
		if status.Code == nil || *status.Code != "Ok" {
			s.l.Error("Failed to send SMS", zap.Error(err))
			return err
		}
	}
	responseTime := time.Since(start)
	s.recordStats(true, responseTime)
	if s.needAsync() {
		return s.repo.Add(ctx, domain.AsyncSms{
			TplId:    tplId,
			Args:     args,
			Numbers:  numbers,
			RetryMax: 3,
		})
	}
	return nil
}

// recordStats 记录统计数据
func (s *sendCodeService) recordStats(success bool, responseTime time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().Unix()
	currentIdx := int(now) % statsWindowSize
	if currentIdx != s.currentSecond {
		s.stats[currentIdx] = statsEntry{
			timestamp:       time.Now(),
			requests:        0,
			errorRequests:   0,
			responseTimeSum: 0,
		}
		s.currentSecond = currentIdx
	}
	entry := &s.stats[currentIdx]
	entry.requests++
	if !success {
		entry.errorRequests++
	}
	entry.responseTimeSum += responseTime
}

// getStatistics 获取统计数据
func (s *sendCodeService) getStatistics(duration time.Duration) (totalRequests int, errorRequests int, totalResponseTime time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-duration)
	for i := 0; i < statsWindowSize; i++ {
		entry := s.stats[i]
		if entry.timestamp.After(cutoff) {
			totalRequests += entry.requests
			errorRequests += entry.errorRequests
			totalResponseTime += entry.responseTimeSum
		}
	}
	return totalRequests, errorRequests, totalResponseTime
}

// needAsync 判断是否需要使用异步发送
func (s *sendCodeService) needAsync() bool {
	totalRequests, errorRequests, totalResponseTime := s.getStatistics(checkDuration)
	if totalRequests == 0 {
		return false
	}
	averageResponseTime := totalResponseTime / time.Duration(totalRequests)
	errorRate := float64(errorRequests) / float64(totalRequests)
	return averageResponseTime > responseTimeThreshold || errorRate > errorRateThreshold || float64(totalRequests)*syncTrafficPercentage < 1
}
