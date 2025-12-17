// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package pipeline

import (
	"chatwiki/internal/pkg/lib_define"
	"reflect"
	"runtime"

	"github.com/zhimaAi/go_tools/logs"
)

type PipeResult bool

const PipeStop PipeResult = true      //终止后面的pipe执行
const PipeContinue PipeResult = false //继续执行后面的pipe

type Pipe[In, Out any] func(in *In, out *Out) PipeResult

type Pipeline[In, Out any] struct {
	In    *In
	Out   *Out
	pipes []Pipe[In, Out]
}

func (pipeline *Pipeline[In, Out]) Pipe(pipe Pipe[In, Out]) *Pipeline[In, Out] {
	pipeline.pipes = append(pipeline.pipes, pipe)
	return pipeline
}

func (pipeline *Pipeline[In, Out]) Process() (result PipeResult) {
	defer func() {
		pipeline.pipes = make([]Pipe[In, Out], 0) //清空pipes
	}()
	for _, pipe := range pipeline.pipes {
		if lib_define.IsDev {
			logs.Debug(`pipe:%s`, runtime.FuncForPC(reflect.ValueOf(pipe).Pointer()).Name())
		}
		if pipe != nil && pipe(pipeline.In, pipeline.Out) {
			return PipeStop
		}
	}
	return PipeContinue
}

func NewPipeline[In, Out any](in *In, out *Out) *Pipeline[In, Out] {
	return &Pipeline[In, Out]{
		In: in, Out: out,
		pipes: make([]Pipe[In, Out], 0),
	}
}
