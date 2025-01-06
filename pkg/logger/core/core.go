package core

import (
	"crypto/md5"
	"fmt"

	"go.uber.org/zap/zapcore"
)

var (
	_ zapcore.Core = (*MaskCore)(nil)
)

type MaskCore struct {
	splitLen int
	zapcore.LevelEnabler
	maskers *Maskers
	enc     zapcore.Encoder
	out     zapcore.WriteSyncer
}

func NewMaskCore(splitLen int, rules MaskRules, lv zapcore.LevelEnabler, enc zapcore.Encoder, out zapcore.WriteSyncer) zapcore.Core {
	if splitLen <= 0 {
		splitLen = 5 * 1024
	}

	return &MaskCore{
		splitLen:     splitLen,
		LevelEnabler: lv,
		maskers:      NewMaskers(rules),
		enc:          enc,
		out:          out,
	}
}

func (c *MaskCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.enc, fields)
	return clone
}

func (c *MaskCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	return ce
}

func (c *MaskCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	data := c.maskers.Mask(buf.Bytes())

	if dataLen := len(data); dataLen > c.splitLen {
		h := md5.New()
		h.Write(data)
		logId := []byte(fmt.Sprintf("[logId=%x]", h.Sum(nil)))

		for c.splitLen < len(data) {
			_, err = c.out.Write(append(logId, data[0:c.splitLen:c.splitLen]...))
			_, err = c.out.Write([]byte("\n"))

			data = data[c.splitLen:]
		}

		if len(data) > 0 {
			_, err = c.out.Write(append(logId, data...))
		}
	} else {
		_, err = c.out.Write(data)
	}

	buf.Free()
	if err != nil {
		return err
	}
	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		_ = c.Sync()
	}
	return nil
}

func (c *MaskCore) Sync() error {
	return c.out.Sync()
}

func (c *MaskCore) clone() *MaskCore {
	return &MaskCore{
		splitLen:     c.splitLen,
		LevelEnabler: c.LevelEnabler,
		maskers:      c.maskers,
		enc:          c.enc.Clone(),
		out:          c.out,
	}
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
