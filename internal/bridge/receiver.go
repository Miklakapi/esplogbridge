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

	buf := make([]byte, 65535)

	for {
		_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))

		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				select {
				case <-ctx.Done():
					return nil
				default:
					continue
				}
			}
			return fmt.Errorf("read: %w", err)
		}

		packet := strings.TrimSpace(string(buf[:n]))
		if packet == "" {
			continue
		}

		ipStr := remoteIPString(addr)
		deviceID, ok := cfg.Devices[ipStr]
		if !ok {
			continue
		}

		clean := stripSyslogPrefix(packet)
		clean = normalizeWhitespace(clean)
		if clean == "" {
			continue
		}

		fmt.Println(clean)

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
