package informer

import "rpi_stat_tg_bot/internal/finder"

type RealInformer struct {
	finder    finder.Finder
	root_user string
}

type InformerConf struct {
	Finder finder.Finder
	User   string
}

func NewInformer(conf InformerConf) Informer {
	return &RealInformer{
		finder:    conf.Finder,
		root_user: conf.User,
	}
}
