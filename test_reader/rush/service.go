package reader

import (
	"fmt"
	"github.com/ebfe/scard"
	"go.uber.org/atomic"
	"time"
)

const (
	SearchItv = 1 * time.Second
	LenCmdOK  = 4
)

type Service struct {
	configValue   atomic.Value
	ctx           *scard.Context

}

func NewService(c Config) *Service {

	s := &Service{
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) Model() interface{} {
	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	go s.search()

	return nil
}

func (s *Service) Close() error {
	if s.ctx != nil {
		err := s.ctx.Release()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) waitUntilCardPresent(ctx *scard.Context, readers []string) (int, error) {
	rs := make([]scard.ReaderState, len(readers))
	for i := range rs {
		rs[i].Reader = readers[i]
		rs[i].CurrentState = scard.StateUnaware
	}

	for {
		for i := range rs {
			if rs[i].EventState&scard.StatePresent != 0 {
				return i, nil
			}
			rs[i].CurrentState = rs[i].EventState
		}
		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return -1, err
		}
	}
}

func (s *Service) DeviceType() string {
	return "reader"
}

func (s *Service) Status() string {
	return "online"
}

func (s *Service) Config() interface{} {
	return nil
}

func (s *Service) Data() interface{} {
	return nil
}

func (s *Service) search() {
	var err error

	for {
		s.ctx, err = scard.EstablishContext()
		if err == nil {
			break
		} else {
			time.Sleep(SearchItv)
		}
	}

	for {
		// List available readers
		ctx, err := scard.EstablishContext()
		if err != nil {
			time.Sleep(SearchItv)
			continue
		}

		readers, err := ctx.ListReaders()
		if err != nil {
			//s.diag.Debug("reader lost")
			_ = ctx.Release()
			time.Sleep(SearchItv)
			continue
		}

		for {
			if len(readers) > 0 {
				index, err := s.waitUntilCardPresent(ctx, readers)
				if err != nil {
					break
				}

				// connect to card
				card, err := ctx.Connect(readers[index], scard.ShareExclusive, scard.ProtocolAny)
				if err != nil {
					time.Sleep(SearchItv)
					continue
				}

				for {
					var cmd = []byte{0xff, 0xca, 0x00, 0x00, 0x00} // SELECT uid

					_, err := card.Status()
					if err != nil {
						_ = card.Disconnect(scard.ResetCard)
						break
					}

					rsp, err := card.Transmit(cmd)
					if err != nil {
						// card lost
						//s.diag.Debug("card lost")
						_ = card.Disconnect(scard.ResetCard)
						break
					}

					uid := fmt.Sprintf("%x", string(rsp))
					if len(uid) > LenCmdOK {
						uid = uid[0 : len(uid)-LenCmdOK]
					}

					fmt.Printf("uid:%s\n", uid)

					// ws notify

					_ = card.Disconnect(scard.ResetCard)

					time.Sleep(SearchItv)
				}
			}
		}
	}
}
