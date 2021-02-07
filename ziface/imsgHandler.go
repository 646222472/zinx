package ziface

// IMsgHandler 消息管理抽象层
type IMsgHandler interface {
	// 调度/执行对应的 Router 消息处理方法
	DoMsgHandler(IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(uint32, IRouter)

	// 启动 Worker 工作池
	StartWorkerPool()

	// 发送消息到任务队列 TaskQueue 中，由 Worker 进行处理
	SendMsgToTaskQueue(IRequest)
}
