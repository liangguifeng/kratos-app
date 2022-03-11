package helper

import (
	"context"
	"github.com/liangguifeng/kratos-app/logging"
)

const (
	err_tag = "err"
	bus_tag = "business"
)

func NewLogger() (*logger, error) {
	errLogger, err := logging.GetErrLogger(err_tag)
	if err != nil {
		return nil, err
	}

	businessLogger, err := logging.GetBusinessLogger(bus_tag)
	if err != nil {
		return nil, err
	}

	l := make(map[string]interface{})
	l[err_tag] = errLogger
	l[bus_tag] = businessLogger
	return &logger{logS: l}, nil
}

type logger struct {
	logS map[string]interface{}
}

func (d *logger) GetErrLogger() logging.ErrLoggerContextIface {
	return d.logS[err_tag].(logging.ErrLoggerContextIface)
}

func (d *logger) GetBusinessLogger() logging.BusinessLoggerContextIface {
	return d.logS[bus_tag].(logging.BusinessLoggerContextIface)
}

func (d *logger) Debug(ctx context.Context, args ...interface{}) {
	d.GetBusinessLogger().Debug(ctx, args...)
}

func (d *logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	d.GetBusinessLogger().Debugf(ctx, format, args...)
}

func (d *logger) Debugw(ctx context.Context, err error) {
	d.GetBusinessLogger().Debugw(ctx, err)
}

func (d *logger) Info(ctx context.Context, args ...interface{}) {
	d.GetBusinessLogger().Info(ctx, args...)
}

func (d *logger) Infof(ctx context.Context, format string, args ...interface{}) {
	d.GetBusinessLogger().Infof(ctx, format, args...)
}

func (d *logger) Infow(ctx context.Context, err error) {
	d.GetBusinessLogger().Infow(ctx, err)
}

func (d *logger) Warn(ctx context.Context, args ...interface{}) {
	d.GetErrLogger().Warn(ctx, args...)
}

func (d *logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	d.GetErrLogger().Warnf(ctx, format, args...)
}

func (d *logger) Warnw(ctx context.Context, err error) {
	d.GetErrLogger().Warnw(ctx, err)
}

func (d *logger) Error(ctx context.Context, args ...interface{}) {
	d.GetErrLogger().Error(ctx, args...)
}

func (d *logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	d.GetErrLogger().Errorf(ctx, format, args...)
}

func (d *logger) Errorw(ctx context.Context, err error) {
	d.GetErrLogger().Errorw(ctx, err)
}
