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

package logger

func Debug(args ...interface{}) {
	instance.RawLogger().Debug(args...)
}

func Debugf(s string, args ...interface{}) {
	instance.RawLogger().Debugf(s, args...)
}

func Info(args ...interface{}) {
	instance.RawLogger().Info(args...)
}

func Infof(s string, args ...interface{}) {
	instance.RawLogger().Infof(s, args...)
}

func Warn(args ...interface{}) {
	instance.RawLogger().Warn(args...)
}

func Warnf(s string, args ...interface{}) {
	instance.RawLogger().Warnf(s, args...)
}

func Fatal(args ...interface{}) {
	instance.RawLogger().Fatal(args...)
}

func Fatalf(s string, args ...interface{}) {
	instance.RawLogger().Fatalf(s, args...)
}

func Error(args ...interface{}) {
	instance.RawLogger().Error(args...)
}

func Errorf(s string, args ...interface{}) {
	instance.RawLogger().Errorf(s, args...)
}

func Flush() {
	instance.Flush()
}
