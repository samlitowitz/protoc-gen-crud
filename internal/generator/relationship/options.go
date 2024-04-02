package relationship

type options struct{}

type Option interface {
	apply(*options)
}
