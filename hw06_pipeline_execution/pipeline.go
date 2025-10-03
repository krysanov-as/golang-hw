package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) == 0 {
		return nil
	}

	resultOut := in
	for _, stage := range stages {
		resultOut = executeStageWithDone(stage(resultOut), done)
	}

	return resultOut
}

func executeStageWithDone(in In, done In) Out {
	ch := make(Bi)

	go func() {
		defer func() {
			close(ch)
			for range in {
				continue
			}
		}()

		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case ch <- val:
				}
			}
		}
	}()

	return ch
}
