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

}