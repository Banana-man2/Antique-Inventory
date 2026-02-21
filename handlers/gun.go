package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"Antique-Inventory/ent"
	"Antique-Inventory/ent/gun"
)

type GunHandler struct {
	Client *ent.Client
}

// HTML handlers

func (h *GunHandler) ListGuns(c *gin.Context) {
	ctx := c.Request.Context()
	q := strings.TrimSpace(c.Query("q"))

	query := h.Client.Gun.Query().Order(ent.Asc(gun.FieldID))
	if q != "" {
		query = query.Where(gun.Or(
			gun.GunNameContainsFold(q),
			gun.SerialNumberContainsFold(q),
			gun.DescriptionContainsFold(q),
			gun.MiscAttachmentsContainsFold(q),
		))
	}

	guns, err := query.All(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching guns: %v", err)
		return
	}

	var totalValue float64
	for _, g := range guns {
		if g.Value != nil {
			totalValue += *g.Value
		}
	}

	c.HTML(http.StatusOK, "guns_list.html", gin.H{
		"guns":       guns,
		"count":      len(guns),
		"totalValue": totalValue,
		"query":      q,
	})
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

func (h *GunHandler) ServeImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	g, err := h.Client.Gun.Get(c.Request.Context(), id)
	if err != nil || g.Image == nil {
		c.String(http.StatusNotFound, "Image not found")
		return
	}
	contentType := http.DetectContentType(*g.Image)
	c.Data(http.StatusOK, contentType, *g.Image)
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
	if v := c.PostForm("serial_number"); v != "" {
		create.SetSerialNumber(v)
	}
	if v := c.PostForm("description"); v != "" {
		create.SetDescription(v)
	}
	if v := c.PostForm("misc_attachments"); v != "" {
		create.SetMiscAttachments(v)
	}
	if v := c.PostForm("value"); v != "" {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			create.SetValue(val)
		}
	}
	if file, err := c.FormFile("image"); err == nil {
		f, err := file.Open()
		if err == nil {
			defer f.Close()
			data, err := io.ReadAll(f)
			if err == nil {
				create.SetImage(data)
			}
		}
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
	if v := c.PostForm("serial_number"); v != "" {
		update.SetSerialNumber(v)
	} else {
		update.ClearSerialNumber()
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
	if v := c.PostForm("value"); v != "" {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			update.SetValue(val)
		}
	} else {
		update.ClearValue()
	}
	if file, err := c.FormFile("image"); err == nil {
		f, err := file.Open()
		if err == nil {
			defer f.Close()
			data, err := io.ReadAll(f)
			if err == nil {
				update.SetImage(data)
			}
		}
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
		GunName         string   `json:"gun_name"`
		Year            *int     `json:"year"`
		Condition       *int     `json:"condition"`
		SerialNumber    string   `json:"serial_number"`
		Description     string   `json:"description"`
		MiscAttachments string   `json:"misc_attachments"`
		Value           *float64 `json:"value"`
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
	if input.SerialNumber != "" {
		create.SetSerialNumber(input.SerialNumber)
	}
	if input.Description != "" {
		create.SetDescription(input.Description)
	}
	if input.MiscAttachments != "" {
		create.SetMiscAttachments(input.MiscAttachments)
	}
	if input.Value != nil {
		create.SetValue(*input.Value)
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
		GunName         string   `json:"gun_name"`
		Year            *int     `json:"year"`
		Condition       *int     `json:"condition"`
		SerialNumber    string   `json:"serial_number"`
		Description     string   `json:"description"`
		MiscAttachments string   `json:"misc_attachments"`
		Value           *float64 `json:"value"`
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
	if input.SerialNumber != "" {
		update.SetSerialNumber(input.SerialNumber)
	} else {
		update.ClearSerialNumber()
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
	if input.Value != nil {
		update.SetValue(*input.Value)
	} else {
		update.ClearValue()
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

func (h *GunHandler) DeleteGunForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	err = h.Client.Gun.DeleteOneID(id).Exec(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/guns")
}
