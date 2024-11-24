package informer

type Informer interface {
	Basic() (string, string, error)
	FullState() (string, error)
	CMDMount(md string) (string, []string)
	CMDChmod() (string, []string)
	IPFormatter() string
}
