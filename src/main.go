package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lovusoft/salvation/src/entity"
	"github.com/unrolled/secure"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	log.Println("Test")
	db, err := gorm.Open("mysql", "lovu:1314@/salvation?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		panic("连接数据库出错")
	}
	defer db.Close()
	if !db.HasTable(entity.User{}) {
		log.Println("User表不存在，新建中......")
		db.CreateTable(entity.User{})
		log.Println("User创建成功")
	}
	if !db.HasTable(entity.Secret{}) {
		log.Println("Secret表不存在，新建中......")
		db.CreateTable(entity.Secret{})
		log.Println("Secret创建成功")
	}
	if !db.HasTable(entity.Salvation{}) {
		log.Println("Salvation表不存在，新建中......")
		db.CreateTable(entity.Salvation{})
		log.Println("Salvation创建成功")
	}
}

func main() {
	router := gin.Default()
	router.Use(Cors())
	router.Use(TLSHandler())
	group1 := router.Group("/user")
	{
		group1.POST("/login", func(context *gin.Context) {
			var user entity.User
			context.BindJSON(&user)
			var noP entity.User
			noP.Name = user.Name
			hasP := user.UserFind(noP)
			fmt.Println(noP.Password)
			fmt.Println(hasP.Password)
			if strings.EqualFold(hasP.Password, user.Password) {
				context.JSON(http.StatusOK, true)
			} else {
				context.JSON(http.StatusBadRequest, "账号或密码错误，请重新输入")
			}
			fmt.Println()
		})
		group1.POST("/join", func(context *gin.Context) {
			var user entity.User
			context.BindJSON(&user)
			user.Coin = 3

			db, err := gorm.Open("mysql", "lovu:1314@/salvation?charset=utf8&parseTime=True&loc=Local")
			if err != nil {
				panic("连接数据库出错")
			}
			defer db.Close()
			check := entity.User{}
			db.Where("name = ?", user.Name).First(&check)
			fmt.Println(check.ID)
			if check.ID == 0 {
				db.Create(&user)
				context.JSON(http.StatusOK, "欢迎您："+user.Name+" 加入我们")
			} else {
				context.JSON(http.StatusBadRequest, "很抱歉，这个昵称已经有人使用了")

			}
		})
	}
	group2 := router.Group("/secret")
	{
		group2.POST("/all", func(context *gin.Context) {
			db, err := gorm.Open("mysql", "lovu:1314@/salvation?charset=utf8&parseTime=True&loc=Local")
			if err != nil {
				panic("连接数据库出错")
				context.JSON(http.StatusBadRequest, "服务器出现了一些问题，请稍后重试")
			} else {
				var s [] entity.Secret
				db.Not("status", 0).Find(&s)
				context.JSON(http.StatusOK, s)
			}
		})
		group2.POST("/add", func(context *gin.Context) {
			var secret entity.Secret
			content := context.PostForm("content")
			user_id := context.PostForm("user_id")
			if user_id == "" {
				context.JSON(http.StatusBadRequest, "没有登录是不可以发小秘密的哦")
			} else {
				secret.Content = content
				id, _ := strconv.ParseUint(user_id, 10, 64)
				secret.UserID = uint(id)
				db, err := gorm.Open("mysql", "lovu:1314@/salvation?charset=utf8&parseTime=True&loc=Local")
				if err != nil {
					panic("连接数据库出错")
					context.JSON(http.StatusBadRequest, "服务器出现了一些问题，请稍后重试")
				} else {
					defer db.Close()
					fmt.Println(secret)
					db.Create(&secret)
					context.JSON(http.StatusOK, true)
				}
			}
		})
	}
	//_ = router.RunTLS(":1314", "src/ssl.pem", "src/ssl.key")
	router.Run(":1314")
}

//Cors 跨域处理
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next()
	}
}

// TLSHandler tls中间件
func TLSHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
