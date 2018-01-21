package models

import (
	"github.com/yubo/falcon/ctrl"
	"golang.org/x/net/context"
)

type ModelsModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *ModelsModule) PreStart(c *ctrl.Ctrl) error {
	return nil
}

func (p *ModelsModule) Start(c *ctrl.Ctrl) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())

	conf := &c.Conf.Ctrl
	Db = ctrl.Orm

	if err := initConfig(conf); err != nil {
		return err
	}
	if err := initAuth(conf); err != nil {
		return err
	}
	if err := initCache(conf); err != nil {
		return err
	}

	SysOp = &Operator{
		O:     Db.Ctrl,
		Token: SYS_F_A_TOKEN | SYS_F_O_TOKEN | SYS_F_A_TOKEN,
	}
	SysOp.User, _ = GetUser(1, SysOp.O)

	putEtcdConfig()

	return nil
}

func (p *ModelsModule) Stop(c *ctrl.Ctrl) error {
	p.cancel()
	return nil
}

func (p *ModelsModule) Reload(c *ctrl.Ctrl) error {
	// TODO
	return nil
}
