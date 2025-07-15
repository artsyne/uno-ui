package di

import (
	"fmt"
	"reflect"
	"sync"

	errors "github.com/biuren/hebe/internal/core/error"
)

// FactoryFunc 工厂函数类型
type FactoryFunc func(ContainerInterface) (interface{}, error)

// provider 简化的提供者实现
type provider struct {
	serviceType reflect.Type
	factory     FactoryFunc
	singleton   bool
	instance    interface{}
	initOnce    sync.Once
	mu          sync.RWMutex
}

// newProvider 创建提供者
func newProvider(serviceType reflect.Type, factory FactoryFunc, singleton bool) *provider {
	return &provider{
		serviceType: serviceType,
		factory:     factory,
		singleton:   singleton,
	}
}

// provide 提供服务实例
func (p *provider) provide(c ContainerInterface) (interface{}, error) {
	if !p.singleton {
		return p.factory(c)
	}

	var err error
	p.initOnce.Do(func() {
		p.instance, err = p.factory(c)
	})

	if err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.instance, nil
}

// isSingleton 检查是否为单例
func (p *provider) isSingleton() bool {
	return p.singleton
}

// ValidateFactory 验证工厂函数签名
func ValidateFactory(factory interface{}) error {
	if factory == nil {
		return newInvalidFactoryError(reflect.TypeOf(factory))
	}

	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		return newInvalidFactoryError(factoryType)
	}

	if factoryType.NumIn() != 1 || factoryType.NumOut() != 2 {
		return newInvalidFactoryError(factoryType)
	}

	containerType := reflect.TypeOf((*ContainerInterface)(nil)).Elem()
	if !factoryType.In(0).Implements(containerType) {
		return newInvalidFactoryError(factoryType)
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if !factoryType.Out(1).Implements(errorType) {
		return newInvalidFactoryError(factoryType)
	}

	return nil
}

// WrapFactory 包装工厂函数
func WrapFactory(factory interface{}) (FactoryFunc, error) {
	if err := ValidateFactory(factory); err != nil {
		return nil, err
	}

	factoryValue := reflect.ValueOf(factory)

	return func(c ContainerInterface) (interface{}, error) {
		results := factoryValue.Call([]reflect.Value{reflect.ValueOf(c)})

		if err := results[1].Interface(); err != nil {
			return nil, err.(error)
		}

		return results[0].Interface(), nil
	}, nil
}

// newInvalidFactoryError 创建无效工厂函数错误
func newInvalidFactoryError(factoryType reflect.Type) error {
	return errors.DIInvalidFactoryErr().WithDetails(
		fmt.Sprintf("invalid factory function type: %s", factoryType.String()),
	)
}
