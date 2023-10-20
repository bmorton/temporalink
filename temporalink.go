package temporalink

import (
	"context"
	"fmt"

	"github.com/pborman/uuid"
	"github.com/temporalio/cli/server"
	cliconfig "github.com/temporalio/tctl-kit/pkg/config"
	uiserver "github.com/temporalio/ui-server/v2/server"
	uiconfig "github.com/temporalio/ui-server/v2/server/config"
	uiserveroptions "github.com/temporalio/ui-server/v2/server/server_options"
	_ "go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/temporal"
)

type EmbeddedTemporal struct {
	IP         string
	ServerPort int
	UIPort     int

	interruptChan chan interface{}
	server        *server.Server
}

func NewEmbeddedTemporal(ip string, serverPort int, uiPort int) (*EmbeddedTemporal, error) {
	et := &EmbeddedTemporal{
		IP:            ip,
		ServerPort:    serverPort,
		UIPort:        uiPort,
		interruptChan: make(chan interface{}, 1),
	}

	return et, et.setup()
}

func (et *EmbeddedTemporal) setup() error {
	opts := []server.ServerOption{
		server.WithDynamicPorts(),
		server.WithFrontendPort(et.ServerPort),
		server.WithFrontendIP(et.IP),
		server.WithNamespaces("default"),
		server.WithUpstreamOptions(
			temporal.InterruptOn(et.interruptChan),
		),
		server.WithBaseConfig(&config.Config{}),
	}

	frontendAddr := fmt.Sprintf("%s:%d", et.IP, et.ServerPort)

	uiBaseCfg := &uiconfig.Config{
		Host:                et.IP,
		Port:                et.UIPort,
		TemporalGRPCAddress: frontendAddr,
		EnableUI:            true,
		EnableOpenAPI:       true,
	}
	opts = append(opts, server.WithUI(uiserver.NewServer(uiserveroptions.WithConfigProvider(uiBaseCfg))))

	opts = append(opts, server.WithPersistenceDisabled())
	if clusterCfg, err := cliconfig.NewConfig("temporalio", "version-info"); err == nil {
		defaultEnv := "default"
		clusterIDKey := "cluster-id"

		clusterID, _ := clusterCfg.EnvProperty(defaultEnv, clusterIDKey)

		if clusterID == "" {
			// fallback to generating a new cluster Id in case of errors or empty value
			clusterID = uuid.New()
			cErr := clusterCfg.SetEnvProperty(defaultEnv, clusterIDKey, clusterID)
			if cErr != nil {
				return fmt.Errorf("failed to set cluster ID: %w", cErr)
			}
		}

		opts = append(opts, server.WithCustomClusterID(clusterID))
	}

	var err error
	et.server, err = server.NewServer(opts...)

	return err
}

func (et *EmbeddedTemporal) Start(ctx context.Context) error {
	go func() {
		if doneChan := ctx.Done(); doneChan != nil {
			s := <-doneChan
			et.interruptChan <- s
		} else {
			s := <-temporal.InterruptCh()
			et.interruptChan <- s
		}
	}()

	return et.server.Start()
}
