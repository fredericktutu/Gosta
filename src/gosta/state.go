package gosta

type Op int
const (
	OP_SEND Op = iota
	OP_RECV 
	OP_EPSILON          // 空边
)

type NodeType int
const (
	NODE_START NodeType = iota
	NODE_END
	NODE_MIDDLE
)

type ChanID int        // Channel id
type GID int           // GFSM id
const (
	NO_MACHINE GID = -1  // 对于RGoroutine,从未初始化的意思
)

type SID struct {      // State id <which GFSM, which State>
	GID GID
	ID int
}

// 协程
// 自动机的Go调用 & Goroutine运行时状态
type Goroutine struct {
	Num int           // Goroutine序号
	Machine GID       // 是从哪个函数中spawn出来的
	Bindings []ChanID  // channel的id，和StaticProgram里的Channels下标一致
	//Current SID        // 当前的状态，只在运行时有用
}


// 协程自动机
type GFSM struct {
	ID GID            // 在StaticProgram中自动机的下标
	States []*State    
}

// 协程自动机中的状态
type State struct {
	ID SID			  // GFSM中状态的下标
	Gos []*Goroutine   // 在这里即将调用的Goroutine，在运行到这里时，会复制       
	In, Out []*Edge
	Type NodeType
}


type Edge struct {
	From, To SID
	Op Op
	ChanID ChanID
}

type Channel struct {
	ID ChanID
	Cap int
	Type string
}


// 构建时，一个BasicBlock对应的自动机
type SmallFSM struct {
	First *State
	Last *State
	Outs []int
}


func GetChan(self StaticProgram, cid ChanID) *Channel{
	return &(self.Channels[cid])
}

/*
func GetState(self StaticProgram, sid SID) *State {
	return &(self.GFSMs[sid.GID].States[sid.ID])
}
*/
/*
func GetGFSM(self StaticProgram, gid GID) *GFSM {
	return &(self.GFSMs[gid])
}
*/



//begin SmallFSM构建方法

func LinkState(from *State, to *State, chanID ChanID, op Op) {
	/*
	linkState
	顾名思义，连接两个State
	在构建SmallFSM和构建bootstrap GFSM时用到
	*/
	var e *Edge = new(Edge)
	e.From = from.ID
	e.To = to.ID
	e.Op = op
	if op != OP_EPSILON {
		e.ChanID = chanID
	}
	from.Out = append(from.Out, e)
	to.In = append(to.In, e)
}

func MakeSmallFSM(gfsm *GFSM, gid int, sid int) *SmallFSM {
	/*
	创建一个SmallFSM
	*/
	var sfsm *SmallFSM = new(SmallFSM)

	var startState *State = new(State)
	gfsm.States = append(gfsm.States, startState)
	
	startState.Type = NODE_MIDDLE
	startState.ID = SID {GID(gid), sid}

	sfsm.First = startState
	sfsm.Last = sfsm.First

	return sfsm
}

func AddState(self *SmallFSM, gfsm *GFSM, op Op, chanID ChanID, gid int, sid int) {
	/*
	对一个SmallFSM追加一个State
	op,chanID 代表是通过哪个channel上的哪个操作
	gid,sid用于给新增的状态分配id
	*/
	var nextState *State = new(State)
	gfsm.States = append(gfsm.States, nextState)

	nextState.Type = NODE_MIDDLE
	nextState.ID = SID {GID(gid), sid}

	LinkState(self.Last, nextState, chanID, op)
	self.Last = nextState
}

func AddGoroutine(self *SmallFSM) *Goroutine {
	var g *Goroutine = new(Goroutine)
	self.Last.Gos = append(self.Last.Gos, g)
	return g 
}
