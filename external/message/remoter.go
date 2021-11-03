package message

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

type Config struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	clientSecureCres credentials.TransportCredentials
	serverSecureCres credentials.TransportCredentials
}

func (c *Config) WithClientSecureCres(s credentials.TransportCredentials) *Config {
	c.clientSecureCres = s
	return c
}

func (c *Config) WithServerSecureCres(s credentials.TransportCredentials) *Config {
	c.serverSecureCres = s
	return c
}

func NewRemoter(system *actor.ActorSystem, cfg Config) *Remoter {
	if cfg.Host == "" {
		cfg.Host = "0.0.0.0"
	}
	if cfg.Port == 0 {
		cfg.Port = 9000
	}
	remoterConfig := remote.Configure(cfg.Host, cfg.Port)
	clientCres := []grpc.DialOption{grpc.WithBlock()}
	serverCres := []grpc.ServerOption{}
	if cfg.clientSecureCres != nil {
		clientCres = append(clientCres, grpc.WithTransportCredentials(cfg.clientSecureCres))
	}
	if cfg.serverSecureCres != nil {
		serverCres = append(serverCres, grpc.Creds(cfg.serverSecureCres))
	}
	remoterConfig.WithDialOptions(clientCres...).WithServerOptions(serverCres...)
	return &Remoter{
		remoter: remote.NewRemote(system, remoterConfig),
		context: system.Root,
		mapPID:  make(map[string]*actor.PID),
	}
}
