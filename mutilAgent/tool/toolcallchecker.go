package tool

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/pkg/errors"
	"io"
)

func ToolCallChecker(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
	defer sr.Close()
	for {
		msg, err := sr.Recv()
		if err != nil {
			fmt.Println("ioEOF", err.Error())
			if errors.Is(err, io.EOF) {
				break
			}
			return false, err
		}
		if len(msg.ToolCalls) > 0 {
			return true, nil
		}
	}
	return false, nil
}
