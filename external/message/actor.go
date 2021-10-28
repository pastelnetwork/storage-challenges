package message

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

var (
	rootContext *actor.RootContext
)

type Callback func(actor.Context, interface{})

func init() {
	system := actor.NewActorSystem()
	rootContext = system.Root
	config := remote.Configure("127.0.0.1", 0)
	remoter := remote.NewRemote(system, config)
	remoter.Start()
}

func Do(clients *actor.PIDSet, message interface{}, callback Callback) (cancel func()) {
	props := actor.PropsFromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *actor.Started:
			for _, client := range clients.Values() {
				context.Send(client, message)
			}
		case *actor.Stopped:
		case *actor.Stopping:
		case *actor.Restarting:
		case *actor.Restart:
		default:
			callback(context, msg)
		}
	})
	pid := rootContext.Spawn(props)

	return func() { rootContext.Stop(pid) }
}

func NewActorPIDSet(addresses map[string][]string) *actor.PIDSet {
	var pids = make([]*actor.PID, 0)
	for id, addressesByID := range addresses {
		for _, address := range addressesByID {
			pids = append(pids, actor.NewPID(address, id))
		}
	}
	return actor.NewPIDSet(pids...)
}
