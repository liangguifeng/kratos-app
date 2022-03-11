package logging

import (
	"context"
	"github.com/liangguifeng/kratos-app/logging/applog"
)

const (
	// 日志打印时默认的堆栈跳过层数
	DefaultCallerSkip = applog.APPLOG_CALLER_SKIP + 1
)

type ErrLoggerContextIface interface {
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Warnw(ctx context.Context, err error)
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Errorw(ctx context.Context, err error)
	WithCommonFields(commonFields applog.Fields) ErrLoggerContextIface
	WithFields(fields applog.Fields) ErrLoggerContextIface
	WithCaller(skip int) ErrLoggerContextIface
}

type ErrLoggerContext struct {
	logger *applog.LoggerContext
}

func GetErrLogger(loggerTag string) (*ErrLoggerContext, error) {
	errLogger, err := applog.GetErrLogger(loggerTag)
	if err != nil {
		return nil, err
	}

	return &ErrLoggerContext{logger: errLogger}, nil
}

func (l *ErrLoggerContext) Warn(ctx context.Context, args ...interface{}) {
	l.logger.Warn(ctx, args...)
}

func (l *ErrLoggerContext) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Warnf(ctx, format, args...)
}

func (l *ErrLoggerContext) Warnw(ctx context.Context, err error) {
	l.logger.Warnw(ctx, err)
}

func (l *ErrLoggerContext) Error(ctx context.Context, args ...interface{}) {
	l.logger.Error(ctx, args...)
}

func (l *ErrLoggerContext) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Errorf(ctx, format, args...)
}

func (l *ErrLoggerContext) Errorw(ctx context.Context, err error) {
	l.logger.Errorw(ctx, err)
}

func (l *ErrLoggerContext) WithCommonFields(commonFields applog.Fields) ErrLoggerContextIface {
	ll := l.logger.WithCommonFields(commonFields)
	newLog := l.clone()
	newLog.logger = ll.(*applog.LoggerContext)
	return newLog
}

func (l *ErrLoggerContext) WithFields(fields applog.Fields) ErrLoggerContextIface {
	ll := l.logger.WithFields(fields)
	newLog := l.clone()
	newLog.logger = ll.(*applog.LoggerContext)
	return newLog
}

func (l *ErrLoggerContext) WithCaller(skip int) ErrLoggerContextIface {
	ll := l.logger.WithCaller(skip)
	newLog := l.clone()
	newLog.logger = ll.(*applog.LoggerContext)
	return newLog
}

func (l *ErrLoggerContext) clone() *ErrLoggerContext {
	newLog := *l
	return &newLog
}

type BusinessLoggerContextIface interface {
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Debugw(ctx context.Context, err error)
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Infow(ctx context.Context, err error)
	WithCommonFields(commonFields applog.Fields) BusinessLoggerContextIface
	WithFields(fields applog.Fields) BusinessLoggerContextIface
	WithCaller(skip int) BusinessLoggerContextIface
}

type BusinessLoggerContext struct {
	logger *applog.LoggerContext
}

func GetBusinessLogger(loggerTag string) (*BusinessLoggerContext, error) {
	businessLogger, err := applog.GetBusinessLogger(loggerTag)
	if err != nil {
		return nil, err
	}

	return &BusinessLoggerContext{logger: businessLogger}, nil
}

func (l *BusinessLoggerContext) Debug(ctx context.Context, args ...interface{}) {
	l.logger.Debug(ctx, args...)
}

func (l *BusinessLoggerContext) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debugf(ctx, format, args...)
}

func (l *BusinessLoggerContext) Debugw(ctx context.Context, err error) {
	l.logger.Debugw(ctx, err)
}

func (l *BusinessLoggerContext) Info(ctx context.Context, args ...interface{}) {
	l.logger.Info(ctx, args...)
}

func (l *BusinessLoggerContext) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logger.Infof(ctx, format, args...)
}

func (l *BusinessLoggerContext) Infow(ctx context.Context, err error) {
	l.logger.Infow(ctx, err)
}

func (l *BusinessLoggerContext) WithCommonFields(commonFields applog.Fields) BusinessLoggerContextIface {
	ll := l.logger.WithCommonFields(commonFields)
	newLog := l.clone()
	newLog.logger = ll.(*applog.LoggerContext)
	return newLog
}

func (l *BusinessLoggerContext) WithFields(fields applog.Fields) BusinessLoggerContextIface {
	ll := l.logger.WithFields(fields)
	newLog := l.clone()
	newLog.logger = ll.(*applog.LoggerContext)
	return newLog
}

func (l *BusinessLoggerContext) WithCaller(skip int) BusinessLoggerContextIface {
	ll := l.logger.WithCaller(skip + DefaultCallerSkip) // 外部不需要关心内部需要跳过几个层数
	newLog := l.clone()
	newLog.logger = ll.(*applog.LoggerContext)
	return newLog
}

func (l *BusinessLoggerContext) clone() *BusinessLoggerContext {
	newLog := *l
	return &newLog
}
