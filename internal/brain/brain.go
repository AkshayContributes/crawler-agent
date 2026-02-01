package brain

type AI interface {
	Summarize(text string) (string, error)
	Close() error
}
