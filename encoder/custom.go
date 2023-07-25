package encoder

type customEncoder struct {
	enc zapcore.Encoder
}

func newCustomEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "" // Omit timestamps from the output
	return &customEncoder{
		enc: zapcore.NewJSONEncoder(encoderConfig), // Use JSON encoder as a base
	}
}

func (c *customEncoder) Clone() zapcore.Encoder {
	return &customEncoder{
		enc: c.enc.Clone(),
	}
}

func (c *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer, error) {
	// Customize the log output format here
	// For simplicity, we'll print the log level and message only
	line := fmt.Sprintf("[%s] %s\n", entry.Level.String(), entry.Message)
	return &buffer{
		str: line,
	}, nil
}

func (c *customEncoder) EncodeEntrySync(entry zapcore.Entry, fields []zapcore.Field) (*buffer, error) {
	// Since we don't buffer logs, this is the same as EncodeEntry
	return c.EncodeEntry(entry, fields)
}

func (c *customEncoder) EncodeKey(key string) {
	c.enc.EncodeKey(key)
}

func (c *customEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	return c.enc.AddArray(key, marshaler)
}

func (c *customEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	return c.enc.AddObject(key, marshaler)
}

func (c *customEncoder) AddBinary(key string, value []byte) {
	c.enc.AddBinary(key, value)
}

func (c *customEncoder) AddByteString(key string, value []byte) {
	c.enc.AddByteString(key, value)
}

func (c *customEncoder) AddBool(key string, value bool) {
	c.enc.AddBool(key, value)
}

func (c *customEncoder) AddComplex128(key string, value complex128) {
	c.enc.AddComplex128(key, value)
}

func (c *customEncoder) AddComplex64(key string, value complex64) {
	c.enc.AddComplex64(key, value)
}

func (c *customEncoder) AddDuration(key string, value time.Duration) {
	c.enc.AddDuration(key, value)
}

func (c *customEncoder) AddFloat64(key string, value float64) {
	c.enc.AddFloat64(key, value)
}

func (c *customEncoder) AddFloat32(key string, value float32) {
	c.enc.AddFloat32(key, value)
}

func (c *customEncoder) AddInt(key string, value int) {
	c.enc.AddInt(key, value)
}

func (c *customEncoder) AddInt64(key string, value int64) {
	c.enc.AddInt64(key, value)
}

func (c *customEncoder) AddInt32(key string, value int32) {
	c.enc.AddInt32(key, value)
}

func (c *customEncoder) AddInt16(key string, value int16) {
	c.enc.AddInt16(key, value)
}

func (c *customEncoder) AddInt8(key string, value int8) {
	c.enc.AddInt8(key, value)
}

func (c *customEncoder) AddString(key string, value string) {
	c.enc.AddString(key, value)
}

func (c *customEncoder) AddTime(key string, value time.Time) {
	c.enc.AddTime(key, value)
}

func (c *customEncoder) AddUint(key string, value uint) {
	c.enc.AddUint(key, value)
}

func (c *customEncoder) AddUint64(key string, value uint64) {
	c.enc.AddUint64(key, value)
}

func (c *customEncoder) AddUint32(key string, value uint32) {
	c.enc.AddUint32(key, value)
}

func (c *customEncoder) AddUint16(key string, value uint16) {
	c.enc.AddUint16(key, value)
}

func (c *customEncoder) AddUint8(key string, value uint8) {
	c.enc.AddUint8(key, value)
}

func (c *customEncoder) AddUintptr(key string, value uintptr) {
	c.enc.AddUintptr(key, value)
}

func (c *customEncoder) AddReflected(key string, value interface{}) error {
	return c.enc.AddReflected(key, value)
}

func (c *customEncoder) OpenNamespace(key string) {
	c.enc.OpenNamespace(key)
}

func (c *customEncoder) ClonePool() zapcore.Encoder {
	return &customEncoder{
		enc: c.enc.ClonePool(),
	}
}
