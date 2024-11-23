package informer

import "rpi_stat_tg_bot/internal/finder"

type KekInformer struct {
	finder    finder.Finder
	root_user string
}

type KekInformerConf struct {
	Finder finder.Finder
	User   string
}

func NewKekInformer(conf KekInformerConf) Informer {
	return &KekInformer{
		finder:    conf.Finder,
		root_user: conf.User,
	}
}
