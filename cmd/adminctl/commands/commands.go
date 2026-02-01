package commands

type Commands interface {
	Name() string
	Run([]string) error
}
