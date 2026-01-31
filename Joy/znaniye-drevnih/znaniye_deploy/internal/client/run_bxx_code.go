package client

import (
	"context"
	"fmt"

	piston "github.com/milindmadhukar/go-piston"
	"go.uber.org/zap"
)

const (
	pistonLanguage = "cpp"
	pistonVersion  = "10.2.0"

	cppFileName = "main.cpp"
)

func (c *Client) RunBXXCode(ctx context.Context,
	bxxCode string,
	stdin string,
) (*piston.PistonExecution, error) {
	res, err := c.c.Execute(ctx,
		pistonLanguage,
		pistonVersion,
		[]piston.Code{
			*c.bxxLib,
			{
				Name:    cppFileName,
				Content: bxxCode,
			},
		},
		piston.CompileMemoryLimit(int(c.pistonMemoryLimit)),
		piston.RunMemoryLimit(int(c.pistonMemoryLimit)),
		piston.Stdin(stdin),
	)
	if err != nil {
		c.l.Error("piston execution error", zap.Error(err))
		return nil, fmt.Errorf("c.c.Execute: %w", err)
	}
	c.l.Info("piston executed",
		zap.Int("codeLen", len(bxxCode)),
		zap.Float64("compileMemory", res.Compile.Memory),
		zap.Float64("compileWallTime", res.Compile.WallTime),
		zap.String("compileMessage", res.Compile.Message),
		zap.String("compileStatus", res.Compile.Status),
		zap.Float64("runMemory", res.Run.Memory),
		zap.Float64("runWallTime", res.Run.WallTime),
		zap.String("runMessage", res.Run.Message),
		zap.String("runStatus", res.Run.Status),
	)
	return res, nil
}
