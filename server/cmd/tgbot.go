package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

// 保存用户输入和操作步骤
var userInputMap = make(map[int64]string) // 保存用户输入的内容
var userStepMap = make(map[int64]string)  // 保存用户当前操作步骤

// 定义菜单
const (
	// 一级菜单
	menuMain = "代理管理|域名管理|客服管理"
	// 代理管理菜单
	subMenuProxy = "代理列表|创建代理|删除代理|充值代理"
	// 域名管理菜单
	subMenuDomain = "域名列表|批量创建|删除域名"
	// 客服管理菜单
	subMenuSupport = "创建客服|充值|扣除|客服列表|搜索"
)

func init() {
	var err error
	// 直接将 Token 写在代码中
	bot, err = tgbotapi.NewBotAPI("7667408957:AAH5T-e8TdJ8a7deaC1cAssCSl-s-F5hoFk") // 直接写入你的 Token
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		// 获取消息并判断是菜单点击
		if update.Message != nil {
			text := update.Message.Text
			chatID := update.Message.Chat.ID
			msg := tgbotapi.NewMessage(chatID, "")

			switch {
			case text == "/start":
				// 显示主菜单，使用按钮代替文本菜单
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
				bot.Send(msg)
			}
		}

		// 处理按钮点击的回调
		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID
			msg := tgbotapi.NewMessage(chatID, "")

			// 确认用户点击的按钮
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			bot.AnswerCallbackQuery(callback) // 给 Telegram 发送响应，避免点击后没有反馈

			// 根据点击的按钮执行操作
			switch callbackData {
			case "proxy":
				msg.Text = "您已选择代理管理操作。请选择一个具体操作："
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

			case "domain":
				msg.Text = "您已选择域名管理操作。请选择一个具体操作："
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("域名列表", "domain_list"),
						tgbotapi.NewInlineKeyboardButtonData("批量创建", "batch_create"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("删除域名", "delete_domain"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("返回首页", "main_menu"),
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
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("扣除", "deduct"),
						tgbotapi.NewInlineKeyboardButtonData("客服列表", "support_list"),
					),
					tgbotapi.NewInlineKeyboardRow(
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
			}
		}

		// 处理用户输入的内容（记录并反馈）
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			text := update.Message.Text

			// 根据用户当前步骤继续操作
			switch userStepMap[chatID] {
			case "creating_support":
				// 输入客服账号数量
				accountCount := text
				userInputMap[chatID] = accountCount
				msg := tgbotapi.NewMessage(chatID, "请输入充值天数：")
				userStepMap[chatID] = "enter_recharge_days"
				bot.Send(msg)

			case "enter_recharge_days":
				// 输入充值天数
				rechargeDays := text
				accountCount := userInputMap[chatID]
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("客服账号创建成功！\n创建了 %s 个账号，充值了 %s 天", accountCount, rechargeDays))
				bot.Send(msg)
				// 清除用户记录
				userStepMap[chatID] = ""  // 清空步骤
				userInputMap[chatID] = "" // 清空输入

			case "recharge":
				// 输入客服账号
				accountID := text
				userInputMap[chatID] = accountID
				msg := tgbotapi.NewMessage(chatID, "请输入充值天数：")
				userStepMap[chatID] = "recharge_days"
				bot.Send(msg)

			case "recharge_days":
				// 输入充值天数
				rechargeDays := text
				accountID := userInputMap[chatID]
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("客服账号 %s 充值了 %s 天成功", accountID, rechargeDays))
				bot.Send(msg)
				// 清除用户记录
				userStepMap[chatID] = ""  // 清空步骤
				userInputMap[chatID] = "" // 清空输入

			case "deduct":
				// 输入客服账号
				accountID := text
				userInputMap[chatID] = accountID
				msg := tgbotapi.NewMessage(chatID, "请输入扣除天数：")
				userStepMap[chatID] = "deduct_days"
				bot.Send(msg)

			case "deduct_days":
				// 输入扣除天数
				deductDays := text
				accountID := userInputMap[chatID]
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("客服账号 %s 扣除 %s 天成功", accountID, deductDays))
				bot.Send(msg)
				// 清除用户记录
				userStepMap[chatID] = ""  // 清空步骤
				userInputMap[chatID] = "" // 清空输入

			case "batch_create":
				// 输入域名类型
				if text == "1" || text == "2" {
					domainType := "入口"
					if text == "2" {
						domainType = "落地"
					}
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("您选择了 %s 域名，请输入域名列表（多个域名请换行输入）：", domainType))
					userStepMap[chatID] = "enter_domains"
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(chatID, "请输入有效的选择：1 为入口，2 为落地。")
					bot.Send(msg)
				}

			case "enter_domains":
				// 输入域名列表
				domains := text
				domainList := strings.Split(domains, "\n")
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("共 %d 个 %s 域名：\n%s", len(domainList), text, strings.Join(domainList, "\n")))
				bot.Send(msg)
				// 清除用户记录
				userStepMap[chatID] = ""  // 清空步骤
				userInputMap[chatID] = "" // 清空输入
			}
		}
	}
}
