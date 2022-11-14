package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

// Custom errors.
var (
	ErrPingExpect    = fmt.Errorf("ping/icmp expect must contain three integers separated by two colons. ex: 3:2:500")
	ErrNoPingVal     = fmt.Errorf("ping or icmp 'check' must not be empty")
	ErrPingCountZero = fmt.Errorf("none of the ping expect values may be set to 0")
)

/*
 * These all run once at startup to fill our check data.
 * The service Lock is acquired before running any of this code.
 */

// pingExpect is setup for each 'ping'/'icmp' service from input data on initialization.
type pingExpect struct {
	count    int
	min      int
	interval int
	icmp     bool
}

func (s *Service) checkPingValues(icmp bool) error {
	if s.Value == "" {
		return ErrNoPingVal
	}

	if err := s.fillPingExpect(icmp); err != nil {
		return fmt.Errorf("ping expect format is count:min:interval where interval is in milliseconds "+
			" and min is the minimum packets that must be returned: %w", err)
	}

	if s.svc.ping.count == 0 || s.svc.ping.min > s.svc.ping.count || s.svc.ping.interval <= 0 {
		return fmt.Errorf("%w, ensure minimum(%d) is less than or equal to count(%d) and interval(%d) > 0",
			ErrPingExpect, s.svc.ping.min, s.svc.ping.count, s.svc.ping.interval)
	}

	return nil
}

func (s *Service) fillPingExpect(icmp bool) (err error) {
	s.svc.ping = &pingExpect{
		icmp: icmp,
	}

	splitStr := strings.Split(s.Expect, ":")
	if len(splitStr) != 3 { //nolint:gomnd
		return ErrPingExpect
	}

	if s.svc.ping.count, err = strconv.Atoi(splitStr[0]); err != nil {
		return fmt.Errorf("invalid packet send count: %s: %w", splitStr[0], err)
	}

	if s.svc.ping.min, err = strconv.Atoi(splitStr[1]); err != nil {
		return fmt.Errorf("invalid minimum packet receive count: %s: %w", splitStr[1], err)
	}

	if s.svc.ping.interval, err = strconv.Atoi(splitStr[2]); err != nil {
		return fmt.Errorf("invalid packet send interval: %s: %w", splitStr[2], err)
	}

	return nil
}

func (s *Service) checkPING() *result {
	pinger, err := ping.NewPinger(s.Value)
	if err != nil {
		return &result{
			state:  StateUnknown,
			output: "invalid ping value: " + err.Error(),
		}
	}

	pinger.SetPrivileged(s.svc.ping.icmp)
	pinger.Timeout = s.Timeout.Duration
	pinger.Count = s.svc.ping.count
	pinger.Interval = time.Duration(s.svc.ping.interval) * time.Millisecond

	if err = pinger.Run(); err != nil { // blocks.
		return &result{
			state:  StateCritical,
			output: "error pinging service: " + err.Error(),
		}
	}

	stats := pinger.Statistics()
	state := StateOK
	msg := fmt.Sprintf("rcvd(%d) >= min(%d)", stats.PacketsRecv, s.svc.ping.min)

	// Check if we received our minimum packet count responses.
	if stats.PacketsRecv < s.svc.ping.min {
		state = StateCritical
		msg = fmt.Sprintf("rcvd(%d) < min(%d)", stats.PacketsRecv, s.svc.ping.min)
	}

	return &result{
		state: state,
		output: fmt.Sprintf("(%s) pkts sent:%d, rcvd:%d, loss:%.01f, max:%s, avg:%s",
			msg, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
			stats.MaxRtt.Round(time.Millisecond), stats.AvgRtt.Round(time.Millisecond)),
	}
}
