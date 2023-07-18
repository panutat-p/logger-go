package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type MyEncoder struct {
	zapcore.Encoder
}

func (m *MyEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	filtered := make([]zapcore.Field, 0, len(fields))
	for _, field := range fields {
		if field.Key == "password" {
			continue
		}
		filtered = append(filtered, field)
	}
	return m.Encoder.EncodeEntry(entry, filtered)
}

func CustomEncoder(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
	encoder := zapcore.NewJSONEncoder(config)
	return &MyEncoder{encoder}, nil
}

func main() {
	err := zap.RegisterEncoder("custom", CustomEncoder)
	if err != nil {
		fmt.Println("🔴 Zap encoder err:", err)
		os.Exit(1)
	}
	config := zap.NewDevelopmentConfig()
	config.Encoding = "custom"
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		MessageKey:     "message",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		LevelKey:       "level",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	logger, err := config.Build()
	if err != nil {
		fmt.Println("🔴 Zap config err:", err)
		os.Exit(1)
	}
	logger.Info("Hello",
		zap.String("animal", "🐵"),
		zap.Int("weight", 42),
		zap.Int("password", 1234),
	)
}
