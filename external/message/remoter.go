package message

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

type ActorProperties struct {
	Address, Name, Kind string
}

type Remoter struct {
	remoter *remote.Remote
	context *actor.RootContext
	mapPID  map[string]*actor.PID
}

func (r *Remoter) Send(ctx appcontext.Context, properties ActorProperties, message proto.Message) error {
	pid := actor.NewPID(properties.Address, properties.Name)
	if actorContext := ctx.GetActorContext(); actorContext != nil {
		actorContext.Send(pid, message)
	} else {
		r.context.Send(pid, message)
	}
	return nil
}

func (r *Remoter) SendMany(ctx appcontext.Context, properties []ActorProperties, message proto.Message) error {
	clients := NewActorPIDSet(properties)
	for _, client := range clients.Values() {
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
type Config struct {
	Remoter Address `yaml:"remoter"`
}

func NewRemoter(system *actor.ActorSystem, cfg Config) *Remoter {
	return &Remoter{
		remoter: remote.NewRemote(system, remote.Configure(cfg.Remoter.Host, cfg.Remoter.Port)),
		context: system.Root,
		mapPID:  make(map[string]*actor.PID),
	}
}

func NewActorPIDSet(properties []ActorProperties) *actor.PIDSet {
	var pids = make([]*actor.PID, 0)
	for _, p := range properties {
		pids = append(pids, actor.NewPID(p.Address, p.Kind))
	}
	return actor.NewPIDSet(pids...)
}
