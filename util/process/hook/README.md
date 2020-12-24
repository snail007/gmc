## Demo 

hook package does prevent your program main function to exit when all your worker code in goroutine.

```golang
package hook

import (
	"github.com/snail007/gmc/process/hook"
)

func main() {
	//your business code here
	
	ghook.RegistShutdown(func() {
		//do something before program exit
	})
    //this will waiting for singal, prevent program main function to exit
	ghook.WaitShutdown()
}
```
