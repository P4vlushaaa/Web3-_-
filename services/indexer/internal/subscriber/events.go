package subscriber

import (
	"context"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient"
	"time"
	"web3-onlyfans/services/indexer/internal/utils"
)

type BlockSubscriber struct {
	cfg    *utils.Config
	logger utils.Logger
	rpc    *rpcclient.Client
}

func NewBlockSubscriber(cfg *utils.Config, logger utils.Logger) (*BlockSubscriber, error) {
	cli, err := rpcclient.New(context.Background(), cfg.NeoRPC, rpcclient.Options{})
	if err != nil {
		return nil, err
	}
	return &BlockSubscriber{
		cfg:    cfg,
		logger: logger,
		rpc:    cli,
	}, nil
}

func (s *BlockSubscriber) Start() {
	pollInterval := time.Duration(s.cfg.PollIntervalSeconds) * time.Second
	for {
		err := s.pollOnce()
		if err != nil {
			s.logger.Errorf("poll error: %v", err)
		}
		time.Sleep(pollInterval)
	}
}

func (s *BlockSubscriber) pollOnce() error {
	// Здесь логика: узнаём текущий блок, сверяемся с локальным бд,
	// проходимся по новым блокам, анализируем транзакции/нотификации.
	// Для упрощения покажем только "пример" — в реале будет много кода.

	height, err := s.rpc.GetBlockCount()
	if err != nil {
		return err
	}
	s.logger.Debugf("current block count: %d", height)
	// ... далее - обработка новых блоков (Tx, Notifications, Transfers).
	return nil
}
