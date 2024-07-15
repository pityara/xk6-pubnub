package pubnub
import (
    "log"
    "sync"

    "go.k6.io/k6/js/modules"
    "github.com/pubnub/go/v7"
)

func init() {
    modules.Register("k6/x/pubnub", new(PubNub))
}

type PubNub struct {
    client *pubnub.PubNub
    channel string
    mu sync.Mutex
}


type Config struct {
    PublishKey   string `json:"publishKey"`
    SubscribeKey string `json:"subscribeKey"`
    Channel      string `json:"channel"`
}

func (p *PubNub) Configure(config Config) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.channel = config.Channel

    pubnubConfig := pubnub.NewConfig("")
    pubnubConfig.PublishKey = config.PublishKey
    pubnubConfig.SubscribeKey = config.SubscribeKey
    p.client = pubnub.NewPubNub(pubnubConfig)

    go p.startPubNubListener()
}

func (p *PubNub) startPubNubListener() {
    listener := p.client.Subscribe().
        Channels([]string{p.channel}).
        Execute()

    for {
        select {
        case msg := <-listener.Message:
            log.Printf("Received message: %s", msg.Message)
            // Save the message or send it to the K6 test
        }
    }
}
