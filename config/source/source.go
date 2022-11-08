// Package source is the interface for sources
package source

import (
	"crypto/md5"
	"fmt"
	"time"
)

type (
	// Source is the source from which config is loaded
	Source interface {
		Read() (*ChangeSet, error)
		Watch() (Watcher, error)
		String() string
	}

	Watcher interface {
		Next() (*ChangeSet, error)
		Stop() error
	}

	// ChangeSet represents a set of changes from a source
	ChangeSet struct {
		Data      []byte
		Checksum  string
		Format    string
		Source    string
		Timestamp time.Time
	}
)

// Sum returns the md5 checksum of the ChangeSet data
func (c *ChangeSet) Sum() string {
	h := md5.New()
	h.Write(c.Data)
	c.Checksum = fmt.Sprintf("%x", h.Sum(nil))

	return c.Checksum
}
