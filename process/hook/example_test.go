package hook

import (
	"github.com/snail007/gmc/process/hook"
)

func Example() {
	//the code block should be last in your main function code

	hook.RegistShutdown(func() {
		//do something before program exit
	})
	hook.WaitShutdown()
}
