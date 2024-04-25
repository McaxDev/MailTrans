package handler

import (
	"sort"

	"github.com/McaxDev/MailTrans/util"
	"github.com/emersion/go-imap"
	"github.com/gin-gonic/gin"
)

type Email struct {
	Time    string
	UID     uint32
	From    string
	Subject string
	Preview string
}

// 获取所有邮件的handler
func AllMail(c *gin.Context) {

	// 与邮件服务器建立连接
	conn, err := util.ConnectMail()
	if err != nil {
		util.Error(c, 500, "与邮件服务器建立连接失败", err)
		return
	}
	defer conn.Logout()

	// 创建检索条件变量
	criteria := imap.NewSearchCriteria()

	/*
		// 通过配置文件里的关键词白名单过滤
		whitelist := config.Config.Filter
		if len(whitelist) > 0 {
			orCriteria := make([]*imap.SearchCriteria, len(whitelist))
			for i, subject := range whitelist {
				orCriteria[i] = imap.NewSearchCriteria()
				orCriteria[i].Header.Add("Subject", subject)
			}
			criteria.Or = orCriteria
		}
	*/

	// 从查询字符串参数获取收件人邮箱过滤条件
	receiver := c.Query("receiver")
	if receiver != "" {
		criteria.Header.Add("To", receiver)
	}

	// 根据条件搜索满足条件的邮件ID
	ids, err := conn.UidSearch(criteria)
	if err != nil {
		util.Error(c, 500, "从邮件服务器获取邮件UID列表失败", err)
		return
	}

	// 对UID进行排序并截取前五个元素
	sort.Slice(ids, func(i, j int) bool { return ids[i] > ids[j] })
	if len(ids) > 5 {
		ids = ids[:5]
	}

	// 获取邮件的概要信息
	messages, err := util.GetContent(conn, ids...)
	if err != nil {
		util.Error(c, 500, "获取邮件信息失败", err)
		return
	}

	// 将邮件列表变成切片
	var emails []Email
	for msg := range messages {
		if msg == nil {
			continue
		}
		preview, err := util.ExtractText(msg)
		if err != nil {
			preview = "加载邮件详细信息失败"
		}
		email := Email{
			Time:    msg.Envelope.Date.Format("2006-01-02 15:04"),
			UID:     msg.Uid,
			From:    msg.Envelope.From[0].PersonalName,
			Subject: msg.Envelope.Subject,
			Preview: preview,
		}
		emails = append(emails, email)
	}

	// 将切片返回
	c.AbortWithStatusJSON(200, util.Resp("查询成功", emails))
}
