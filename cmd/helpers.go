package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/alecthomas/kingpin"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	zapopt "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type commonOptions struct {
	MetricAddress         string
	MetricPort            uint16
	ManagerLeaderElection bool
	Debug                 bool
}

func NewCommonOptions(cmd *kingpin.Application) *commonOptions {
	opt := &commonOptions{}

	cmd.Flag("manager-leader-election", "Enable manager leader election").Default("false").BoolVar(&opt.ManagerLeaderElection)
	cmd.Flag("debug", "Enable debug logging").Default("false").BoolVar(&opt.Debug)

	return opt
}

func (opt *commonOptions) WithMetrics(cmd *kingpin.Application) *commonOptions {
	cmd.Flag("metrics-address", "Address to bind HTTP metrics listener").Default("127.0.0.1").StringVar(&opt.MetricAddress)
	cmd.Flag("metrics-port", "Port to bind HTTP metrics listener").Default("9525").Uint16Var(&opt.MetricPort)

	return opt
}

func (opt *commonOptions) Logger() logr.Logger {
	// While debugging, it may be useful to provide debug log lines that
	// include sensitive or large payloads.
	logLevel := zapcore.InfoLevel
	if opt.Debug {
		logLevel = zapcore.DebugLevel
	}

	logger := zap.New(
		zap.Encoder(zaplogfmt.NewEncoder(zapcore.EncoderConfig{
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			TimeKey:       "ts",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			EncodeTime:    zapcore.RFC3339NanoTimeEncoder,
		})),
		zap.WriteTo(os.Stderr),
		zap.Level(logLevel),
		zap.RawZapOpts(zapopt.AddCaller()),
	)
	ctrl.SetLogger(logger)

	return logger
}

// Set via compiler flags
var (
	Version   = "dev"
	Commit    = "none"
	Date      = "unknown"
	GoVersion = runtime.Version()
)

func VersionStanza() string {
	return fmt.Sprintf(
		"Version: %v\nGit SHA: %v\nGo Version: %v\nGo OS/Arch: %v/%v\nBuilt at: %v",
		Version, Commit, GoVersion, runtime.GOOS, runtime.GOARCH, Date,
	)
}

var (
	BuildInfo = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "theatre_build_info",
			Help: "A metric with a constant '1' value labeled by version, commit, goversion from which %s was built, and the goos and goarch for the build.",
			ConstLabels: prometheus.Labels{
				"version":   Version,
				"commit":    Commit,
				"goversion": GoVersion,
				"goos":      runtime.GOOS,
				"goarch":    runtime.GOARCH,
			},
		},
		func() float64 { return 1 },
	)
)
