package core

import (
	"go.uber.org/zap/zapcore"
)

var (
	_ zapcore.Core = (*MaskCore)(nil)
)

type MaskCore struct {
	zapcore.LevelEnabler
	maskers *Maskers
	enc     zapcore.Encoder
	out     zapcore.WriteSyncer
}

func NewMaskCore(rules MaskRules, lv zapcore.LevelEnabler, enc zapcore.Encoder, out zapcore.WriteSyncer) zapcore.Core {
	return &MaskCore{
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

	_, err = c.out.Write(data)
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
