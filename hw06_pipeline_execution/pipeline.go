package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		// если стадий нет, мы не можем закрыть in, поэтому создаем прокси канал и возвращаем его
		out := make(Bi)
		go func() {
			defer close(out)
			for {
				select {
				case val, ok := <-in:
					if !ok {
						return
					}
					select {
					case out <- val:
					case <-done:
						return
					}
				case _, _ = <-done:
					return
				}
			}
		}()
		return out
	}

	currentIn := in
	for _, stage := range stages {
		currentIn = stage(currentIn)
	}

	result := make(Bi)

	go func() {
		defer close(result)
		for {
			select {
			case val, ok := <-currentIn:
				if !ok {
					return
				}
				select {
				case result <- val:
				case _, _ = <-done:
					return
				}
			case _, _ = <-done:
				return
			}
		}
	}()

	return result
}
