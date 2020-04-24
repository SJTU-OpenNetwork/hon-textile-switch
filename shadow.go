package main

import (
	"fmt"
	"os"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/cmd"
)


// The context governs the lifetime of the libp2p node
func main() {
	/*
	ctx := context.Background()
	shadowCtx, _ := context.WithCancel(ctx)
	h, err := host.NewHost(ctx)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	shadowService := service.NewShadowService(ctx, h)
	shadowService.Start()
	fmt.Printf("Hello World, my hosts ID is %s\n", h.ID())
	select{
		case <-shadowCtx.Done():
			fmt.Printf("Routine end.")
			break
	}
	 */
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
