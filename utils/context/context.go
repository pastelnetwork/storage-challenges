package context

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
	"gorm.io/gorm"
)

type Context interface {
	context.Context
	WithDBTx(tx *gorm.DB) Context
	GetDBTx() *gorm.DB
	WithActorContext(actorCtx actor.Context) Context
	GetActorContext() (actorContext actor.Context)
}

type appContext struct {
	context.Context
}

type ctxKey struct {
	name string
}

var (
	dbTxKey     = ctxKey{name: "dbtx_keyname"} // private key is important for the design
	actorCtxKey = ctxKey{name: "actorcontext_keyname"}
)

func (ctx *appContext) WithDBTx(tx *gorm.DB) Context {
	ctx.Context = context.WithValue(ctx.Context, dbTxKey, tx)
	return ctx
}

func (ctx *appContext) GetDBTx() (tx *gorm.DB) {
	tx, _ = ctx.Context.Value(dbTxKey).(*gorm.DB)
	return
}

func (ctx *appContext) WithActorContext(actorCtx actor.Context) Context {
	ctx.Context = context.WithValue(ctx.Context, actorCtxKey, actorCtx)
	return ctx
}

func (ctx *appContext) GetActorContext() (actorContext actor.Context) {
	actorContext, _ = ctx.Context.Value(actorCtxKey).(actor.Context)
	return
}

func FromContext(ctx context.Context) Context {
	return &appContext{Context: ctx}
}
