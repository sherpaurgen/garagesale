package graceful

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	done := make(chan bool, 1)

	go func() {

		sig := <-sigs
		fmt.Println("sigbelow")
		fmt.Println("-----", sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}
