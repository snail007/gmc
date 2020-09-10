package profiling

func ExampleStartArg() {
	//StartArg will search command line arguments "profiling" , using it's value as store data folder.
	StartArg("profiling")
	defer Stop()
}

func ExampleStart() {
	//Start will using the argument as store data folder
	Start("debug")
	defer Stop()
}
