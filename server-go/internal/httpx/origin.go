package httpx

import "net/http"

// Origin 根据 HTTP 请求头信息构建完整的请求来源 URL
// 优先使用 X-Forwarded-* 头（反向代理环境），否则使用请求的 Host 和 TLS 状态
//
// 适用场景：
//   - 应用部署在反向代理（Nginx、Caddy）后面时，代理会传递原始协议和主机信息
//   - 用于生成重定向 URL 或构建 API 响应中的完整链接
//
// 参数：
//   - r: HTTP 请求对象
//
// 返回值：
//   - string: 完整的来源 URL，格式为 "https://example.com" 或 "http://example.com"
func Origin(r *http.Request) string {
	// 优先从代理头获取协议（X-Forwarded-Proto 通常由代理设置）
	proto := r.Header.Get("X-Forwarded-Proto")
	// 优先从代理头获取主机（X-Forwarded-Host 通常由代理设置）
	host := r.Header.Get("X-Forwarded-Host")

	// 如果代理未提供主机信息，使用请求的 Host 头
	if host == "" {
		host = r.Host
	}

	// 如果代理未提供协议信息，根据 TLS 状态判断
	if proto == "" {
		if r.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}

	// 组合完整 URL
	return proto + "://" + host
}
