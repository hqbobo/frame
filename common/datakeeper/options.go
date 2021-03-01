package datakeeper

import (
	"time"
)

type options struct {
	//负责定时reload
	hour, min, sec int
	of             bool
	//负责定时器reload
	timer          time.Duration
	tf             bool
	//负责定时update
	uptimer          time.Duration
	upf             bool
}

// Option overrides behavior of Connect.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithTriggerTimer 按照定时器时间Reload
func WithTriggerTimer(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.timer = t
		o.tf = true
	})
}

// WithTriggerDaily 每天固定时间Reload
func WithTriggerDaily(hour, min, second int) Option {
	return optionFunc(func(o *options) {
		o.hour = hour
		o.min = min
		o.sec = second
		o.of = true
	})
}


// WithUpdateTimer 按照定时器时间Update
func WithUpdateTimer(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.upf = true
		o.uptimer = t
	})
}