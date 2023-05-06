

package cb

import (
	"context"
	"errors"

	"github.com/cloudwego/kitex/pkg/endpoint"
)

type ctxKey int

const (
	ctxPass ctxKey = iota
)

var (
	noPass  = "NoPass"
	pass    = "Pass"
	errFail = errors.New("you shall not pass")
)

func FailMW(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		if status, _ := ctx.Value(ctxPass).(string); status == noPass {
			return errFail
		}
		return nil
	}
}
