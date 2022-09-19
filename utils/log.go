package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

func init() {
	//writeSyncer := getLogWriter()
	//encoder := getEncoder()
	//
	//core := zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)
	//Logger = zap.New(core)
	Logger, _ = zap.NewProduction()
	defer Logger.Sync()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.log")

	return zapcore.AddSync(file)
}
