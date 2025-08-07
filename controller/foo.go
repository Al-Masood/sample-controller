package controller

import (
	"context"
	"fmt"
	"time"

	clientset "github.com/al-masood/sample-controller/pkg/generated/clientset/versioned"
	informers "github.com/al-masood/sample-controller/pkg/generated/informers/externalversions/sample.com/v1alpha1"
	listers "github.com/al-masood/sample-controller/pkg/generated/listers/sample.com/v1alpha1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type Controller struct {
	kubeclientset   kubernetes.Interface
	sampleclientset clientset.Interface

	foosLister listers.FooLister
	foosSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
}

func NewController(
	ctx context.Context,
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	fooInformer informers.FooInformer) *Controller {

	controller := &Controller{
		kubeclientset:   kubeclientset,
		sampleclientset: sampleclientset,

		foosLister: fooInformer.Lister(),
		foosSynced: fooInformer.Informer().HasSynced,

		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ""),
	}

	fooInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.handleAdd,
		DeleteFunc: controller.handleDelete,
	})
	return controller
}

func (controller *Controller) handleAdd(obj interface{}) {
	controller.workqueue.Add(obj)
}

func (controller *Controller) handleDelete(obj interface{}) {
	controller.workqueue.Add(obj)
}

func (c *Controller) Run(ctx context.Context, workers int) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	if ok := cache.WaitForCacheSync(ctx.Done(), c.foosSynced); !ok {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runWorker, time.Second)
	}

	<-ctx.Done()
	klog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker(ctx context.Context) {
	for c.processNextWorkItem(ctx) {
	}
}

func (c *Controller) processNextWorkItem(ctx context.Context) bool {
	objRef, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	defer c.workqueue.Done(obj)

	key, ok := obj.(string)
	if !ok {
		c.workqueue.Forget(obj)
		runtime.HandleError(fmt.Errorf(obj))
		return true
	}

	if err := c.syncHandler(ctx, key); err != nil {
		runtime.HandleError(fmt.Errorf(key, err))
	}

	c.workqueue.Forget(obj)
	klog.Infof(key)
	return true
}

func (c *Controller) syncHandler(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespacesKey(key)
	if err != nil {
		return err
	}

	foo, err := c.foosLister.Foos(namespace).Get(name)
	if err != nil {
		klog.Infof(namespace, name)
		return nil
	}

	klog.Infof(namespace, name)

	return nil
}
