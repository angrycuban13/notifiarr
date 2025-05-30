package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	ping "github.com/prometheus-community/pro-bing"
)

// Custom errors.
var (
	ErrPingExpect    = errors.New("ping/icmp expect must contain three integers separated by two colons. ex: 3:2:500")
	ErrNoPingVal     = errors.New("ping or icmp 'check' must not be empty")
	ErrPingCountZero = errors.New("none of the ping expect values may be set to 0")
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

func (s *ServiceConfig) checkPingValues(icmp bool) error {
	if s.Value == "" {
		return ErrNoPingVal
	}

	if err := s.fillPingExpect(icmp); err != nil {
		return fmt.Errorf("ping expect format is count:min:interval where interval is in milliseconds "+
			" and min is the minimum packets that must be returned: %w", err)
	}

	if s.ping.count == 0 || s.ping.min > s.ping.count || s.ping.interval <= 0 {
		return fmt.Errorf("%w, ensure minimum(%d) is less than or equal to count(%d) and interval(%d) > 0",
			ErrPingExpect, s.ping.min, s.ping.count, s.ping.interval)
	}

	return nil
}

func (s *ServiceConfig) fillPingExpect(icmp bool) error {
	s.ping = &pingExpect{
		icmp: icmp,
	}

	splitStr := strings.Split(s.Expect, ":")
	if len(splitStr) != 3 { //nolint:mnd
		return ErrPingExpect
	}

	var err error
	if s.ping.count, err = strconv.Atoi(splitStr[0]); err != nil {
		return fmt.Errorf("invalid packet send count: %s: %w", splitStr[0], err)
	}

	if s.ping.min, err = strconv.Atoi(splitStr[1]); err != nil {
		return fmt.Errorf("invalid minimum packet receive count: %s: %w", splitStr[1], err)
	}

	if s.ping.interval, err = strconv.Atoi(splitStr[2]); err != nil {
		return fmt.Errorf("invalid packet send interval: %s: %w", splitStr[2], err)
	}

	return nil
}

func (s *ServiceConfig) checkPING() *result {
	pinger, err := ping.NewPinger(s.Value)
	if err != nil {
		return &result{
			state:  StateUnknown,
			output: &Output{str: "invalid ping value: " + err.Error()},
		}
	}

	pinger.SetPrivileged(s.ping.icmp)
	pinger.Timeout = s.Timeout.Duration
	pinger.Count = s.ping.count
	pinger.Interval = time.Duration(s.ping.interval) * time.Millisecond

	if err = pinger.Run(); err != nil { // blocks.
		return &result{
			state:  StateCritical,
			output: &Output{str: "error pinging service: " + err.Error()},
		}
	}

	stats := pinger.Statistics()
	state := StateOK
	msg := fmt.Sprintf("rcvd(%d) >= min(%d)", stats.PacketsRecv, s.ping.min)

	// Check if we received our minimum packet count responses.
	if stats.PacketsRecv < s.ping.min {
		state = StateCritical
		msg = fmt.Sprintf("rcvd(%d) < min(%d)", stats.PacketsRecv, s.ping.min)
	}

	return &result{
		state: state,
		output: &Output{str: fmt.Sprintf("(%s) pkts sent:%d, rcvd:%d, loss:%.01f, max:%s, avg:%s",
			msg, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
			stats.MaxRtt.Round(time.Millisecond), stats.AvgRtt.Round(time.Millisecond))},
	}
}
