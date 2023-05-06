

package cb

import (
	"context"
	"log"

	"github.com/bytedance/gopkg/cloud/circuitbreaker"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
)

func ChangeHandler(key string, oldState, newState circuitbreaker.State, m circuitbreaker.Metricer) {
	log.Printf("circuitbreaker status change, old: %v, new: %v\n", oldState, newState)
}

func GetKey(ctx context.Context, request interface{}) (key string, enabled bool) {
	return "1234", true
}

func GetErrorType(ctx context.Context, request, response interface{}, err error) circuitbreak.ErrorType {
	if err != nil {
		return circuitbreak.TypeFailure
	}
	return circuitbreak.TypeSuccess
}

func DecorateError(ctx context.Context, request interface{}, err error) error {
	return err
}
