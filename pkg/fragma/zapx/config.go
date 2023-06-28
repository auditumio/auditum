package zapx

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func defaultConfig() zap.Config {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig = defaultEncoderConfig()

	return conf
}

func defaultEncoderConfig() zapcore.EncoderConfig {
	conf := zap.NewProductionEncoderConfig()
	conf.TimeKey = "time"
	conf.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	return conf
}

func getConfig(format string, level string) (zap.Config, error) {
	encoderConfig, encoding, err := getEncoder(format)
	if err != nil {
		return zap.Config{}, err
	}

	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return zap.Config{}, err
	}

	zap.AddCaller()

	conf := zap.NewProductionConfig()
	conf.Level = atomicLevel
	conf.Encoding = encoding
	conf.EncoderConfig = encoderConfig

	if encoding == zapFormatConsole {
		conf.DisableCaller = true
	}

	return conf, nil
}

func getEncoder(format string) (zapcore.EncoderConfig, string, error) {
	switch format {
	case FormatJSON:
		ec, e := jsonEncoder()
		return ec, e, nil
	case FormatText:
		ec, e := textEncoder()
		return ec, e, nil
	default:
		return zapcore.EncoderConfig{}, "", fmt.Errorf("unknown format: %v", format)
	}
}

func jsonEncoder() (zapcore.EncoderConfig, string) {
	return defaultEncoderConfig(), zapFormatJSON
}

func textEncoder() (zapcore.EncoderConfig, string) {
	conf := defaultEncoderConfig()
	conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	conf.EncodeTime = zapcore.RFC3339TimeEncoder

	return conf, zapFormatConsole
}
