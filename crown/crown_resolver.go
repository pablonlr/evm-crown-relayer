package crown

import (
	"log"
	"time"

	crownd "github.com/pablonlr/go-rpc-crownd"
	"github.com/pablonlr/poly-crown-relayer/config"
)

type CrownResolver struct {
	client     *crownd.Client
	unlockPass string
	unlockFor  int
	lastUnlock time.Time
}

func NewCrownResolver(conf config.CrownClientConfig) (*CrownResolver, error) {
	client, err := crownd.NewClient(conf.ClientAdress, conf.ClientPort, conf.Secrets.RPC_USER, conf.Secrets.RPC_PASS, conf.ClientRequestTimeout)
	if err != nil {
		return nil, err
	}

	log.Println("trying to connect to crown deamon...")
	blockCount, err := client.GetBlockCount()
	if err != nil && blockCount == 0 {
		return nil, err
	}
	log.Println("connected to crown deamon, blockcount: ", blockCount)
	return &CrownResolver{
		client:     client,
		unlockPass: conf.Secrets.UnlockPass,
		unlockFor:  conf.UnlockNSeconds,
	}, nil
}

func (crw *CrownResolver) tryToUnclockWallet() error {
	if crw.lastUnlock.Add(time.Duration(crw.unlockFor) * time.Second).After(time.Now()) {
		return nil
	}
	err := crw.client.Unlock(crw.unlockPass, crw.unlockFor)
	if err != nil {
		return err
	}
	crw.lastUnlock = time.Now()
	return nil
}
