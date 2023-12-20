package metrics

type Reporter interface {
	Schedule(forceRunNow bool) error
}
