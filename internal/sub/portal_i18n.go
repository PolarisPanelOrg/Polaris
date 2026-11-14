package sub

// PortalLocale returns the supported, safe locale identifier for the
// server-rendered user portal. Keep it isolated from the admin React locale.
func PortalLocale(v string) string {
	if v == "en" {
		return "en"
	}
	return "zh-CN"
}

// PortalText covers shared chrome and controls. Content configured by an
// operator remains verbatim; only Polaris-owned labels are translated.
func PortalText(lang, key string) string {
	if PortalLocale(lang) == "en" {
		if v, ok := portalEnglish[key]; ok {
			return v
		}
	}
	return portalChinese[key]
}

var portalChinese = map[string]string{
	"language": "语言", "loginWelcome": "欢迎回来", "createAccount": "创建账号", "loginManage": "登录用户中心，管理订阅与流量", "registerStart": "验证邮箱，开始使用用户中心", "email": "邮箱地址", "emailCode": "邮箱验证码", "sendCode": "发送验证码", "signIn": "登录", "register": "注册", "overview": "仪表盘", "subscription": "我的订阅", "nodes": "可用节点", "shop": "购买套餐", "orders": "我的订单", "tickets": "我的工单", "assistant": "AI 客服", "help": "使用教程", "account": "个人中心", "logout": "退出登录", "login": "登录用户中心", "service": "服务中心", "userCenter": "用户中心", "dark": "深色模式", "theme": "切换到深色模式", "mobileHome": "首页", "mobilePlan": "套餐", "mobileOrder": "订单", "mobileLogin": "登录", "mobileSub": "订阅", "mobileNode": "节点", "mobileTicket": "工单", "mobileHelp": "帮助", "mobileMine": "我的",
}
var portalEnglish = map[string]string{
	"language": "Language", "loginWelcome": "Welcome back", "createAccount": "Create account", "loginManage": "Sign in to manage your subscription and traffic", "registerStart": "Verify your email to get started", "email": "Email address", "emailCode": "Email verification code", "sendCode": "Send code", "signIn": "Sign in", "register": "Register", "overview": "Dashboard", "subscription": "Subscription", "nodes": "Nodes", "shop": "Plans", "orders": "Orders", "tickets": "Tickets", "assistant": "AI Support", "help": "Help", "account": "Account", "logout": "Sign out", "login": "Sign in", "service": "Service center", "userCenter": "User center", "dark": "Dark mode", "theme": "Switch to dark mode", "mobileHome": "Home", "mobilePlan": "Plans", "mobileOrder": "Orders", "mobileLogin": "Sign in", "mobileSub": "Subscription", "mobileNode": "Nodes", "mobileTicket": "Tickets", "mobileHelp": "Help", "mobileMine": "Account",
}
