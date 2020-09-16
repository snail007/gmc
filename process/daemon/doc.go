/*
daemon package does add daemon, logging and forever function to your program.

the below command line arguments you can passed to your program after you using daemon package.

--forver or -forver

the argument will fork a worker process and master process monit worker process , and restart it when it crashed.

--daemon or -daemon

the argument will put your program running in background , only working in linux uinx etc.

--flog or -flog <filename.log>

the argument will logging your program stdout to the log file.

notice:

before you maybe execute your program like this:

./foobar -u root

after using daemon package, execute your program can be like this:

./foobar -u root -forever -daemon -flog foobar.log

*/
package daemon
