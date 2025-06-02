package main

import sshmanager "github.com/zine0/gsshm/internal/SSHManager"


func main(){
	mgr := sshmanager.NewSSHManager()
	mgr.Connect("127.0.0.1","22","zine","zine","")
	mgr.StartTerminal()	
}