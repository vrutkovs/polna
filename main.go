package main

import (
    "fmt"
    "io/ioutil"
    "net/http"

    "html/template"
    "path/filepath"

    "github.com/gin-contrib/location"
    "github.com/gin-contrib/multitemplate"
    "github.com/gin-gonic/gin"
)

const uploadDir = "upload/"

func loadTemplates(templatesDir string) multitemplate.Render {
    r := multitemplate.New()

    base, err := filepath.Abs(templatesDir + "/base.tmpl")
    if err != nil {
        panic(err.Error())
    }

    layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
    if err != nil {
        panic(err.Error())
    }

    includes, err := filepath.Glob(templatesDir + "/includes/*.tmpl")
    if err != nil {
        panic(err.Error())
    }

    // Generate our templates map from our layouts/ and includes/ directories
    for _, layout := range layouts {
        baseAndLayouts := append([]string{base}, layout)
        files := append(baseAndLayouts, includes...)
        r.Add(filepath.Base(layout), template.Must(template.ParseFiles(files...)))
    }
    return r
}

func addRoutes(r *gin.Engine) {
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{})
    })

    r.GET("/paste/:id", func(c *gin.Context) {
        id := c.Param("id")
        paste, err := getPaste(id)
        if err != nil {
            c.AbortWithError(http.StatusNotFound, err)
        }
        c.HTML(http.StatusOK, "retrieve.tmpl", gin.H{
            "id":    id,
            "paste": paste,
        })
    })

    r.GET("/paste/:id/raw", func(c *gin.Context) {
        paste, err := getPaste(c.Param("id"))
        if err != nil {
            c.AbortWithError(http.StatusNotFound, err)
        }
        c.String(http.StatusOK, paste)
    })

    r.POST("/api", func(c *gin.Context) {
        paste, err := ioutil.ReadAll(c.Request.Body)
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
        }
        id, err := savePasteInFile(string(paste))
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
        }
        url := fmt.Sprintf("%s/paste/%s\n", location.Get(c), id)
        c.String(http.StatusOK, url)
    })

    r.POST("/", func(c *gin.Context) {
        paste := c.PostForm("code")
        id, err := savePasteInFile(paste)
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
        }
        url := fmt.Sprintf("%s/paste/%s", location.Get(c), id)
        c.Redirect(301, url)
    })
}

func main() {
    r := gin.New()

    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(location.Default())

    r.Static("/static", "./static")

    r.HTMLRender = loadTemplates("./templates")
    addRoutes(r)
    r.Run()
}
