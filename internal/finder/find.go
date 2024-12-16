package finder

type RealFinder struct {
	fileSearch string // Example: md
}

type FinderConf struct {
	FileSearch string
}

func NewFinder(conf FinderConf) Finder {
	return &RealFinder{
		fileSearch: conf.FileSearch,
	}
}
