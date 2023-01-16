package log

import (
	"sync"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// zapLogger 是Logger接口的具体实现。它底层封装了zap.Logger
type zapLogger struct {
	z *zap.Logger
}

var (
	mu sync.Mutex
	//std 定义了默认的全局Logger
	std = NewLogger(NewOptions())
)

// 确保zapLogger 实现了Logger 接口。以下变量赋值，可以使错误在编译期间被发现
var _ Logger = &zapLogger{}

func NewLogger(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}
	// 将文本格式的日志级别，例如 info 转换为 zapcore.Level 类型以供后面使用
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		// 如果指定了非法的日志级别，则默认使用 info 级别
		zapLevel = zapcore.InfoLevel
	}
	// 创建一个默认的 encoder 配置
	encoderConfig := zap.NewProductionEncoderConfig()
	//自定义 MessageKey 为message,message 语义更明确
	encoderConfig.MessageKey = "message"
	// 自定义 TimeKey 为 timestamp,timestamp 语义更明确
	encoderConfig.TimeKey = "timestamp"
	//指定时间序列化函数，将时间序列化为 `2006-01-02 15:04:05.000`格式，更易读
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	//指定time.Duration 序列化函数，将time.Duration 序列化为经过的毫秒数的浮点数
	//毫秒数比默认的秒数更精确
	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}
	//创建构建 zap.Logger 需要的配置
	cfg := &zap.Config{
		// 指定日志级别
		Level: zap.NewAtomicLevelAt(zapLevel),
		//是否在日志中显示调用日志所在的文件行号，例如：`"caller":"miniblog/miniblog.go:79"`
		DisableCaller: opts.DisableCaller,
		//是否禁止在 panic 及以上级别打印堆栈信息
		DisableStacktrace: opts.DisableStacktrace,
		//指定日志显示格式，可选值：console ,json
		Encoding:      opts.Format,
		EncoderConfig: encoderConfig,
		// 指定日志输出位置
		OutputPaths: opts.OutputPaths,
		// 设置zap内部错误输出位置
		ErrorOutputPaths: []string{"stderr"},
	}
	// 使用cfg 创建 *zap.Logger 对象
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{
		z: z,
	}
	//把标准库中的 log.Logger 的 info 级别的输出重定向到zap.Logger
	zap.RedirectStdLog(z)
	return logger
}

// 使用指定的选项初始化Logger

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = NewLogger(opts)
}

// Logger 定义了miniblog项目的日志接口，该接口只包含了支持的日志记录方法
type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Sync()
}

// Debugw 输出debug 级别的日志
func Debugw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Debugw(msg, keysAndValues...)
}
func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Debugw(msg, keysAndValues...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Infow(msg, keysAndValues...)
}
func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Infow(msg, keysAndValues...)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Warnw(msg, keysAndValues...)
}
func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Warnw(msg, keysAndValues...)
}
func Errorw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Errorw(msg, keysAndValues...)
}
func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Errorw(msg, keysAndValues...)
}
func Panicw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Panicw(msg, keysAndValues...)
}

// Panicw 输出 panic级别的日志
func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Fatalw(msg, keysAndValues...)
}
func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Fatalw(msg, keysAndValues...)
}

// Sync 调用底层zap.Logger的Sync方法，将缓存中的日志刷新到磁盘文件中.主程序需要在退出前调用Sync.
func Sync() {
	std.Sync()
}
func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}
