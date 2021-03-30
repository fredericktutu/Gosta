package gosta

//静态构建的目标，动态运行的初始
type StaticProgram struct {  //程序的静态特征
	GFSMs []*GFSM             //所有协程自动机
	Channels []Channel       //所有管道，一次运行中，只有一套管道

	GoroutineCount int
}

type Runtime struct {
	Goroutines []RGoroutine
	Channels []RChannel
	Logs ActionStack     //日志栈

	LogCounter int    //当前日志编号
}


type RChannel struct {
	Cap int
	Now IntStack        //容量栈 TODO: 发现它好像是不需要的，用正常加减就可以了
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
)

type Action struct {
	StepID int       // 一步可能有多个动作
	Score int        // 打分值
	
	Onwhich int    // 在哪个Goroutine上操作
	ToState SID      // 到哪个State

	Type ActionType  // 哪种动作
	ChanID int        // SEND/RECV时，用于改变Chan
	Goroutines []*Goroutine  // GO时，用于绑定Goroutine初值
}

func doAction(rt *Runtime, a Action) {  
	/*
	上一层： 
	1. 记得之前要设置StepID
	2. 每次状态转移，现在只用打一个日志
	3. a.ChanID 是需要先计算好的，如果是非Main Goroutine，在外面看到它的Bindings是负数
	*/
	switch a.Type {

	case ACTION_SEND:  
		// 1. 转移RGoroutine的状态
		rt.Goroutines[a.Onwhich].Current.Push(a.ToState) 

		// 2. 增加RChannel的当前量
		curCap := rt.Channels[a.ChanID].Now.Top()
		rt.Channels[a.ChanID].Now.Push(curCap  + 1)  // 到这里了就一定能加成功


	case ACTION_RECV:
		// 1. 转移RGoroutine的状态
		rt.Goroutines[a.Onwhich].Current.Push(a.ToState) 

		// 2. 减少RChannel的当前量
		curCap := rt.Channels[a.ChanID].Now.Top()
		rt.Channels[a.ChanID].Now.Push(curCap - 1)  // 到这里了就一定能减成功

	case ACTION_EPSILON:  
		/* 只变协程状态*/

		// 1. 转移RGoroutine的状态
		rt.Goroutines[a.Onwhich].Current.Push(a.ToState) 
	}


	// -2. 启动若干Goroutine，设置绑定列表
	for _, gp := range a.Goroutines {
		// gp.Num 是Goroutine在运行时的下标

		if rt.Goroutines[gp.Num].Machine == NO_MACHINE {  
			// 以前没有设置过Goroutine，则初始化一下
			rt.Goroutines[gp.Num].Bindings = gp.Bindings  // 这样是值拷贝还是指针拷贝？ 不过无论是哪种，大概都没问题，因为不会修改
			rt.Goroutines[gp.Num].Machine = gp.Machine
		}
		rt.Goroutines[gp.Num].Current.Push(SID{GID(gp.Machine), 0})  // 对应的GFSM的初始状态(0)
	}


	// -1. 追加日志
	rt.Logs.Push(a)
	return
}

func undoAction(rt *Runtime) {
	la := rt.Logs.Top()
	
	switch la.Type {

	case ACTION_SEND:  
		// 1. 恢复RGoroutine的状态
		rt.Goroutines[la.Onwhich].Current.Pop()

		// 2. 恢复原本增加的RChannel当前量
		rt.Channels[la.ChanID].Now.Pop()


	case ACTION_RECV:
		// 1. 恢复RGoroutine的状态
		rt.Goroutines[la.Onwhich].Current.Pop()

		// 2. 恢复原本减少的RChannel的当前量
		rt.Channels[la.ChanID].Now.Pop() 

	case ACTION_EPSILON:  
		/* 只变协程状态*/

		// 1. 转移RGoroutine的状态
		rt.Goroutines[la.Onwhich].Current.Pop()

	}

	// -2. 将启动过的Goroutine撤销
	for _, gp := range la.Goroutines {
		// gp.Num 是Goroutine在运行时的下标

		rt.Goroutines[gp.Num].Current.Pop()  // 把这个Goroutine的最后一个状态Pop掉，现在它的Current为空
	}	
	
	// -1. pop掉这个日志
	rt.Logs.Pop()  
	return 
}

func GenAction(e *Edge) *Action {
	// TODO
	return nil
}



func MakeRuntime(sp *StaticProgram, totalScore int) *Runtime{
	/*
	根据StaticProgram创建一个新的运行时
	并初始化
	*/

	rt := new(Runtime)

	// 1. 初始化Channels
	for _, c := range(sp.Channels) {
		var rc RChannel
		rc.Now.Push(0)
		rc.Cap = c.Cap

		rt.Channels = append(rt.Channels, rc)
	}

	// 2. 初始化Goroutines
	// 它的bindings是空
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

	// 3. 初始化日志
	var firstAct Action
	firstAct.StepID = 0
	firstAct.Score = totalScore
	firstAct.ToState = SID {GID(0), 0}  // 第0个自动机的0状态即Main的开始
	firstAct.Type = ACTION_MAIN

	rt.Logs.Push(firstAct)
	rt.LogCounter = 1                   //目前有1个
	
	return rt
}

func Execute(runtime *Runtime, limit int) {
	/*
	递归地执行运行时
	*/
}
