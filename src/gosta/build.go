package gosta

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

	"strings"
	"strconv"
)

type BuildHelper struct {
	Pack *ssa.Package

	ChanEnv map[string]int  
	// REALCHAN{"[func].[reg]" -> chanID(>0)} 
	// PARAM{"[func].[reg]" -> chanID(<0)}

	Channels []Channel

	FuncEnv map[string]int  //函数名-下标映射表 walk前添加
	Funcs []*ssa.Function  //函数列表，walk前添加
	GFSMs []*GFSM  //自动机列表，walk结束后添加

	GoroutineCounter int  // 运行时Goroutine序号

	//构建过程中暂时使用
	CurFuncname string
	StateCounter int
	CurMachine int
}


func makeBuildHelper() *BuildHelper {
	bh := new(BuildHelper)
	bh.ChanEnv = make(map[string]int)
	bh.FuncEnv = make(map[string]int)
	return bh
}

func lookupChan(reg string, bhelper *BuildHelper) ChanID {
	var idInt = bhelper.ChanEnv[reg] 

	return ChanID(idInt)
}

func walkFunc(f *ssa.Function, bhelper *BuildHelper) {
	var sfsmLst []*SmallFSM
	var bbMap = make(map[*ssa.BasicBlock]int)



	if bhelper.CurFuncname != "main" {
		/*
		对于非main的函数，
		需要根据参数列表来设置chanEnv
		这里"[func].[reg]"映射到的下标都是负数，代表参数
		*/

		fmt.Println("[Parameters]")

		countParam := -1
		for _, param := range f.Params {
			fmt.Println("    <Param>", countParam, "|Reg:", param.Name())   			
			// 这个Name()是Value类型的方法
			// 代表参数存在的虚拟寄存器
			// 并且和形参名字一样

			bhelper.ChanEnv[fmt.Sprintf("%s.%s", bhelper.CurFuncname, param.Name())] = countParam 
			// countParam是个负数
			// 到时候映射到实际的Bindings列表应该这么算：
			// Bindings[(-countParam) - 1]

			countParam -= 1
		}
		fmt.Println("---");  // 分析Parameters结束了
	}
	
	//从bb map到它们的下标，方便处理bb之间的边
	for i, bb := range f.Blocks {
		bbMap[bb] = i
	}

	bhelper.StateCounter = 0
	for i, bb := range f.Blocks {
		fmt.Println("[bb]:", i)
		sfsm := walkBB(bb, bhelper)
		sfsmLst = append(sfsmLst, sfsm)

		fmt.Println("[bb]:", i, "'s sfsm:", sfsm)
		fmt.Println("---"); fmt.Println()

		// 处理bb之间的边，EPSILON边
		for _, succ := range bb.Succs {
			sfsm.Outs = append(sfsm.Outs, bbMap[succ])
		}
	}

	// 根据BB之间的连接，将SmallFSM里面各个State连起来
	for _, sfsm := range sfsmLst {
		sourceState := sfsm.Last
		for _, num := range sfsm.Outs {
			targetState := sfsmLst[num].First  // 要连接的目标state
			LinkState(sourceState, targetState, 0, OP_EPSILON)  // 连接两个State

		}
	}

	fmt.Println("-------------------------")  // 一个函数的Analysis结束了

}

func walkBB(bb *ssa.BasicBlock, bhelper *BuildHelper) *SmallFSM {

	var sfsm *SmallFSM = MakeSmallFSM(bhelper.GFSMs[bhelper.CurMachine], bhelper.CurMachine, bhelper.StateCounter)
	bhelper.StateCounter += 1

	for i, instr := range bb.Instrs {
		switch instr.(type) {
		case *ssa.Send:
			/*
			SEND
			所需信息： 1. SEND到哪个管道(寄存器 -> 管道实体)
			*/
			instrSend, _ := instr.(*ssa.Send)
			chanID := lookupChan(fmt.Sprintf("%s.%s", bhelper.CurFuncname, instrSend.Chan.Name()), bhelper)
			fmt.Println("    <instr>", i, "SEND | To Reg:", instrSend.Chan.Name(), "| Actually is Chan", chanID)

			AddState(sfsm, bhelper.GFSMs[bhelper.CurMachine], OP_SEND, chanID, bhelper.CurMachine, bhelper.StateCounter)  //增加状态

		case *ssa.UnOp:
			/*
			RECV
			所需信息：1. 从哪个channel进行RECV(寄存器->channel实体)
			*/
			instrUnOp, _ := instr.(*ssa.UnOp)

			if instrUnOp.Op.String() == "<-" {
				chanID := lookupChan(fmt.Sprintf("%s.%s", bhelper.CurFuncname, instrUnOp.X.Name()), bhelper)
				fmt.Println("    <instr>", i, "RECV | From Reg:", instrUnOp.X.Name(), "| Actually is Chan", chanID)  // Name()返回存channel的寄存器，接下来需要确认是哪个channel

				AddState(sfsm, bhelper.GFSMs[bhelper.CurMachine], OP_RECV, chanID, bhelper.CurMachine, bhelper.StateCounter)  // 增加状态

			}

		case *ssa.MakeChan:
			/*
			MAKECHAN
			所需信息： 1. 创建的管道的容量(int)  2. 创建的管道存入的寄存器(寄存器 + channel实体)
			*/
			instrMakeChan, _ := instr.(*ssa.MakeChan)
			fmt.Println("    <instr>", i, "MAKECHAN |", instrMakeChan)
			fmt.Println("        (size)", instrMakeChan.Size.Name())
			fmt.Println("        (reg)", instrMakeChan.Name())

			bhelper.ChanEnv[fmt.Sprintf("%s.%s", bhelper.CurFuncname, instrMakeChan.Name())] = len(bhelper.Channels)  // func.tN -> index_of_chan
			var c Channel
			c.ID = ChanID(len(bhelper.Channels))  // chan的下标

			sizeStr := instrMakeChan.Size.Name()  // 把代表容量的字符串转成int
			sizeInt64, _ := strconv.ParseInt(sizeStr[0: strings.Index(sizeStr, ":")], 10, 64)
			c.Cap = int(sizeInt64)

			bhelper.Channels = append(bhelper.Channels, c)
		
		case *ssa.Go:
			/*
			Go
			所需信息： 1. 启动哪个函数(名字->函数实体)  2. 传入的参数有哪些channel(寄存器->channel实体)
			*/
			instrGo, _ := instr.(*ssa.Go)
			fmt.Println("    <instr>", i, "GO | Func name:", instrGo.Call.Value.Name(), "| Args:")

			//追加一个Goroutine,并填充一下
			var g *Goroutine = AddGoroutine(sfsm)
			g.Num = bhelper.GoroutineCounter
			bhelper.GoroutineCounter += 1

			if funcIndex, ok := bhelper.FuncEnv[instrGo.Call.Value.Name()]; ok {  // 找的到
				g.Machine = GID(funcIndex)

			} else {  // 没找到
				newFunc := bhelper.Pack.Func(instrGo.Call.Value.Name())

				g.Machine = GID(len(bhelper.Funcs))
				bhelper.FuncEnv[instrGo.Call.Value.Name()] = len(bhelper.Funcs)

				bhelper.Funcs = append(bhelper.Funcs, newFunc)
			}

			for j, param := range instrGo.Call.Args {
				if strings.HasPrefix(param.Type().String(), "chan") {  // 类型为chan的
					chanID := lookupChan(fmt.Sprintf("%s.%s", bhelper.CurFuncname, param.Name()), bhelper)
					fmt.Println("        (param)", j, "| Name:",  param.Name(), "| Type:", param.Type().String(), "| Actually is Chan", chanID)
					
					g.Bindings = append(g.Bindings, chanID)

				}
			}
		
		default:
			/*
			其他指令，忽略掉
			*/
			fmt.Println("    <instr>", i, "other")
		}
	} 
	return sfsm
}



func BuildFSMs(fileName string) *StaticProgram{
	// Parse the source files.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Print(err) // parse error
		return nil
	}
	files := []*ast.File{f}

	// Create the type-checker's package.
	pkg := types.NewPackage("main", "")

	// Type-check the package, load dependencies.
	// Create and build the SSA program.
	obj, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()}, fset, pkg, files, ssa.SanityCheckFunctions)
	if err != nil {
		fmt.Print(err) // type error in some package
		return nil
	}

	/***********************************以上是抄来的代码，用来构建SSA****************************************/

	/***********************************参考以下代码，根据SSA构建自动机****************************************/
	var bhelper *BuildHelper = makeBuildHelper()
	bhelper.Pack = obj
	bhelper.GoroutineCounter = 0
	// 
	for name, member := range obj.Members {
		// 处理main
		if member.Token().String() == "func" && name == "main"{
			obj.Func(name).WriteTo(os.Stdout)

			bhelper.FuncEnv[name] = 0
			bhelper.Funcs = append(bhelper.Funcs, obj.Func(name))
			bhelper.CurFuncname = name

			bhelper.GFSMs = append(bhelper.GFSMs, new(GFSM))
			bhelper.GFSMs[0].ID = GID(0)
			bhelper.CurMachine = 0

			fmt.Println("<->")

			walkFunc(obj.Func(name), bhelper)
			break
		}
	}

	var curFuncIndex int = 1

	for {
		if curFuncIndex >= len(bhelper.Funcs) {
			break
		}
		curFunc := bhelper.Funcs[curFuncIndex]
		curFunc.WriteTo(os.Stdout)

		bhelper.CurFuncname = curFunc.Name()

		bhelper.GFSMs = append(bhelper.GFSMs, new(GFSM))
		bhelper.GFSMs[curFuncIndex].ID = GID(curFuncIndex)
		bhelper.CurMachine = curFuncIndex

		fmt.Println("<->")

		walkFunc(curFunc, bhelper)

		curFuncIndex += 1
	}

	var sp *StaticProgram = new(StaticProgram)
	sp.Channels = bhelper.Channels
	sp.GFSMs = bhelper.GFSMs
	sp.GoroutineCount = bhelper.GoroutineCounter

	return sp

}