package bridge

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Miklakapi/esplogbridge/internal/config"
)

type Event struct {
	Timestamp time.Time
	DeviceID  string
	RawLine   string
}

func runUDPReceiver(ctx context.Context, cfg config.Config, out chan Event) error {
	conn, err := net.ListenPacket("udp", cfg.Listen)
	if err != nil {
		return fmt.Errorf("listen %s: %w", cfg.Listen, err)
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	buf := make([]byte, 65535)

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("read: %w", err)
		}

		ipStr := remoteIPString(addr)
		deviceID, ok := cfg.Devices[ipStr]
		if !ok {
			continue
		}

		packet := string(buf[:n])
		if packet == "" {
			continue
		}

		clean := stripSyslogPrefix(packet)
		clean = normalizeWhitespace(clean)
		clean = trimEspHomePrefix(clean)
		if clean == "" {
			continue
		}

		ev := Event{
			Timestamp: time.Now().UTC(),
			DeviceID:  deviceID,
			RawLine:   clean,
		}

		pushDropOldest(out, ev)

	}
}

func remoteIPString(addr net.Addr) string {
	switch a := addr.(type) {
	case *net.UDPAddr:
		if a.IP == nil {
			return ""
		}
		if v4 := a.IP.To4(); v4 != nil {
			return v4.String()
		}
		return a.IP.String()
	default:
		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			return ""
		}
		ip := net.ParseIP(host)
		if ip == nil {
			return ""
		}
		if v4 := ip.To4(); v4 != nil {
			return v4.String()
		}
		return ip.String()
	}
}

func pushDropOldest(ch chan Event, ev Event) {
	select {
	case ch <- ev:
		return
	default:
		select {
		case <-ch:
		default:
		}
		select {
		case ch <- ev:
		default:
		}
	}
}

func stripSyslogPrefix(s string) string {
	s = strings.TrimLeft(s, "\r\n\t ")
	if len(s) == 0 || s[0] != '<' {
		return s
	}

	end := strings.IndexByte(s, '>')
	if end <= 1 {
		return s
	}

	return strings.TrimLeft(s[end+1:], " \t")
}

func normalizeWhitespace(s string) string {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return ""
	}
	return strings.Join(fields, " ")
}

func trimEspHomePrefix(s string) string {
	idx := strings.IndexByte(s, '[')
	if idx <= 0 {
		return s
	}
	return strings.TrimSpace(s[idx:])
}
