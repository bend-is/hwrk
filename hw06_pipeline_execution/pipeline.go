package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return wrapWithCancel(in, done)
	}

	for _, stage := range stages {
		in = stage(wrapWithCancel(in, done))
	}

	return in
}

func wrapWithCancel(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				out <- v
			}
		}
	}()

	return out
}
