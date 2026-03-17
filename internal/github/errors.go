package github

import "errors"

// 预定义错误，便于调用方进行类型判断
var (
	ErrUnauthorized = errors.New("github: token 无效或已过期")
	ErrNotFound     = errors.New("github: 资源不存在或仓库为私有")
	ErrRateLimited  = errors.New("github: API 请求已限流，请稍后重试")
	ErrNetworkError = errors.New("github: 网络请求失败")
)

// IsNotFound 判断错误是否为"未找到"类型
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsUnauthorized 判断错误是否为"未授权"类型
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsRateLimited 判断错误是否为"限流"类型
func IsRateLimited(err error) bool {
	return errors.Is(err, ErrRateLimited)
}
