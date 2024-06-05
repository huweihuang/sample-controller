func main(){
	sc := NewSampleController()
	sc.Run()
}

func NewSampleController() *SampleController{
	// 初始化队列
	sc := &SampleController{
		kubeClient:       kubeClient,
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), queueName),
	}

	// 添加event handler
	scInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			sc.addSample(logger, obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			sc.updateSample(logger, oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			sc.deleteSample(logger, obj)
		},
	})

	sc.scListerSynced = scInformer.Informer().HasSynced
	// 添加syncHandler
	sc.syncHandler = sc.syncSample
}

func (sc *SampleController)Run(ctx context.Context, workers int) {
	defer sc.queue.ShutDown()

	// 同步cache
	if !cache.WaitForNamedCacheSync(sc.Kind, ctx.Done(), sc.scListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, sc.worker, time.Second)
	}

	<-ctx.Done()
}

func (sc *SampleController) worker(ctx context.Context) {
	for sc.processNextWorkItem(ctx) {
	}
}

func (sc *SampleController) processNextWorkItem(ctx context.Context) bool {
	key, quit := sc.queue.Get()
	if quit {
		return false
	}
	defer sc.queue.Done(key)

	err := sc.syncHandler(ctx, key.(string))
	if err == nil {
		sc.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("sync %q failed with %v", key, err))
	sc.queue.AddRateLimited(key)

	return true
}

func (sc *SampleController) syncSample() error{
	// 定制具体的业务逻辑
}
