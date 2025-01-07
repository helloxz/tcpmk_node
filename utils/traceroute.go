// utils/route.go
package utils

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pixelbender/go-traceroute/traceroute"
)

// TraceRouteResult 定义了单行路由跟踪结果的结构
type TraceRouteResult struct {
	Hop      int
	IP       net.IP
	RTTs     []time.Duration
	Reached  bool
	ErrorMsg string
}

// TraceRoute 发起路由跟踪并逐行返回结果。
// 参数:
//   - ip: 目标 IP 地址
//
// 返回值:
//   - resultChan: 只读的通道，用于接收每一行的跟踪结果
//   - error: 如果初始化失败，则返回错误
func TraceRoute(ip net.IP) (<-chan TraceRouteResult, error) {
	// 创建一个通道用于发送结果
	resultChan := make(chan TraceRouteResult)

	// 创建一个新的 Tracer 实例，配置 Count=1 避免重复发送
	tracer := &traceroute.Tracer{
		Config: traceroute.Config{
			Delay:    50 * time.Millisecond,
			Timeout:  2 * time.Second,
			MaxHops:  30,
			Count:    1, // 设置为1避免重复执行
			Networks: []string{"ip4:icmp", "ip4:ip"},
		},
	}

	// 创建一个上下文，可以根据需要设置超时或取消
	ctx, cancel := context.WithCancel(context.Background())

	// 启动一个 goroutine 执行路由跟踪
	go func() {
		defer close(resultChan) // 结束时关闭通道
		defer cancel()          // 确保取消上下文

		// 初始化一个空的结果映射，键为跳数，值为 *traceroute.Hop
		hopMap := make(map[int]*traceroute.Hop)
		var mu sync.Mutex  // 保护 hopMap 的互斥锁
		var once sync.Once // 确保只处理一次到达目标的情况
		reached := false   // 标记是否已经到达目标

		// 定义一个回调函数，用于处理每一个回复
		callback := func(reply *traceroute.Reply) {
			mu.Lock()
			defer mu.Unlock()

			// 如果已经到达目标且当前回复的跳数大于等于目标跳数，忽略后续回复
			if reached && reply.Hops >= hopMap[reply.Hops].Distance {
				return
			}

			hop, exists := hopMap[reply.Hops]
			if !exists {
				hop = &traceroute.Hop{
					Distance: reply.Hops,
				}
				hopMap[reply.Hops] = hop
			}
			node := hop.Add(reply)

			isReached := false
			if reply.IP.Equal(ip) {
				isReached = true
				once.Do(func() {
					reached = true
					cancel() // 到达目标后取消上下文，停止进一步探测
				})
			}

			result := TraceRouteResult{
				Hop:     reply.Hops,
				IP:      node.IP,
				RTTs:    node.RTT,
				Reached: isReached,
			}

			// 将结果发送到通道
			select {
			case resultChan <- result:
				// 成功发送
			default:
				// 如果通道已满，忽略或采取其他措施
			}
		}

		// 执行 tracer.Trace
		err := tracer.Trace(ctx, ip, callback)
		if err != nil && err != context.Canceled {
			// 如果发生错误，发送错误信息到通道
			resultChan <- TraceRouteResult{
				ErrorMsg: fmt.Sprintf("Traceroute error: %v", err),
			}
		}
	}()

	return resultChan, nil
}
