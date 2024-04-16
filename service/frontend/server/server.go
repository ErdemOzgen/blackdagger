package server

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/logger/tag"
	"github.com/ErdemOzgen/blackdagger/service/frontend/restapi"
	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"

	pkgmiddleware "github.com/ErdemOzgen/blackdagger/service/frontend/middleware"
	"github.com/ErdemOzgen/blackdagger/service/frontend/restapi/operations"

	"github.com/go-chi/chi/v5"
)

type BasicAuth struct {
	Username string
	Password string
}

type AuthToken struct {
	Token string
}

type Params struct {
	Host      string
	Port      int
	BasicAuth *BasicAuth
	AuthToken *AuthToken
	TLS       *config.TLS
	Logger    logger.Logger
	Handlers  []New
	AssetsFS  fs.FS
}

type Server struct {
	host      string
	port      int
	basicAuth *BasicAuth
	authToken *AuthToken
	tls       *config.TLS
	logger    logger.Logger
	server    *restapi.Server
	handlers  []New
	assets    fs.FS
}

type New interface {
	Configure(api *operations.BlackdaggerAPI)
}

func NewServer(params Params) *Server {
	return &Server{
		host:      params.Host,
		port:      params.Port,
		basicAuth: params.BasicAuth,
		authToken: params.AuthToken,
		tls:       params.TLS,
		logger:    params.Logger,
		handlers:  params.Handlers,
		assets:    params.AssetsFS,
	}
}

func (svr *Server) Shutdown() {

	if svr.server == nil {
		return
	}
	err := svr.KillGotty()
	if err != nil {
		svr.logger.Warn("GOTTY shutdown", tag.Error(err))
	}
	err = svr.server.Shutdown()
	if err != nil {
		svr.logger.Warn("Server shutdown", tag.Error(err))
	}

}
func (svr *Server) KillGotty() error {
	// Command to list all gotty processes
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		svr.logger.Warn("Failed to list processes: " + err.Error())
		return err
	}

	// Parse the output to find gotty processes
	processes := strings.Split(out.String(), "\n")
	for _, process := range processes {
		if strings.Contains(process, "gotty") {
			fields := strings.Fields(process)
			if len(fields) > 1 {
				pid, err := strconv.Atoi(fields[1]) // PID is usually the second field in 'ps aux' output
				if err != nil {
					svr.logger.Warn("Failed to parse PID for process: " + process + ", error: " + err.Error())
					continue
				}
				// Kill the process by PID
				killCmd := exec.Command("kill", strconv.Itoa(pid))
				err = killCmd.Run()
				if err != nil {
					svr.logger.Warn("Failed to kill gotty process with PID " + strconv.Itoa(pid) + ": " + err.Error())
				} else {
					svr.logger.Info("Killed gotty process with PID " + strconv.Itoa(pid))
				}
			}
		}
	}

	return nil
}

func (svr *Server) Serve(ctx context.Context) (err error) {
	middlewareOptions := &pkgmiddleware.Options{
		Handler: svr.defaultRoutes(chi.NewRouter()),
	}
	if svr.authToken != nil {
		middlewareOptions.AuthToken = &pkgmiddleware.AuthToken{
			Token: svr.authToken.Token,
		}
	}
	if svr.basicAuth != nil {
		middlewareOptions.AuthBasic = &pkgmiddleware.AuthBasic{
			Username: svr.basicAuth.Username,
			Password: svr.basicAuth.Password,
		}
	}
	pkgmiddleware.Setup(middlewareOptions)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		svr.logger.Error("failed to load API spec", tag.Error(err))
		return err
	}
	api := operations.NewBlackdaggerAPI(swaggerSpec)
	for _, h := range svr.handlers {
		h.Configure(api)
	}

	svr.server = restapi.NewServer(api)
	defer svr.Shutdown()

	svr.server.Host = svr.host
	svr.server.Port = svr.port
	svr.server.ConfigureAPI()

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(ctx)

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Trigger graceful shutdown
		err := svr.server.Shutdown()
		if err != nil {
			svr.logger.Error("server shutdown error", tag.Error(err))
		}
		serverStopCtx()
	}()

	if svr.tls != nil {
		svr.server.TLSCertificate = flags.Filename(svr.tls.CertFile)
		svr.server.TLSCertificateKey = flags.Filename(svr.tls.KeyFile)
		svr.server.EnabledListeners = []string{"https"}
		svr.server.TLSHost = svr.host
		svr.server.TLSPort = svr.port
	}

	// Run the server
	err = svr.server.Serve()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		svr.logger.Error("server error", tag.Error(err))
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	svr.logger.Info("server closed")

	return nil
}
