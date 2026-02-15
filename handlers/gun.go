package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Antique-Inventory/ent"
	"Antique-Inventory/ent/gun"
)

type GunHandler struct {
	Client *ent.Client
}

// HTML handlers

func (h *GunHandler) ListGuns(c *gin.Context) {
	guns, err := h.Client.Gun.Query().Order(ent.Asc(gun.FieldID)).All(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching guns: %v", err)
		return
	}
	c.HTML(http.StatusOK, "guns_list.html", gin.H{"guns": guns})
}

func (h *GunHandler) GetGun(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	g, err := h.Client.Gun.Get(c.Request.Context(), id)
	if err != nil {
		c.String(http.StatusNotFound, "Gun not found")
		return
	}
	c.HTML(http.StatusOK, "gun_detail.html", gin.H{"gun": g})
}

func (h *GunHandler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "gun_form.html", gin.H{"gun": nil})
}

func (h *GunHandler) ShowEditForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	g, err := h.Client.Gun.Get(c.Request.Context(), id)
	if err != nil {
		c.String(http.StatusNotFound, "Gun not found")
		return
	}
	c.HTML(http.StatusOK, "gun_form.html", gin.H{"gun": g})
}

func (h *GunHandler) CreateGunForm(c *gin.Context) {
	create := h.Client.Gun.Create().
		SetGunName(c.PostForm("gun_name"))

	if v := c.PostForm("year"); v != "" {
		if year, err := strconv.Atoi(v); err == nil {
			create.SetYear(year)
		}
	}
	if v := c.PostForm("condition"); v != "" {
		if cond, err := strconv.Atoi(v); err == nil {
			create.SetCondition(cond)
		}
	}
	if v := c.PostForm("description"); v != "" {
		create.SetDescription(v)
	}
	if v := c.PostForm("misc_attachments"); v != "" {
		create.SetMiscAttachments(v)
	}

	_, err := create.Save(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating gun: %v", err)
		return
	}
	c.Redirect(http.StatusFound, "/guns")
}

func (h *GunHandler) UpdateGunForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}

	update := h.Client.Gun.UpdateOneID(id).
		SetGunName(c.PostForm("gun_name"))

	if v := c.PostForm("year"); v != "" {
		if year, err := strconv.Atoi(v); err == nil {
			update.SetYear(year)
		}
	} else {
		update.ClearYear()
	}
	if v := c.PostForm("condition"); v != "" {
		if cond, err := strconv.Atoi(v); err == nil {
			update.SetCondition(cond)
		}
	} else {
		update.ClearCondition()
	}
	if v := c.PostForm("description"); v != "" {
		update.SetDescription(v)
	} else {
		update.ClearDescription()
	}
	if v := c.PostForm("misc_attachments"); v != "" {
		update.SetMiscAttachments(v)
	} else {
		update.ClearMiscAttachments()
	}

	_, err = update.Save(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error updating gun: %v", err)
		return
	}
	c.Redirect(http.StatusFound, "/guns")
}

// JSON API handlers

func (h *GunHandler) ListGunsJSON(c *gin.Context) {
	guns, err := h.Client.Gun.Query().Order(ent.Asc(gun.FieldID)).All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, guns)
}

func (h *GunHandler) GetGunJSON(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	g, err := h.Client.Gun.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gun not found"})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *GunHandler) CreateGun(c *gin.Context) {
	var input struct {
		GunName         string `json:"gun_name"`
		Year            *int   `json:"year"`
		Condition       *int   `json:"condition"`
		Description     string `json:"description"`
		MiscAttachments string `json:"misc_attachments"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	create := h.Client.Gun.Create().SetGunName(input.GunName)
	if input.Year != nil {
		create.SetYear(*input.Year)
	}
	if input.Condition != nil {
		create.SetCondition(*input.Condition)
	}
	if input.Description != "" {
		create.SetDescription(input.Description)
	}
	if input.MiscAttachments != "" {
		create.SetMiscAttachments(input.MiscAttachments)
	}

	g, err := create.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *GunHandler) UpdateGun(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input struct {
		GunName         string `json:"gun_name"`
		Year            *int   `json:"year"`
		Condition       *int   `json:"condition"`
		Description     string `json:"description"`
		MiscAttachments string `json:"misc_attachments"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := h.Client.Gun.UpdateOneID(id).SetGunName(input.GunName)
	if input.Year != nil {
		update.SetYear(*input.Year)
	} else {
		update.ClearYear()
	}
	if input.Condition != nil {
		update.SetCondition(*input.Condition)
	} else {
		update.ClearCondition()
	}
	if input.Description != "" {
		update.SetDescription(input.Description)
	} else {
		update.ClearDescription()
	}
	if input.MiscAttachments != "" {
		update.SetMiscAttachments(input.MiscAttachments)
	} else {
		update.ClearMiscAttachments()
	}

	g, err := update.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *GunHandler) DeleteGun(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err = h.Client.Gun.DeleteOneID(id).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
