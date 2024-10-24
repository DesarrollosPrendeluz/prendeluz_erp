package controllers

import (

	//"prendeluz/erp/internal/dtos"

	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/userrepo"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	// Obtener los parámetros de la solicitud; se asume que se envían en el cuerpo de la solicitud
	var loginReq LoginRequest

	// Parsear el JSON del cuerpo de la solicitud a la estructura loginReq
	if err := c.BindJSON(&loginReq); err != nil {
		// Si hay un error en el parseo, devolver un error 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Crear una instancia del repositorio
	repo := *userrepo.NewUsersRepository(db.DB)

	// Verificar las credenciales del usuario
	token, id, err := repo.CheckCredentials(loginReq.Email, loginReq.Password)
	if err != nil {
		// Manejo del error si las credenciales no son correctas
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Responder con el token si las credenciales son correctas
	c.JSON(http.StatusOK, gin.H{"token": token, "id": id})
}
