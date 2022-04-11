package tightening_device

import (
	"fmt"
	"net/http"
)

func (s *Service) handlerSocketSelector(data interface{}) {
	if data == nil {
		return
	}

	req := data.(SocketSelectorReq)
	//step, err := s.storageService.GetStepByCodeAndWorkorderID(req.StepCode, req.WorkorderID)
	//if err != nil {
	//	s.diag.Error("Step Not Found", err)
	//	return
	//}
	//
	//consume, err := s.storageService.GetConsumeBySeqInStep(step, int(req.Sequence))
	//if err != nil {
	//	s.diag.Error("Consume Not Found", err)
	//	return
	//}
	//
	//if consume.Socket.IOSNInput == "" {
	//	s.diag.Debug("Socket Not Configured")
	//	return
	//}
	//
	//req.Socket = consume.Socket

	switch req.Type {
	case SocketSelectorTrigger:
		if err := s.triggerSocketSelector(req.PSetSet); err != nil {
			s.diag.Error("triggerSocketSelector Failed", err)
			return
		}

	case SocketSelectorClear:
		if err := s.triggerSocketSelectorClear(); err != nil {
			s.diag.Error("triggerSocketSelector Failed", err)
			return
		}
	}

}

func (s *Service) triggerSocketSelector(req PSetSet) error {

	if !s.config().SocketSelector.Enable {
		return nil
	}

	r := s.httpClient.R().SetBody(req)
	url := s.config().SocketSelector.Endpoint
	resp, err := r.Put(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("triggerSocketSelector Failed: %d", resp.StatusCode())
	}

	return nil
}

func (s *Service) triggerSocketSelectorClear() error {
	if !s.config().SocketSelector.Enable {
		return nil
	}

	r := s.httpClient.R()
	url := s.config().SocketSelector.Endpoint
	resp, err := r.Delete(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("triggerSocketSelectorClear Failed: %d", resp.StatusCode())
	}

	return nil
}
