package main
import (
	"gosta"
	"os"
	"fmt"
	"strconv"
)


func main() {
	
	if len(os.Args) != 4 {
		fmt.Println("Input format:\n[1] File Name(string) [2] Weight(int) [3] '--noexec' or '--exec'")
		return
	}
	//Args[1]为文件名
	sp := gosta.BuildFSMs(os.Args[1])
	fmt.Println("Static Program after build:\n", sp)
	
	if os.Args[3] == "--noexec" {
		return
	} else {
		
	}
	

	fmt.Println("\n*****Begin Execution*****")
	intval64, _ :=  strconv.ParseInt(os.Args[2], 10, 0)
	N := int(intval64)
	
	rt := gosta.MakeRuntime(sp, N)
	gosta.Execute(rt, sp, N)
	fmt.Println("*****End Execution*****")

	fmt.Println("\n---REPORT---")

	fmt.Println("Detected", rt.Result, "BUGs in [", os.Args[1], "].")
	fmt.Println("Generate", sp.AllStateCount, "STATEs in total.")
	fmt.Println("Execute", rt.Paths, "PATHs in total")


	fmt.Println("-----------")
	os.Exit(rt.Result)
	


}