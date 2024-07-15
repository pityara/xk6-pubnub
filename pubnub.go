package pubnub
import (
    "fmt"
    "sync"

    "go.k6.io/k6/js/modules"
    pubnub "github.com/pubnub/go/v7"
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

    pubnubConfig := pubnub.NewConfigWithUserId("123")
    pubnubConfig.PublishKey = config.PublishKey
    pubnubConfig.SubscribeKey = config.SubscribeKey
    p.client = pubnub.NewPubNub(pubnubConfig)

    listener := pubnub.NewListener()
    doneConnect := make(chan bool)
    donePublish := make(chan bool)

    go func() {
        for {
            select {
            case status := <-listener.Status:
                switch status.Category {
                case pubnub.PNDisconnectedCategory:
                    // This event happens when radio / connectivity is lost
                case pubnub.PNConnectedCategory:
                    // Connect event. You can do stuff like publish, and know you'll get it.
                    // Or just use the connected event to confirm you are subscribed for
                    // UI / internal notifications, etc
                    doneConnect <- true
                case pubnub.PNReconnectedCategory:
                    // Happens as part of our regular operation. This event happens when
                    // radio / connectivity is lost, then regained.
                }
            case message := <-listener.Message:
                // Handle new message stored in message.message
                if message.Channel != "" {
                    // Message has been received on channel group stored in
                    // message.Channel
                } else {
                    // Message has been received on channel stored in
                    // message.Subscription
                }
                if msg, ok := message.Message.(map[string]interface{}); ok {
                    fmt.Println(msg["msg"])
                }
                /*
                    log the following items with your favorite logger
                        - message.Message
                        - message.Subscription
                        - message.Timetoken
                */

                donePublish <- true
            case <-listener.Presence:
                // handle presence
            }
        }
    }()

    p.client.AddListener(listener)

    p.client.Subscribe().
       Channels([]string{config.Channel}).
       Execute()

    <-doneConnect
}
