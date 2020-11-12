package errcode

import (
	"fmt"
	"sort"
	"sync"

	"github.com/boxgo/box/pkg/insight"
	"github.com/gin-gonic/gin"
)

type (
	err struct {
		Code int32  `json:"code"`
		Msg  string `json:"msg"`
	}
	errs []err
)

var (
	locker        sync.Mutex
	registered    sync.Map
	registerError = func(key string) error { return fmt.Errorf("errcode %s is registered", key) }
)

func init() {
	insight.Get("/errors", func(ctx *gin.Context) {
		ctx.JSON(200, RegisteredErrors())
	})
}

func RegisteredErrors() []err {
	locker.Lock()
	defer locker.Unlock()

	var errs errs
	registered.Range(func(key, value interface{}) bool {
		if ecb, ok := value.(*ErrorCodeBuilder); ok {
			errs = append(errs, err{
				Code: ecb.Code(),
				Msg:  ecb.Message(),
			})
		}

		return true
	})

	sort.Sort(errs)

	return errs
}

func (es errs) Len() int {
	return len(es)
}

func (es errs) Less(i, j int) bool {
	return es[i].Code < es[j].Code
}

func (es errs) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}
