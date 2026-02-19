package main

import (
	"context"
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	"Antique-Inventory/ent"
	"Antique-Inventory/handlers"
)

func main() {
	db, err := sql.Open("sqlite", "file:C:\\Coding\\antique_inventory.db?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer db.Close()

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	h := &handlers.GunHandler{Client: client}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// HTML routes (order matters: /guns/new before /guns/:id)
	r.GET("/guns", h.ListGuns)
	r.GET("/guns/new", h.ShowCreateForm)
	r.GET("/guns/:id", h.GetGun)
	r.GET("/guns/:id/edit", h.ShowEditForm)
	r.GET("/guns/:id/image", h.ServeImage)
	r.POST("/guns", h.CreateGunForm)
	r.POST("/guns/:id", h.UpdateGunForm)
	r.POST("/guns/:id/delete", h.DeleteGunForm)

	// JSON API routes
	r.GET("/api/guns", h.ListGunsJSON)
	r.GET("/api/guns/:id", h.GetGunJSON)
	r.POST("/api/guns", h.CreateGun)
	r.PUT("/api/guns/:id", h.UpdateGun)
	r.DELETE("/api/guns/:id", h.DeleteGun)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
