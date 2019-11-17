package whois

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	whois "github.com/likexian/whois-go"
	parser "github.com/likexian/whois-parser-go"

	"github.com/naiba/nbdomain/model"
)

//Whois whois 查询
func Whois(c *gin.Context) {
	domain := c.Param("domain")
	if len(domain) < 4 {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "域名格式不符合规范",
		})
		return
	}
	result, err := whois.Whois(domain)
	if err == nil {
		var parsed parser.WhoisInfo
		parsed, err = parser.Parse(result)
		if err == nil {
			c.JSON(http.StatusOK, model.Response{
				Code:   http.StatusOK,
				Result: parsed,
			})
			return
		}
	}
	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("获取 Whois 错误：%s", err.Error()),
	})
}
