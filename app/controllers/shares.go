package controllers

import (
	"net/http"
	"strconv"

	"github.com/chalvern/apollo/app/model"
	"github.com/chalvern/apollo/app/service"
	"github.com/chalvern/sugar"
	"github.com/gin-gonic/gin"
)

// ShareDetailGet 点击一个分享
func ShareDetailGet(c *gin.Context) {
	c.Set(PageTitle, "分享详情")
	share, err := service.ShareQueryByID(c.Query("id"))
	if err != nil || share.ID == 0 {
		sugar.Warnf("ShareGet 未检索到对应的分享，id=%v", c.Query("id"))
		html(c, http.StatusOK, "notify/error.tpl", gin.H{
			"Info": "对应条目找不到了，请看看别的吧",
		})
		return
	}
	share.Click(share.ID)
	htmlOfOk(c, "shares/detail.tpl", gin.H{
		"Share": share,
	})
}

// ShareNewGet 创建分享的表单
func ShareNewGet(c *gin.Context) {
	c.Set(PageTitle, "创建分享")
	htmlOfOk(c, "shares/new.tpl", gin.H{})
}

// ShareNewPost 创建分享
func ShareNewPost(c *gin.Context) {
	c.Set(PageTitle, "创建分享")
	form := struct {
		URL    string `form:"url" binding:"required,url"`
		Title  string `form:"title" binding:"required,lengte=3"`
		Review string `form:"review" binding:"required,lenlte=30"`
		Tag    string `form:"tag"`
	}{}

	if errs := c.ShouldBind(&form); errs != nil {
		sugar.Warnf("ShareNewPost Bind form Error: %s", errs.Error())
		html(c, http.StatusOK, "shares/new.tpl", gin.H{
			"FlashError": "请检查URL、标题、评论等内容及格式是否填写正确",
		})
		return
	}

	user := c.MustGet("user").(*model.User)
	share := model.Share{
		UserID: user.ID,
		URL:    form.URL,
		Title:  form.Title,
		Review: form.Review,
		Tag:    form.Tag,
	}

	if err := service.ShareCreate(&share); err != nil {
		sugar.Errorf("ShareNewPost Create Error: %s", err.Error())
		htmlOfOk(c, "shares/new.tpl", gin.H{
			"FlashError": "保存出现了问题，请检查提交内容后稍后重试",
		})
		return
	}
	htmlOfOk(c, "notify/success.tpl", gin.H{
		"Info":         "发布成功!",
		"Timeout":      3,
		"RedirectURL":  "/user/info/" + strconv.Itoa(int(user.ID)),
		"RedirectName": "用户页",
	})
}