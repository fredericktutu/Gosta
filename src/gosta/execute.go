package gosta

import (
	"fmt"
)

//静态构建的目标，动态运行的初始
type StaticProgram struct {  //程序的静态特征
	GFSMs []*GFSM             //所有协程自动机
	Channels []Channel       //所有管道，一次运行中，只有一套管道

	GoroutineCount int  // 运行时可能的总协程数
	AllStateCount int   // 静态程序的总状态数
}

type Runtime struct {
	Goroutines []RGoroutine
	Live map[int]int

	Channels []RChannel
	Logs ActionStack     //日志栈

	LogCounter int    //当前日志编号

	Result int // 抽象执行过程中，产生协程死锁的数量
	Paths int  // 抽象执行过程中，实际分析的路径数
}


type RChannel struct {
	Cap int
	Now int        //容量,发现它不需要用栈，用正常加减就可以了
}

type RGoroutine struct {
	Machine GID
	Bindings []ChanID
	Current SIDStack      //状态栈
}

type ActionType int
const (
	ACTION_MAIN  ActionType = iota   //初始化Main函数
	ACTION_SEND
	ACTION_RECV
	ACTION_EPSILON
	ACTION_MATCH
	ACTION_ILLEGAL

	
	ACTION_SINGLE_SEND
	ACTION_SINGLE_RECV
)

type Action struct {
	StepID int       // 一步可能有多个动作
	Score int        // 打分值
	
	Onwhich int    // 在哪个Goroutine上操作
	ToState SID      // 到哪个State

	Onwhich2 int  //ACTION_MATCH情况下，RECV的
	ToState2 SID  //ACTION_MATCH情况下,RECV的

	Type ActionType  // 哪种动作
	ChanID int        // SEND/RECV时，用于改变Chan
	Goroutines []*Goroutine  // GO时，用于绑定Goroutine初值
}

func genAction(e *Edge, step int, onwhich int, sp *StaticProgram, rt *Runtime) Action {
	var a Action
	a.StepID = step      // 当前的stepID
	a.Onwhich = onwhich  // 在哪个协程上进行操作
	a.Onwhich2 = onwhich  // fix bug:用于保证容量为0的channel

	realChanID := e.ChanID
	if e.ChanID < 0 {
		realChanID = rt.Goroutines[onwhich].Bindings[-e.ChanID - 1]  // Binding列表中获取
	}

	// 需要判断是否可以进行操作
	switch e.Op {
	case OP_SEND:
		if rt.Channels[realChanID].Cap == 0 {
			a.Type = ACTION_SINGLE_SEND
			a.ChanID = int(realChanID)

			a.ToState = e.To
			toState := getState(sp, e.To)
			a.Goroutines = toState.Gos  
			return a
		}

		if rt.Channels[realChanID].Now >= rt.Channels[realChanID].Cap {
			a.Type = ACTION_ILLEGAL
			return a
		}
	
		a.Type = ACTION_SEND
		a.ChanID = int(realChanID)


	case OP_RECV:
		if rt.Channels[realChanID].Cap == 0 {
			a.Type = ACTION_SINGLE_RECV
			a.ChanID = int(realChanID)

			a.ToState2 = e.To
			toState2 := getState(sp, e.To)
			a.Goroutines = toState2.Gos  
			return a
		}

		if rt.Channels[realChanID].Now <= 0 {
			a.Type = ACTION_ILLEGAL
			return a
		}

		a.Type = ACTION_RECV
		a.ChanID = int(realChanID)


	case OP_EPSILON:
		a.Type = ACTION_EPSILON
	}

	a.ToState = e.To
	toState := getState(sp, e.To)
	a.Goroutines = toState.Gos  // CHECK: 拷贝方式

	return a
}

func doAction(rt *Runtime, a Action) {  
	/*
	上一层： 
	1. 记得之前要设置StepID
	2. 每次状态转移，现在只用打一个日志
	3. a.ChanID 是需要先计算好的，如果是非Main Goroutine，在外面看到它的Bindings是负数
	*/

	fmt.Println("[Do Action]", "ActionType", a.Type, "| On Goroutine", a.Onwhich, "| ToState", a.ToState, "|ChanID", a.ChanID, "|Chan now", rt.Channels[a.ChanID].Now)
	switch a.Type {

	case ACTION_SEND:  
		// 1. 转移RGoroutine的状态
		rt.Goroutines[a.Onwhich].Current.Push(a.ToState) 

		// 2. 增加RChannel的当前量
		rt.Channels[a.ChanID].Now += 1  // 到这里一定成功

	case ACTION_RECV:
		// 1. 转移RGoroutine的状态
		rt.Goroutines[a.Onwhich].Current.Push(a.ToState) 

		// 2. 减少RChannel的当前量
		rt.Channels[a.ChanID].Now -= 1  // 到这里一定成功

	case ACTION_EPSILON:  
		/* 只变协程状态*/

		// 1. 转移RGoroutine的状态
		rt.Goroutines[a.Onwhich].Current.Push(a.ToState)

	case ACTION_MATCH:
		// 1. 对于操作SEND和RECV的这对协程，改变状态
		rt.Goroutines[a.Onwhich].Current.Push(a.ToState)
		rt.Goroutines[a.Onwhich2].Current.Push(a.ToState2)

		// 2. 无需改变空管道的Now

	case ACTION_MAIN:
		break
	}


	// -2. 启动若干Goroutine，设置绑定列表
	for _, gp := range a.Goroutines {
		// gp.Num 是Goroutine在运行时的下标
		fmt.Println("    <Go func>", gp.Num)
		if rt.Goroutines[gp.Num].Machine == NO_MACHINE {  
			// 以前没有设置过Goroutine，则初始化一下
			fmt.Println("    <Binding>", gp.Bindings)
			rt.Goroutines[gp.Num].Bindings = gp.Bindings  // CHECK:这样是值拷贝还是指针拷贝？ 不过无论是哪种，大概都没问题，因为不会修改
			rt.Goroutines[gp.Num].Machine = gp.Machine
		}
		rt.Goroutines[gp.Num].Current.Push(SID{GID(gp.Machine), 0})  // 对应的GFSM的初始状态(0)
		rt.Live[gp.Num] = 1 //这个RGoroutine由死亡 -> 活动
	}


	// -1. 追加日志
	rt.Logs.Push(a)
	return
}

func undoAction(rt *Runtime) {
	la := rt.Logs.Top()
	fmt.Println("[Undo Action]", "ActionType", la.Type, "| On Goroutine", la.Onwhich, "| ToState", la.ToState, "|ChanID", la.ChanID, "|Chan now", rt.Channels[la.ChanID].Now)
	switch la.Type {

	case ACTION_SEND:  
		// 1. 恢复RGoroutine的状态
		rt.Goroutines[la.Onwhich].Current.Pop()

		// 2. 恢复原本增加的RChannel当前量
		rt.Channels[la.ChanID].Now -= 1


	case ACTION_RECV:
		// 1. 恢复RGoroutine的状态
		rt.Goroutines[la.Onwhich].Current.Pop()

		// 2. 恢复原本减少的RChannel的当前量
		rt.Channels[la.ChanID].Now += 1

	case ACTION_EPSILON:  
		/* 只变协程状态*/

		// 1. 转移RGoroutine的状态
		rt.Goroutines[la.Onwhich].Current.Pop()
	
	case ACTION_MATCH:
		// 1. 对于操作SEND和RECV的这对协程，恢复状态
		rt.Goroutines[la.Onwhich].Current.Pop()
		rt.Goroutines[la.Onwhich2].Current.Pop()

		// 2. 无需改变空管道的Now

	}

	// -2. 将启动过的Goroutine撤销
	for _, gp := range la.Goroutines {
		// gp.Num 是Goroutine在运行时的下标

		rt.Goroutines[gp.Num].Current.Pop()  // 把这个Goroutine的最后一个状态Pop掉，现在它的Current为空
		delete(rt.Live, gp.Num)  // 这个RGoroutine由活动->死亡
	}	
	
	// -1. pop掉这个日志
	rt.Logs.Pop()  
	return 
}



func isDead(rg *RGoroutine, sp *StaticProgram) bool{
	/*
	一个RGoroutine已死
	*/
	if rg.Current.Top().ID == len(sp.GFSMs[rg.Machine].States) - 1 {  // 已经在终止状态了
		return true
	} else {
		return false
	}
}

func getState(sp *StaticProgram, sid SID) *State{
	/*
	返回一个sp中的状态
	*/
	return sp.GFSMs[sid.GID].States[sid.ID]
}


func MakeRuntime(sp *StaticProgram, totalScore int) *Runtime{
	/*
	根据StaticProgram创建一个新的运行时
	并初始化
	*/

	rt := new(Runtime)

	rt.Result = 0
	rt.Paths = 0

	// 1. 初始化Channels
	for _, c := range(sp.Channels) {
		var rc RChannel
		rc.Now = 0
		rc.Cap = c.Cap

		rt.Channels = append(rt.Channels, rc)
	}

	//TODO：后面把这个函数里相关的逻辑移到DoAction
	// 2. 初始化Goroutines
	// main的bindings是空
	var mainRoutine RGoroutine
	mainRoutine.Machine = GID(0)
	mainRoutine.Current.Push(SID{GID(0), 0})
	rt.Goroutines = append(rt.Goroutines, mainRoutine)

	// 剩下的Goroutine先占个位
	var i int = 0
	for {  
		if i >= sp.GoroutineCount { // 注意数量
			break
		}
		i += 1

		var rg RGoroutine
		rg.Machine = NO_MACHINE
		rt.Goroutines = append(rt.Goroutines, rg)
	}

	// 3. 初始化Live
	rt.Live = make(map[int]int)
	rt.Live[0] = 1

	// 4. 进行第一个操作 
	var firstAct Action
	firstAct.StepID = 0
	firstAct.Score = totalScore         // 仅记录
	firstAct.ToState = SID {GID(0), 0}  // 第0个自动机的0状态即Main的开始
	firstAct.Type = ACTION_MAIN
	firstAct.Goroutines = getState(sp, SID{GID(0), 0}).Gos  // 一开始就启动的那些协程

	doAction(rt, firstAct)

	rt.LogCounter = 1                   //目前有1个
	
	return rt
}


func rank1(availables *[]Action, limit int) {

	var num = len(*availables)
	portion := (limit / num)
	if portion < 1 {
		portion = 1
	}

	for i, _ := range *availables {
		if limit >= portion {
			(*availables)[i].Score = portion
			limit -= portion
		} else {
			(*availables)[i].Score = limit
			limit = 0
		}
	}
}


func Execute(rt *Runtime, sp *StaticProgram, limit int) {
	/*
	递归地执行运行时
	*/

	// 检查退出条件
	mainRoutine := rt.Goroutines[0]
	if isDead(&mainRoutine, sp) {  // 已经在终止状态了
		fmt.Println("[Execute SUCC]")
		rt.Paths += 1
		return
	}


	var availables []Action
	singleSends := make(map[int]([]Action))
	singleRecvs := make(map[int]([]Action))

	for gindex, _ := range rt.Live {
		rgoroutine := rt.Goroutines[gindex]

		if isDead(&rgoroutine, sp) {  // 已结束的，不考虑动作
			continue
		}

		curState := getState(sp, rgoroutine.Current.Top())
		
		fmt.Println("[Analyze Edges]", "RGoroutine", gindex, "|At State", curState)



		for i, e := range curState.Out{
			fmt.Println("    <Edge>", i, "| To", e.To, "| Op", e.Op, "|On Chan", e.ChanID)

			oneAct := genAction(e, rt.LogCounter, gindex, sp, rt)
			if oneAct.Type == ACTION_ILLEGAL {
				continue
			} else if oneAct.Type == ACTION_SINGLE_SEND {
				sends, ok := singleSends[oneAct.ChanID]
				if ok {
					singleSends[oneAct.ChanID] = append(sends, oneAct)
				}else {
					singleSends[oneAct.ChanID] = make([]Action, 1)
					singleSends[oneAct.ChanID][0] = oneAct
				}

				continue
			} else if oneAct.Type == ACTION_SINGLE_RECV {
				recvs, ok := singleRecvs[oneAct.ChanID]
				if ok {
					singleRecvs[oneAct.ChanID] = append(recvs, oneAct)
				}else {
					singleRecvs[oneAct.ChanID] = make([]Action, 1)
					singleRecvs[oneAct.ChanID][0] = oneAct
				}
				continue
			}

			availables = append(availables, oneAct)
		}
	}
	fmt.Println("[singleSends]", singleSends, "[singleRecvs]", singleRecvs)

	for keyChanID, sends := range singleSends {  // 容量为0的channels
		recvs, ok := singleRecvs[keyChanID]
		if ok {
			for _, sendAction := range sends {
				for _, recvAction := range recvs {
					var matchAction Action
					matchAction.Onwhich = sendAction.Onwhich
					matchAction.ToState = sendAction.ToState
					matchAction.Onwhich2 = recvAction.Onwhich2
					matchAction.ToState2 = recvAction.ToState2
					matchAction.Type = ACTION_MATCH
					matchAction.ChanID = keyChanID
					for _, goro := range sendAction.Goroutines {
						matchAction.Goroutines = append(matchAction.Goroutines, goro)
					}
					for _, goro := range recvAction.Goroutines {
						matchAction.Goroutines = append(matchAction.Goroutines, goro)
					}
					availables = append(availables, matchAction)
				}
			}
		} else {
			continue
		}
	}


	// 打印各个Action
	fmt.Println("[available actions]:")
	for i, action := range availables {
		fmt.Println("<action>", i, action)
	}

	// available空，协程死锁
	if len(availables) == 0 {
		fmt.Println("[ERROR] goroutine deadlock!")
		rt.Paths += 1
		rt.Result += 1
		return
	}



	//打分
	rank1(&availables, limit)
	//排序

	//依次
	for _, action := range availables {
		if action.Score <= 0 {
			continue
		}

		//运行
		rt.LogCounter += 1
		doAction(rt, action)

		//递归
		Execute(rt, sp, action.Score)

		//恢复
		undoAction(rt)
		rt.LogCounter -= 1
	}

}
