package main

import (
	"auth/models"
	"auth/router"
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed static
var static embed.FS

//go:embed templates
var tmpl embed.FS

func main() {
	err := models.StartMySQL()
	if err != nil {
		fmt.Println("数据库启动失败,%v", err)
	}
	defer models.DB.Close()
	r := router.StartRouter()

	// r.Static("/static", "./static")
	r.StaticFS("/static", http.FS(static))
	// r.LoadHTMLGlob("templates/*")
	t, _ := template.ParseFS(tmpl, "templates/*.html")
	r.SetHTMLTemplate(t)

	if err = r.Run(":8080"); err != nil {
		fmt.Println("运行时出错", err)
	}
}
