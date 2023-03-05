# mq
Lightweight generic message queue with timeouts in pure Go
# test
```
go test -race github.com/nullc4t/mq
```
# benchmarks
```
go test -bench=. -run=^# -count=5
```
# usage
```
mq := new(MQ[int])
mq.Push(1)
i := mq.Pop(0)           // without timeout
i := mq.Pop(time.Second) // with timeout
```