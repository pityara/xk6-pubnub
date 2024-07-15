package main

import (
    "os"

    "go.k6.io/xk6"
    "github.com/pityara/xk6-pubnub/pubnub"
)

func main() {
    xk6.Build(xk6.BuildCfg{
        Modules: map[string]xk6.Module{
            "k6/x/pubnub": {FS: os.DirFS("./pubnub")},
        },
    })
}
