package main

import (
    "os"
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

func showError(c *gin.Context, httpCode int, err error) {
    c.HTML(httpCode, "error.tmpl", gin.H{
        "error": httpCode,
        "message": err.Error(),
    })
}

func buildPasteURL(c *gin.Context, id string) string {
    hostname := os.Getenv("PUBLIC_HOSTNAME")
    if (hostname == "") {
        hostname = fmt.Sprintf("%s", location.Get(c))
    }
    return fmt.Sprintf("%s/paste/%s", hostname, id)
}

func addRoutes(r *gin.Engine) {
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{})
    })

    r.GET("/paste/:id", func(c *gin.Context) {
        id := c.Param("id")
        paste, err := getPaste(id)
        if err != nil {
            showError(c, http.StatusNotFound, fmt.Errorf("Paste %s not found", id))
        }
        c.HTML(http.StatusOK, "retrieve.tmpl", gin.H{
            "id":    id,
            "paste": paste,
        })
    })

    r.GET("/paste/:id/raw", func(c *gin.Context) {
        id := c.Param("id")
        paste, err := getPaste(id)
        if err != nil {
            showError(c, http.StatusNotFound, fmt.Errorf("Paste %s not found", id))
        }
        c.String(http.StatusOK, paste)
    })

    r.POST("/api", func(c *gin.Context) {
        paste, err := ioutil.ReadAll(c.Request.Body)
        if err != nil {
            showError(c, http.StatusInternalServerError, err)
        }
        id, err := savePasteInFile(string(paste))
        if err != nil {
            showError(c, http.StatusInternalServerError, err)
        }
        c.String(http.StatusOK, buildPasteURL(c, id) + "\n")
    })

    r.POST("/", func(c *gin.Context) {
        paste := c.PostForm("code")
        id, err := savePasteInFile(paste)
        if err != nil {
            showError(c, http.StatusInternalServerError, err)
        }
        c.Redirect(301, buildPasteURL(c, id))
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
