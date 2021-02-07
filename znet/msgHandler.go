package znet

import (
	"fmt"
	"strconv"

	"e.coding.net/646222472/framework/zinc/utils"
	"e.coding.net/646222472/framework/zinc/ziface"
)

// MsgHandler 消息处理模块的实现
type MsgHandler struct {
	// 存放每个 MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责 Worker 取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作 Worker 池的 worker 数量
	WorkerPoolSize uint32
}

// NewMsgHandler 初始化/创建 MsgHandler 方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

// DoMsgHandler 调度/执行对应的 Router 消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 从 Request 中找到 MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("Api msgID=%d is NOT FOUND! Need Register\n", request.GetMsgID())
	}

	// 根据 MsgID 调度对应的 Router 业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1.判断当前 MsgID 绑定的 API 处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id 已经注册了
		panic("repeat api, MsgID=" + strconv.Itoa(int(msgID)))
	}

	// 2.添加 MsgID 和 API 的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " succ!")
}

// StartWorkerPool 启动一个 Worker 工作池（开启工作池的方法只能发生一次，一个框架只能有一个 Worker 工作池）
func (mh *MsgHandler) StartWorkerPool() {
	// 根据 WorkerPoolSize 分别开启 Worker，每个 Worker 用一个 Goroutine 来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个 Worker 被启动
		// 1、给当前Worker对应的channel消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2、启动当前的 Worker 工作流， 阻塞等待消息从 channel 中传递进来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个 Worker 工作流
func (mh *MsgHandler) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkID=", workerID, " is started ...")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息到来，出列的就是客户端 Request，执行当前 Request 所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 发送消息到任务队列 TaskQueue 中，由 Worker 进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1、将消息平均分配给不同的 Worker
	// 根据客户端建立的 ConnID 来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Printf(
		"Add ConnID=%d, request MsgID=%d to workerID=%d",
		request.GetConnection().GetConnID(), request.GetMsgID(), workerID,
	)
	// 2、将消息发送给 Worker 内的 TaskQueue
	mh.TaskQueue[workerID] <- request
}
