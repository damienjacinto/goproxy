package utils

import (
    "os"
    "flag"
    "os/signal"
    "syscall"
    "github.com/rs/zerolog"
)

func InstallSignalHandler(stop chan struct{}) {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigs
        stop <- struct{}{}
        close(stop)
    }()
}

type Config struct {
    Debug bool
    Backend string
}

func InitLog() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func GetFlag() *Config {
    debug := flag.Bool("debug", false, "sets log level to debug")
    backend := flag.String("backend", "", "proxy backend")
    flag.Parse()

    if *debug {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }
    return &Config{Debug: *debug, Backend: *backend}
}

func GetEnv(key, fallback string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        value = fallback
    }
    return value
}
