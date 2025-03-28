package server

import (
	"log/slog"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"github.com/metal-stack/go-dhcp-relay/config"
)

func Serve(log *slog.Logger, config *config.Config) error {
	listenAddress := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: dhcpv4.ServerPort,
	}

	handler := &handler{
		log:    log,
		config: config,
	}

	server, err := server4.NewServer("", listenAddress, handler.handle)
	if err != nil {
		return err
	}

	log.Info("starting dhcp relay", "listen address", listenAddress)
	err = server.Serve()
	if err != nil {
		return err
	}

	return nil
}
