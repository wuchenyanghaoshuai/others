package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// User 表示登录用户
type User struct {
	ID       int
	Username string
	Password string
}

// Person 表示人员数据
type Person struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

// 数据库连接
var db *sql.DB



func main() {
	// 连接数据库
	var err error
	db, err = sql.Open("mysql", "root:rootpassword@tcp(mysql.default:3306)/people_db?charset=utf8mb4&collation=utf8mb4_unicode_ci")
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	r := gin.Default()

	// 设置HTML模板
	r.LoadHTMLGlob("templates/*")

	// 设置静态文件目录
	r.Static("/static", "./static")

	// 设置session
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 中间件：检查用户是否已登录
	authRequired := func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}

	// 登录页面
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "登录",
		})
	})

	// 根路径重定向到登录页面
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 处理登录请求
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 从数据库验证用户名和密码
		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
		if err == nil && storedPassword == password {
			session := sessions.Default(c)
			session.Set("user", username)
			session.Save()
			c.Redirect(http.StatusFound, "/home")
			return
		}

		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title":   "登录",
			"message": "用户名或密码错误",
		})
	})

	// 主页
	r.GET("/home", authRequired, func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		// 从数据库获取人员列表
		rows, err := db.Query("SELECT id, name, gender, age FROM people")
		if err != nil {
			c.HTML(http.StatusInternalServerError, "home.html", gin.H{
				"title":   "主页",
				"user":    user,
				"message": "获取数据失败: " + err.Error(),
			})
			return
		}
		defer rows.Close()

		var people []Person
		for rows.Next() {
			var p Person
			if err := rows.Scan(&p.ID, &p.Name, &p.Gender, &p.Age); err != nil {
				log.Println("扫描数据失败:", err)
				continue
			}
			people = append(people, p)
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"title":  "主页",
			"user":   user,
			"people": people,
		})
	})

	// 添加人员数据
	r.POST("/add", authRequired, func(c *gin.Context) {
		name := c.PostForm("name")
		gender := c.PostForm("gender")
		ageStr := c.PostForm("age")

		age, err := strconv.Atoi(ageStr)
		if err != nil {
			// 获取人员列表以便在错误页面显示
			rows, _ := db.Query("SELECT id, name, gender, age FROM people")
			var people []Person
			if rows != nil {
				defer rows.Close()
				for rows.Next() {
					var p Person
					if err := rows.Scan(&p.ID, &p.Name, &p.Gender, &p.Age); err == nil {
						people = append(people, p)
					}
				}
			}

			c.HTML(http.StatusBadRequest, "home.html", gin.H{
				"title":   "主页",
				"message": "年龄必须是数字",
				"people":  people,
			})
			return
		}

		// 将人员数据插入数据库
		_, err = db.Exec("INSERT INTO people (name, gender, age) VALUES (?, ?, ?)", name, gender, age)
		if err != nil {
			// 获取人员列表以便在错误页面显示
			rows, _ := db.Query("SELECT id, name, gender, age FROM people")
			var people []Person
			if rows != nil {
				defer rows.Close()
				for rows.Next() {
					var p Person
					if err := rows.Scan(&p.ID, &p.Name, &p.Gender, &p.Age); err == nil {
						people = append(people, p)
					}
				}
			}

			c.HTML(http.StatusInternalServerError, "home.html", gin.H{
				"title":   "主页",
				"message": "添加数据失败: " + err.Error(),
				"people":  people,
			})
			return
		}

		c.Redirect(http.StatusFound, "/home")
	})

	// 登出
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/login")
	})

	// 启动服务器
	r.Run(":8081")
}