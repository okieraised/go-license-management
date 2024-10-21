package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {
	viper.AddConfigPath("conf")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Info(fmt.Sprintf("error reading config file, %s", err))
	}
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)
}

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		slog.Info("app shutdown completed")
	}
}
