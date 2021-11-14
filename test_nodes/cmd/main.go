package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/config"
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/pastelnetwork/storage-challenges/external/storage"
	testnodes "github.com/pastelnetwork/storage-challenges/test_nodes"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
	"github.com/pastelnetwork/storage-challenges/utils/xordistance"
)

func main() {
	cfg := config.Config{}
	if err := cfg.Load(); err != nil {
		panic(fmt.Sprintf("could not load config data: %v", err))
	}
	if cfg.Remoter == nil {
		cfg.Remoter = &message.Config{}
	}

	store, err := storage.NewStore(*cfg.Database)
	if err != nil {
		log.Panicf("could not create new database connection: %v", err)
	}

	db := store.GetDB()
	pastelClient := testnodes.NewMockPastelClient(db)
	remoter := message.NewRemoter(actor.NewActorSystem(), *cfg.Remoter.
		WithClientSecureCreds(credentials.NewClientCreds(pastelClient, &alts.SecInfo{PastelID: "mock pastel id", PassPhrase: "mock passphrase", Algorithm: "mock algorithm"})).
		WithServerSecureCreds(credentials.NewServerCreds(pastelClient, &alts.SecInfo{PastelID: "mock pastel id", PassPhrase: "mock passphrase", Algorithm: "mock algorithm"})))
	remoter.Start()
	defer remoter.GracefulStop()
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	var blockCount int32 = 1
	for ; blockCount <= 100; blockCount++ {
		<-ticker.C
		if err = testnodes.AddPastelBlock(blockCount, db); err != nil {
			log.Panicf("could not add new pastel block to database: %v", err)
			return
		}
		if err = testnodes.AddNIncrementalMasternodesAndKIncrementalSymbolFiles(2, 60, db); err != nil {
			log.Panicf("could not add new incremental masternodes and symbol files to database: %v", err)
			return
		}
		if blockCount%5 == 0 {
			if err = testnodes.RemoveMasternodesAndSymbolFiles(db); err != nil {
				log.Panicf("could not delete existing masternodes and symbol files from database: %v", err)
				return
			}
		}
		sliceOfMasternodes := testnodes.GetMasternodes(db)
		sliceOfMasternodeIDs := make([]string, len(sliceOfMasternodes))
		mapMasternodes := make(map[string]testnodes.Masternode)
		for idx, masternode := range sliceOfMasternodes {
			mapMasternodes[masternode.NodeID] = masternode
			sliceOfMasternodeIDs[idx] = masternode.NodeID
		}
		numberOfMasternodesToIssueChallengesPerBlock := int(math.Ceil(float64(len(sliceOfMasternodes)) / 3))
		challengesPerMasternodePerBlock := int32(math.Ceil(float64(len(sliceOfMasternodes)) / 3))
		blockHash := testnodes.GetPastelBlockHash(blockCount)
		sliceOfChallengingMasternodeIdsForBlock := xordistance.GetNClosestXORDistanceStringToAGivenComparisonString(numberOfMasternodesToIssueChallengesPerBlock, blockHash, sliceOfMasternodeIDs)
		// test with existing SNs and symbol file hashes
		for _, challengingMasternodeID := range sliceOfChallengingMasternodeIdsForBlock {
			generatingChallengeMasternodeID := xordistance.GetNClosestXORDistanceStringToAGivenComparisonString(1, blockHash+fmt.Sprint(challengesPerMasternodePerBlock)+challengingMasternodeID, sliceOfMasternodeIDs)[0]
			generatingChallengeMasternodeAddress := mapMasternodes[generatingChallengeMasternodeID].MasternodeIPAddress
			pid := actor.NewPID(generatingChallengeMasternodeAddress, "storage-challenge")
			remoter.Send(appcontext.FromContext(context.Background()), pid, &dto.GenerateStorageChallengeRequest{CurrentBlockHash: blockHash, ChallengingMasternodeId: challengingMasternodeID, ChallengesPerMasternodePerBlock: challengesPerMasternodePerBlock})
		}
	}
	console.ReadLine()
}
