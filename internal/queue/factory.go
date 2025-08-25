package queue

import (
	"rpa-middleware/internal/config"
	"rpa-middleware/internal/interfaces"

	"github.com/sirupsen/logrus"
)

// NewQueueManager 根据配置创建队列管理器
func NewQueueManager(cfg *config.Config, logger *logrus.Logger) (interfaces.QueueManager, error) {
	switch cfg.Queue.Type {
	case "redis":
		logger.Info("Using Redis queue manager")
		return NewRedisQueueManager(
			cfg.Redis.Addr,
			cfg.Redis.Password,
			cfg.Redis.DB,
			logger,
		)
	case "memory", "":
		fallthrough
	default:
		logger.Info("Using memory queue manager")
		return NewMemoryQueueManager(), nil
	}
}
