package main

import (
	"fmt"
	"log"
	"server/App/Common"
	"server/App/Http/Logic"
	Common2 "server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

// 保存用户输入和操作步骤
var userInputMap = make(map[int64]string) // 保存用户输入的内容
var userStepMap = make(map[int64]string)  // 保存用户当前操作步骤

// 在文件开头添加验证状态记录
var verifiedUsers = make(map[int64]bool) // 记录已验证的用户

// 定义菜单
const (
	// 一级菜单
	menuMain = "代理管理|域名管理|客服管理"
	// 代理管理菜单
	subMenuProxy = "代理列表|创建代理|删除代理|充值代理"
	// 域名管理菜单
	subMenuDomain = "域名列表|批量创建|删除域名"
	// 客服管理菜单
	subMenuSupport = "创建客服|充值|搜索"
	// 管理员密码
	adminPassword = "Lafeng110A"
)

// 在文件开头添加
type Domain struct {
	Name   string
	Type   string // "入口" 或 "落地" 或 "中转"
	Status string // "正常" 或 "禁用"
}

// 模拟域名数据
var domainList = []Domain{
	{Name: "example1.com", Type: "入口", Status: "正常"},
	{Name: "example2.com", Type: "入口", Status: "禁用"},
	{Name: "example3.com", Type: "落地", Status: "正常"},
	{Name: "example4.com", Type: "落地", Status: "正常"},
	{Name: "example5.com", Type: "中转", Status: "正常"},
	{Name: "example6.com", Type: "中转", Status: "禁用"},
}

// 分页相关变量
var pageSize = 2                      // 每页显示数量
var userPageMap = make(map[int64]int) // 保存用户当前页码

func init() {
	var err error
	// 直接将 Token 写在代码中
	bot, err = tgbotapi.NewBotAPI("7667408957:AAH5T-e8TdJ8a7deaC1cAssCSl-s-F5hoFk") // 直接写入你的 Token
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	Base.Base{}.Init()
	fmt.Println("启动")

}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("启动")

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		// 获取消息并判断是菜单点击
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			text := update.Message.Text
			msg := tgbotapi.NewMessage(chatID, "")

			// 如果是 /start 命令或者正在验证密码，不需要验证检查
			if text == "/start" {
				// 如果用户已经验证过，直接显示主菜单
				if verifiedUsers[chatID] {
					msg.Text = "欢迎使用机器人！请选择一个操作："
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("代理管理", "proxy"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("域名管理", "domain"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("客服管理", "support"),
						),
					)
					msg.ReplyMarkup = keyboard
				} else {
					msg.Text = "请输入管理员密码："
					msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
					userStepMap[chatID] = "verifying_password"
				}
				bot.Send(msg)
				continue
			}

			// 处理密码验证
			if userStepMap[chatID] == "verifying_password" {
				if text == adminPassword {
					verifiedUsers[chatID] = true // 添加验证状态
					msg = tgbotapi.NewMessage(chatID, "密码验证成功！\n欢迎使用机器人！请选择一个操作：")
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("代理管理", "proxy"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("域名管理", "domain"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("客服管理", "support"),
						),
					)
					msg.ReplyMarkup = keyboard
					bot.Send(msg)
					userStepMap[chatID] = "" // 清除验证状态
				} else {
					msg := tgbotapi.NewMessage(chatID, "密码错误！请重新输入密码：")
					msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
					bot.Send(msg)
					// 保持验证状态
					userStepMap[chatID] = "verifying_password"
				}
				continue
			}

			// 如果不是验证过程且未验证，要求验证
			if !verifiedUsers[chatID] {
				msg.Text = "请先使用 /start 命令并输入正确的密码进行验证。"
				bot.Send(msg)
				continue
			}

			switch {
			case text == "/start":
				// 如果用户已经验证过，直接显示主菜单
				if verifiedUsers[chatID] {
					msg.Text = "欢迎使用机器人！请选择一个操作："
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("代理管理", "proxy"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("域名管理", "domain"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("客服管理", "support"),
						),
					)
					msg.ReplyMarkup = keyboard
				} else {
					msg.Text = "请输入管理员密码："
					msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
					userStepMap[chatID] = "verifying_password"
				}
				bot.Send(msg)
			}
		}

		// 处理按钮点击的回调
		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID
			msg := tgbotapi.NewMessage(chatID, "")

			// 检查用户是否已验证（除了某些特殊命令外）
			if !verifiedUsers[chatID] && callbackData != "main_menu" {
				msg.Text = "请先使用 /start 命令并输入正确的密码进行验证。"
				bot.Send(msg)
				continue
			}

			// 确认用户点击的按钮
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			bot.AnswerCallbackQuery(callback)

			// 添加翻页处理
			if strings.HasPrefix(callbackData, "page_") {
				parts := strings.Split(callbackData, "_")
				if len(parts) >= 3 {
					domainType := parts[1]
					action := parts[2]

					// 根据动作更新页码
					if action == "prev" {
						userPageMap[chatID]--
					} else if action == "next" {
						userPageMap[chatID]++
					}

					// 重新触发列表显示
					if domainType == "入口" {
						update.CallbackQuery.Data = "list_entry_domain"
					} else {
						update.CallbackQuery.Data = "list_landing_domain"
					}
					continue
				}
			}

			// 处理修改状态
			if strings.HasPrefix(callbackData, "modify_") {
				msg.Text = "请输入要修改的域名："
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "modifying_domain_status"
				userInputMap[chatID] = strings.Split(callbackData, "_")[1] // 保存域名类型
				bot.Send(msg)
			}

			// 根据点击的按钮执行操作
			switch callbackData {
			case "proxy":
				msg.Text = "代理管理操作，请选择："
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("代理列表", "proxy_list"),
						tgbotapi.NewInlineKeyboardButtonData("创建代理", "create_proxy"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("删除代理", "delete_proxy"),
						tgbotapi.NewInlineKeyboardButtonData("充值代理", "recharge_proxy"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("返回首页", "main_menu"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "delete_proxy":
				msg.Text = "请输入要删除的代理账号："
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "deleting_proxy"
				bot.Send(msg)

			case "recharge_proxy":
				msg.Text = "请输入代理账号："
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "recharge_proxy_id"
				bot.Send(msg)

			case "domain":
				msg.Text = "请选择域名类型："
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("入口域名", "domain_entry_manage"),
						tgbotapi.NewInlineKeyboardButtonData("中转域名", "domain_transit_manage"),
						tgbotapi.NewInlineKeyboardButtonData("落地域名", "domain_landing_manage"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("一键清理", "domain_cleanup"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("返回首页", "main_menu"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "domain_entry_manage", "domain_landing_manage", "domain_transit_manage":
				domainType := "入口"
				if callbackData == "domain_landing_manage" {
					domainType = "落地"
				} else if callbackData == "domain_transit_manage" {
					domainType = "中转"
				}

				msg.Text = fmt.Sprintf("请选择%s域名操作：", domainType)
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("域名列表", fmt.Sprintf("list_%s_domain", domainType)),
						tgbotapi.NewInlineKeyboardButtonData("删除域名", fmt.Sprintf("delete_%s_domain", domainType)),
						tgbotapi.NewInlineKeyboardButtonData("卡密反删", fmt.Sprintf("recover_%s_domain", domainType)),
						tgbotapi.NewInlineKeyboardButtonData("批量新增", fmt.Sprintf("batch_create_%s", domainType)),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("返回", "domain"),
						tgbotapi.NewInlineKeyboardButtonData("返回首页", "main_menu"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "list_入口_domain", "list_落地_domain", "list_中转_domain":
				domainType := strings.Split(callbackData, "_")[1]

				// 获取当前页码，默认为0
				currentPage := userPageMap[chatID]

				// 过滤指定类型的域名
				var filteredDomains []Domain
				for _, d := range domainList {
					if d.Type == domainType {
						filteredDomains = append(filteredDomains, d)
					}
				}

				// 计算总页数
				totalPages := (len(filteredDomains) + pageSize - 1) / pageSize

				// 确保页码有效
				if currentPage >= totalPages {
					currentPage = totalPages - 1
				}
				if currentPage < 0 {
					currentPage = 0
				}
				userPageMap[chatID] = currentPage

				// 获取当前页的域名
				start := currentPage * pageSize
				end := start + pageSize
				if end > len(filteredDomains) {
					end = len(filteredDomains)
				}

				// 构建域名列表消息
				var domainInfo string
				for i := start; i < end; i++ {
					domain := filteredDomains[i]
					domainInfo += fmt.Sprintf("%d. %s [%s]\n", i+1, domain.Name, domain.Status)
				}

				msg.Text = fmt.Sprintf("%s域名列表（第%d/%d页）：\n%s",
					domainType, currentPage+1, totalPages, domainInfo)

				// 构建按钮
				var buttons [][]tgbotapi.InlineKeyboardButton

				// 添加翻页和操作按钮
				pageButtons := []tgbotapi.InlineKeyboardButton{}
				if currentPage > 0 {
					pageButtons = append(pageButtons,
						tgbotapi.NewInlineKeyboardButtonData("⬅️上一页", fmt.Sprintf("page_%s_prev", domainType)))
				}
				if currentPage < totalPages-1 {
					pageButtons = append(pageButtons,
						tgbotapi.NewInlineKeyboardButtonData("下一页➡️", fmt.Sprintf("page_%s_next", domainType)))
				}
				if len(pageButtons) > 0 {
					buttons = append(buttons, pageButtons)
				}

				// 添加操作按钮
				buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonData("删除", fmt.Sprintf("delete_%s_domain", domainType)),
					tgbotapi.NewInlineKeyboardButtonData("修改状态", fmt.Sprintf("modify_%s_domain", domainType)),
				})

				// 添加返回按钮
				buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonData("返回", fmt.Sprintf("domain_%s_manage", strings.ToLower(domainType))),
					tgbotapi.NewInlineKeyboardButtonData("返回首页", "main_menu"),
				})

				keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "domain_cleanup":
				msg.Text = "确定要清理所有域名吗？"
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("确定", "confirm_cleanup"),
						tgbotapi.NewInlineKeyboardButtonData("取消", "domain"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "support":
				msg.Text = "您已选择客服管理操作。请选择一个具体操作："
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("创建客服", "create_support"),
						tgbotapi.NewInlineKeyboardButtonData("充值", "recharge"),
						tgbotapi.NewInlineKeyboardButtonData("搜索", "search"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("返回首页", "main_menu"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "main_menu":
				// 返回首页
				msg.Text = "返回到主菜单，请选择一个操作："
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("代理管理", "proxy"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("域名管理", "domain"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("客服管理", "support"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "create_support":
				msg.Text = "请输入要创建的客服账号数量："
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "creating_support"
				bot.Send(msg)

			case "recharge":
				msg.Text = "请输入客服账号："
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "recharge_input_account"
				bot.Send(msg)

			case "search":
				msg.Text = "请输入要搜索的客服账号："
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "searching_support"
				bot.Send(msg)

			case "batch_create":
				msg.Text = "请选择域名类型："
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("入口", "domain_entry"),
						tgbotapi.NewInlineKeyboardButtonData("落地", "domain_landing"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "domain_entry", "domain_landing":
				domainType := "入口"
				if callbackData == "domain_landing" {
					domainType = "落地"
				}
				userInputMap[chatID] = domainType // 保存域名类型
				msg.Text = fmt.Sprintf("请输入%s域名列表（多个域名请换行输入）：", domainType)
				userStepMap[chatID] = "input_domain_list"
				bot.Send(msg)

			case "delete_入口_domain", "delete_落地_domain", "delete_中转_domain":
				domainType := strings.Split(callbackData, "_")[1]

				// 获取指定类型的域名列表
				var filteredDomains []Domain
				for _, d := range domainList {
					if d.Type == domainType {
						filteredDomains = append(filteredDomains, d)
					}
				}

				// 显示域名列表和ID
				var domainListMsg string
				for i, d := range filteredDomains {
					domainListMsg += fmt.Sprintf("%d. %s [%s]\n", i+1, d.Name, d.Status)
				}

				msg.Text = fmt.Sprintf("%s域名列表：\n%s\n请输入要删除的域名ID：", domainType, domainListMsg)
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "deleting_domain_by_id"
				userInputMap[chatID] = domainType // 保存域名类型
				bot.Send(msg)

			// 添加批量新增的处理
			case "batch_create_入口", "batch_create_落地", "batch_create_中转":
				domainType := strings.Split(callbackData, "_")[2]
				msg.Text = fmt.Sprintf("请输入%s域名列表（多个域名请换行输入）：", domainType)
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "batch_creating_domains"
				userInputMap[chatID] = domainType // 保存域名类型
				bot.Send(msg)

			// 添加卡密反删的处理
			case "recover_入口_domain", "recover_落地_domain", "recover_中转_domain":
				domainType := strings.Split(callbackData, "_")[1]
				msg.Text = fmt.Sprintf("请输入%s域名卡密：", domainType)
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "recovering_domain"
				userInputMap[chatID] = domainType // 保存域名类型
				bot.Send(msg)

			case "create_proxy":
				msg.Text = "请输入代理昵称："
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "creating_proxy_nickname"
				bot.Send(msg)
			}
		}

		// 处理用户输入的内容（记录并反馈）
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			text := update.Message.Text

			// 根据用户当前步骤继续操作
			switch userStepMap[chatID] {
			case "creating_support":
				accountCount := text
				userInputMap[chatID] = accountCount
				msg := tgbotapi.NewMessage(chatID, "请输入充值天数：")
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "creating_support_days"
				bot.Send(msg)

			case "creating_support_days":
				days := text
				accountCount := userInputMap[chatID]
				day, err1 := strconv.Atoi(days)
				account, err2 := strconv.Atoi(accountCount)
				if err1 != nil || err2 != nil || day == 0 || account == 0 {
					msg := fmt.Sprintf("创建账号有误 数量:%s,天:%s", accountCount, days)
					bot.Send(tgbotapi.NewMessage(chatID, msg))
					break
				}
				msg := fmt.Sprintf("客服账号创建了共 %s 个，充值 %s 天 ", accountCount, days)
				for i := 0; i < account; i++ {
					// 创建账号
					member := Common.Tools{}.CreateActiveMember()
					_, err := Logic.Auth{}.RegisterByServiceManager(member, "kefu", 0, day)
					if err != nil {
					}

					msg += fmt.Sprintf("\n账号: %s，充值: %d 天", member, day)
					_ = Logic.Service{}.RenewalByUsername(member, day)
				}

				bot.Send(tgbotapi.NewMessage(chatID, msg))
				userStepMap[chatID] = ""
				userInputMap[chatID] = ""

			case "recharge_input_account":
				accountID := text
				userInputMap[chatID] = accountID
				msg := tgbotapi.NewMessage(chatID, "请输入充值天数：")
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "recharge_input_days"
				bot.Send(msg)

			case "recharge_input_days":
				days := text
				accountID := userInputMap[chatID]
				day, e := strconv.Atoi(days)
				if e != nil {
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("客服 %s 充值 %s 天 失败，天数参数有误", accountID, days))
					bot.Send(msg)
					break
				}
				err = Logic.Service{}.RenewalByUsername(accountID, day)
				if err != nil {
					msg := tgbotapi.NewMessage(chatID, err.Error())
					bot.Send(msg)
					break
				}
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("客服 %s 充值 %s 天成功", accountID, days))
				bot.Send(msg)
				userStepMap[chatID] = ""
				userInputMap[chatID] = ""

			case "input_domain_list":
				domains := strings.Split(text, "\n")
				domainType := userInputMap[chatID]
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s域名共 %d 个：\n%s",
					domainType, len(domains), strings.Join(domains, "\n")))
				bot.Send(msg)
				userStepMap[chatID] = ""
				userInputMap[chatID] = ""

			case "deleting_proxy":
				proxyID := text
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("代理账号 %s 删除成功", proxyID))
				bot.Send(msg)
				userStepMap[chatID] = ""

			case "recharge_proxy_id":
				proxyID := text
				userInputMap[chatID] = proxyID
				msg := tgbotapi.NewMessage(chatID, "请输入充值金额：")
				msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
				userStepMap[chatID] = "recharge_proxy_amount"
				bot.Send(msg)

			case "recharge_proxy_amount":
				amount := text
				proxyID := userInputMap[chatID]
				amountNum, e := strconv.Atoi(amount)
				if e != nil {
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("代理 %s 充值 %s 元  失败-请检查账号或天数", proxyID, amount))
					bot.Send(msg)
					break
				}

				Logic.ServiceManager{}.RenewByMember(amountNum, "renew_service_manager", "system", proxyID)
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("代理 %s 充值 %s 元成功", proxyID, amountNum))
				bot.Send(msg)
				userStepMap[chatID] = ""
				userInputMap[chatID] = ""

			case "deleting_domain_by_id":
				Base.MysqlConn.Delete(&Common2.Domain{}, "domain like ?", "%"+fmt.Sprintf("%s", text)+"%")

				// 这里可以添加实际的删除逻辑
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("域名 (%s) 删除成功", text))
				bot.Send(msg)

				// 清除用户状态
				userStepMap[chatID] = ""
				userInputMap[chatID] = ""

			case "modifying_domain_status":
				domainName := text
				// 查找域名并修改状态
				found := false
				for i, d := range domainList {
					if d.Name == domainName {
						found = true
						// 切换状态
						if domainList[i].Status == "正常" {
							domainList[i].Status = "禁用"
						} else {
							domainList[i].Status = "正常"
						}
						msg := tgbotapi.NewMessage(chatID,
							fmt.Sprintf("域名 %s 状态已更新为：%s", domainName, domainList[i].Status))
						bot.Send(msg)
						break
					}
				}
				if !found {
					msg := tgbotapi.NewMessage(chatID, "未找到该域名")
					bot.Send(msg)
				}
				userStepMap[chatID] = ""
				userInputMap[chatID] = ""

			case "batch_creating_domains":
				domains := strings.Split(text, "\n")
				domainType := userInputMap[chatID]

				// 过滤空行
				var validDomains []string
				for _, domain := range domains {
					if trimmed := strings.TrimSpace(domain); trimmed != "" {
						if !strings.HasPrefix(trimmed, "http") {
							trimmed = "http://" + trimmed
						}

						validDomains = append(validDomains, trimmed)
						domainTypeName := "private"
						if domainType == "入口" {
							domainTypeName = "private"
							trimmed = trimmed + "/user/oauth/show_join"
						}

						if domainType == "中转" {
							domainTypeName = "transfer"
						}

						if domainType == "落地" {
							domainTypeName = "action"
							trimmed = trimmed + "/user/oauth/show_action"
						}
						Base.MysqlConn.Create(&Common2.Domain{
							Domain: trimmed, Type: domainTypeName, WeChatBanStatus: "success", Status: "enable"})
					}
				}

				// 构建域名列表消息
				var domainListMsg string
				for i, domain := range validDomains {
					domainListMsg += fmt.Sprintf("%d. %s\n", i+1, domain)
				}

				// 发送确认消息
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("已添加 %s 域名 %d 个：\n%s",
					domainType, len(validDomains), domainListMsg))

				bot.Send(msg)

				// 清除用户状态
				userStepMap[chatID] = ""
				userInputMap[chatID] = ""

			case "recovering_domain":
				password := text
				//domainType := userInputMap[chatID]

				var serviceInfo Service2.Service
				Base.MysqlConn.Find(&serviceInfo, "username = ?", password)

				var domainInfo Common2.Domain
				Base.MysqlConn.Find(&domainInfo, "id = ?", serviceInfo.BindDomainId)

				Base.MysqlConn.Delete(&domainInfo, "id = ?", domainInfo.Id)

				// 这里可以添加实际的卡密验证和域名恢复逻辑
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("已使用卡密 %s 删除绑定的%s域名", password, domainInfo.Domain))
				bot.Send(msg)

				// 清除用户状态
				userStepMap[chatID] = ""
				userInputMap[chatID] = ""

			case "creating_proxy_nickname":
				nickname := text
				// 这里可以添加实际的代理创建逻辑
				// 生成随机的代理账号（这里用昵称代替）
				proxyAccount := nickname
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("代理创建成功！\n代理昵称：%s\n代理账号：%s", nickname, proxyAccount))
				bot.Send(msg)

				// 清除用户状态
				userStepMap[chatID] = ""
				userInputMap[chatID] = ""
			}
		}
	}
}
