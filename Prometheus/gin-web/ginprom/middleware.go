package ginprom

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// 命名空间，用于查询记录的名称。名称组成，Namespace_Subsystem_Name
const namespace = "service"

var (
	// 标签名称
	labels = []string{"status", "endpoint", "method"}

	// 构造需要记录的数据结构

	// 服务启动时间
	uptime = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "uptime",
			Help:      "HTTP service uptime.",
		}, nil,
	)

	// http请求qps
	reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_request_count_total",
			Help:      "Total number of HTTP requests made.",
		}, labels,
	)

	// 响应时长，区间直方图，会生成count(请求记录条数)、sum(每次响应时长的总和)、bucket(每种区间的记录条数)三种数据
	reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
		}, labels,
	)

	// 请求大小，详情点图，会生成count(请求记录条数)、sum(每次请求大小的总和)
	reqSizeBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_request_size_bytes",
			Help:      "HTTP request sizes in bytes.",
		}, labels,
	)

	// 响应大小，详情点图，会生成count(请求记录条数)、sum(每次请求大小的总和)
	respSizeBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP request sizes in bytes.",
		}, labels,
	)
)

// init registers the prometheus metrics
// 初始化数据类型
func init() {
	prometheus.MustRegister(uptime, reqCount, reqDuration, reqSizeBytes, respSizeBytes)
	go recordUptime()
}

// recordUptime increases service uptime per second.
// 每秒记录计数器加一。服务的启动时间=当前时间-计数器总数。已运行时长=计数器总数
func recordUptime() {
	for range time.Tick(time.Second) {
		uptime.WithLabelValues().Inc()
	}
}

// calcRequestSize returns the size of request object.
// 计算请求大小
func calcRequestSize(r *http.Request) float64 {
	size := 0
	if r.URL != nil {
		size = len(r.URL.String())
	}

	size += len(r.Method)
	size += len(r.Proto)

	for name, values := range r.Header {
		size += len(name)
		for _, value := range values {
			size += len(value)
		}
	}
	size += len(r.Host)

	// r.Form and r.MultipartForm are assumed to be included in r.URL.
	if r.ContentLength != -1 {
		size += int(r.ContentLength)
	}
	return float64(size)
}

type RequestLabelMappingFn func(c *gin.Context) string

// PromOpts represents the Prometheus middleware Options.
// It is used for filtering labels by regex.
// 参数信息的参数，用来过滤标签
type PromOpts struct {
	ExcludeRegexStatus     string // 过滤状态
	ExcludeRegexEndpoint   string // 过滤url path
	ExcludeRegexMethod     string // 过滤请求方式
	EndpointLabelMappingFn RequestLabelMappingFn // endpoint标签装饰函数
}

// NewDefaultOpts return the default ProOpts
// 初始化默认的参数信息
func NewDefaultOpts() *PromOpts {
	return &PromOpts{
		EndpointLabelMappingFn: func(c *gin.Context) string {
			//by default do nothing, return URL as is
			// 没有处理，直接返回请求path
			return c.Request.URL.Path
		},
	}
}

// checkLabel returns the match result of labels.
// Return true if regex-pattern compiles failed.
// 检查标签
func (po *PromOpts) checkLabel(label, pattern string) bool {
	if pattern == "" {
		return true
	}
	// 是否匹配标签
	matched, err := regexp.MatchString(pattern, label)
	if err != nil {
		return true
	}
	return !matched
}

// PromMiddleware returns a gin.HandlerFunc for exporting some Web metrics
// prometheus 中间件，用于记录各种数据
func PromMiddleware(promOpts *PromOpts) gin.HandlerFunc {
	// make sure promOpts is not nil
	// 初始化prometheus的参数信息
	if promOpts == nil {
		promOpts = NewDefaultOpts()
	}

	// make sure EndpointLabelMappingFn is callable
	// 确保标签装饰可用，默认直接返回url path
	if promOpts.EndpointLabelMappingFn == nil {
		promOpts.EndpointLabelMappingFn = func(c *gin.Context) string {
			return c.Request.URL.Path
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		// 等待内部执行完毕，包含内部的中间件。这里等待请求返回
		c.Next()

		// 返回状态status
		status := fmt.Sprintf("%d", c.Writer.Status())
		// endpoint标签装饰器
		endpoint := promOpts.EndpointLabelMappingFn(c)
		// 请求方式
		method := c.Request.Method
		// 标签集合
		lvs := []string{status, endpoint, method}

		// 校验标签
		isOk := promOpts.checkLabel(status, promOpts.ExcludeRegexStatus) &&
			promOpts.checkLabel(endpoint, promOpts.ExcludeRegexEndpoint) &&
			promOpts.checkLabel(method, promOpts.ExcludeRegexMethod)

		if !isOk {
			return
		}
		// no response content will return -1
		// 没有返回内容结果为 -1
		respSize := c.Writer.Size()
		if respSize < 0 {
			respSize = 0
		}
		// qps加一
		reqCount.WithLabelValues(lvs...).Inc()
		// 响应时长
		reqDuration.WithLabelValues(lvs...).Observe(time.Since(start).Seconds())
		// 请求大小
		reqSizeBytes.WithLabelValues(lvs...).Observe(calcRequestSize(c.Request))
		// 响应大小
		respSizeBytes.WithLabelValues(lvs...).Observe(float64(respSize))
	}
}

// PromHandler wrappers the standard http.Handler to gin.HandlerFunc
// 提供给 prometheus 拉取数据的装饰器，将http 转 为 gin的方法
func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}