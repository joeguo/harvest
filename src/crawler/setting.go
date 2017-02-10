package crawler

import (
	"time"
	"github.com/PuerkitoBio/purell"
)

const (
	DefaultDepth                                        = 1
	DefaultRequests                                     = 5
	DefaultMaxUrls                                      = 1000
	DefaultUserAgent                                    = `Mozilla/5.0 (Windows NT 6.1; rv:15.0) Gecko/20120716 Firefox/15.0a2`
	DefaultNormalizationFlags purell.NormalizationFlags = purell.FlagsAllGreedy
	DefaultBufferFactor                                 = 100
	DefaultTimeOut                                      = 15 * time.Minute
)

type Setting struct {
	Depth                 int
	Requests              int
	BufferFactor          int
	MaxUrls               int
	TimeOut               time.Duration
	UserAgent             string
	URLNormalizationFlags purell.NormalizationFlags

}

var DefaultSetting = Setting{Depth:DefaultDepth, Requests:DefaultRequests, BufferFactor:DefaultBufferFactor, MaxUrls:DefaultMaxUrls, TimeOut:DefaultTimeOut, UserAgent:DefaultUserAgent, URLNormalizationFlags:DefaultNormalizationFlags}
