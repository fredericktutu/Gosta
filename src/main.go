package main
import (
	"gosta"
	"os"
	"fmt"
)


func main() {
	//Args[1]为文件名
	sp := gosta.BuildFSMs(os.Args[1])
	fmt.Println("BuildHelper after build:\n", sp)
	if len(os.Args) > 2  {
		if os.Args[2] == "--noexec" {
			return
		} else {
			
		}
	} else {
		fmt.Println("add --noexec or --exec")
	}

	var N int = 1000
	rt := gosta.MakeRuntime(sp, N)
	gosta.Execute(rt, sp, N)
	fmt.Println("######################")

	fmt.Println("[", os.Args[1] ,  "] ' s result:", rt.Result, "BUGs")


}