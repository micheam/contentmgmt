package main

import (
    "fmt"
    "github.com/micheam/imgcontent/interfaces/console"
    "os"
)

func main() {
    err := console.NewApp().Run(os.Args)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

