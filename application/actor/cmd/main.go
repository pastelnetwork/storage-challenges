package main

import (
	console "github.com/AsynkronIT/goconsole"
	"github.com/pastelnetwork/storage-challenges/application/actor"
)

func main() {
	pid := actor.StartStorageChallengeHandler()
	defer actor.StopActor(pid)
	defer actor.GracefulStop()
	console.ReadLine()
}
