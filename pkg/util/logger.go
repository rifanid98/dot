package util

import (
	"context"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"os"
)

type Field map[string]any

type Logger struct {
	logger *zap.SugaredLogger
}

func NewLogger() *Logger {
	godotenv.Load()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	envs := []string{"LOCAL", "local"}
	if slices.Contains(envs, os.Getenv("LOGGER_MODE")) {
		log, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}

	log.With()
	defer log.Sync()

	return &Logger{logger: log.Sugar()}
}

func (z *Logger) Warn(ctx context.Context, format string, args ...any) {
	z.logger.Warnf(getTrackerId(ctx)+"::"+format, args...)
}

func (z *Logger) Info(ctx context.Context, format string, args ...any) {
	z.logger.Infof(getTrackerId(ctx)+"::"+format, args...)
}

func (z *Logger) Error(ctx context.Context, format string, args ...any) {
	z.logger.Errorf(getTrackerId(ctx)+"::"+format, args...)
}

func (z *Logger) Fatalln(ctx context.Context, args ...any) {
	z.logger.Fatal(args...)
}

func (z *Logger) WarnWf(ctx context.Context, format string, fields Field) {
	z.logger.With(getFields(bindFields(ctx, fields))...).Warnf(format)
}

func (z *Logger) InfoWf(ctx context.Context, format string, fields Field) {
	z.logger.With(getFields(bindFields(ctx, fields))...).Infof(format)
}

func (z *Logger) ErrorWf(ctx context.Context, format string, fields Field) {
	z.logger.With(getFields(bindFields(ctx, fields))...).Errorf(format)
}

func (z *Logger) FatallnWf(ctx context.Context, fields Field, args ...any) {
	z.logger.With(getFields(bindFields(ctx, fields))...).Fatal(args...)
}

func (z *Logger) WithFields(keyValues Field) *Logger {
	log := z.logger.With(getFields(keyValues)...)
	return &Logger{logger: log}
}

func (z *Logger) WithError(err error) *Logger {
	var log = z.logger.With(err.Error())
	return &Logger{logger: log}
}

func getFields(fields map[string]any) []any {
	var f = make([]any, 0)
	for index, field := range fields {
		f = append(f, index)
		f = append(f, field)
	}
	return f
}

func bindFields(ctx context.Context, fields map[string]any) map[string]any {
	trackerId := getTrackerId(ctx)
	if fields == nil {
		fields = map[string]any{}
	}
	fields["tracker_id"] = trackerId
	return fields
}

func getTrackerId(ctx context.Context) string {
	trackerId := ctx.Value("tracker_id")
	if trackerId == nil {
		return ""
	}
	return trackerId.(string)
}
