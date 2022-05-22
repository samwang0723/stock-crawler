// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package structuredlog

import (
	config "github.com/samwang0723/stock-crawler/configs"
	"github.com/sirupsen/logrus"

	logtest "github.com/sirupsen/logrus/hooks/test"
)

var (
	instance ILogger
)

type ILogger interface {
	RawLogger() *logrus.Logger
}

type structuredLogger struct {
	logger *logrus.Logger
}

func initialize(l ILogger) {
	instance = l
	instance.RawLogger().Info("initialized logger")
}

func Logger(cfg *config.Config) ILogger {
	if instance == nil {
		var level logrus.Level
		switch cfg.Log.Level {
		case "FATAL":
			level = logrus.FatalLevel
		case "INFO":
			level = logrus.InfoLevel
		case "WARN":
			level = logrus.WarnLevel
		case "ERROR":
			level = logrus.ErrorLevel
		default:
			level = logrus.DebugLevel
		}
		slog := &structuredLogger{
			logger: logrus.New(),
		}
		slog.logger.SetLevel(level)
		slog.logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})

		initialize(slog)
	}
	return instance
}

func NullLogger() ILogger {
	l, _ := logtest.NewNullLogger()
	initialize(&structuredLogger{
		logger: l,
	})
	return instance
}

func (l *structuredLogger) RawLogger() *logrus.Logger {
	return l.logger
}
