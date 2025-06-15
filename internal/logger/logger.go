package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Set the global time field format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set the global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Output to console
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Logger = log.Logger.With().CallerWithSkipFrameCount(3).Logger()
}

func Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func Error(msg string) {
	log.Error().Msg(msg)
}

func Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}
