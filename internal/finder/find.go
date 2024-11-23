package finder

type KekFinder struct {
	fileSearch string // Example: md
}

type KekFinderConf struct {
	FileSearch string
}

func NewFinder(conf KekFinderConf) Finder {
	return &KekFinder{
		fileSearch: conf.FileSearch,
	}
}
