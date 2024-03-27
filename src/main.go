// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

var token = os.Getenv("GITLAB_TOKEN")

func main() {
	InitLogger()
	r := gin.New()

	if os.Getenv("GIN_LOG") == "PRETTY" {
		r.Use(gin.Logger())
	} else {
		// Enable slog logging for gin framework
		// https://github.com/samber/slog-gin
		r.Use(sloggin.New(slog.Default()))
	}

	r.Use(gin.Recovery())

	r.GET("/", artifactFile)
	r.GET("/json", artifactJson)
	r.GET("/health", health)
	r.Run()
}
func health(c *gin.Context) {
	c.Status(http.StatusOK)
}
func artifactJson(c *gin.Context) {
	req, err := http.NewRequest(http.MethodGet, "https://gitlab.com/api/v4/projects/54944672/jobs/artifacts/main/raw/report.json?job=harvest", nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	req.Header = http.Header{
		"PRIVATE-TOKEN": {token},
		"Accept":        {"application/json"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var data map[string]interface{}
	json.Unmarshal(bodyBytes, &data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, data)
}

func artifactFile(c *gin.Context) {
	req, err := http.NewRequest(http.MethodGet, "https://gitlab.com/api/v4/projects/54944672/jobs/artifacts/main/raw/report.json?job=harvest", nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	req.Header = http.Header{
		"PRIVATE-TOKEN": {token},
		"Accept":        {"application/json"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var data map[string]interface{}
	json.Unmarshal(bodyBytes, &data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	file, _ := os.Create("report.json")
	defer file.Close()
	as_json, _ := json.MarshalIndent(data, "", "\t")
	file.Write(as_json)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=report.json")
	c.Header("Content-Type", "application/octet-stream")

	c.File("report.json")
}
