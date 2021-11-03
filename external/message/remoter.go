package message

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

type Remoter struct {
	remoter *remote.Remote
	context *actor.RootContext
	mapPID  map[string]*actor.PID
}

func (r *Remoter) Send(ctx appcontext.Context, pid *actor.PID, message proto.Message) error {
	if actorContext := ctx.GetActorContext(); actorContext != nil {
		actorContext.Send(pid, message)
	} else {
		r.context.Send(pid, message)
	}
	return nil
}

func (r *Remoter) SendMany(ctx appcontext.Context, pidSet *actor.PIDSet, message proto.Message) error {
	for _, client := range pidSet.Values() {
		if actorContext := ctx.GetActorContext(); actorContext != nil {
			actorContext.Send(client, message)
		} else {
			r.context.Send(client, message)
		}
	}
	return nil
}

func (r *Remoter) RegisterActor(a actor.Actor, name string) (*actor.PID, error) {
	pid, err := r.context.SpawnNamed(actor.PropsFromProducer(func() actor.Actor { return a }), name)
	if err != nil {
		return nil, err
	}
	r.mapPID[name] = pid
	return nil, nil
}

func (r *Remoter) DeregisterActor(name string) {
	if r.mapPID[name] != nil {
		r.context.Stop(r.mapPID[name])
	}
	delete(r.mapPID, name)
}

func (r *Remoter) Start() {
	r.remoter.Start()
}

func (r *Remoter) GracefulStop() {
	for name, pid := range r.mapPID {
		r.context.Stop(pid)
		delete(r.mapPID, name)
	}
	r.remoter.Shutdown(true)
}

type Address struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type Config = Address

func NewRemoter(system *actor.ActorSystem, cfg Config) *Remoter {
	return &Remoter{
		remoter: remote.NewRemote(system, remote.Configure(cfg.Host, cfg.Port)),
		context: system.Root,
		mapPID:  make(map[string]*actor.PID),
	}
}
