package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/RodBarenco/colab-project-api/recoveruserpass"
	"github.com/RodBarenco/colab-project-api/rsakeys"
	"github.com/RodBarenco/colab-project-api/utils"
)

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var params recoveruserpass.RecoverParams

	// Realizar a validação e extração dos parâmetros do corpo da requisição
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Verificar se a senha é válida
	if !utils.IsValidPassword(params.Password) {
		RespondWithError(w, http.StatusBadRequest, "Invalid password format")
		return
	}

	// Extrair o token da query da URL
	tokenStr := r.URL.Query().Get("token")

	hexRandomPass, err := recoveruserpass.ValidateRecoveryToken(tokenStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Obter a chave privada para descriptografar
	privateKey, err := rsakeys.ReadPrivateKeyFromFile()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error reading private key")
		return
	}

	randomPassBytes, err := hex.DecodeString(hexRandomPass)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid random pass format in token")
		return
	}

	// Descriptografar a senha aleatória
	decryptedRandomPass, err := rsakeys.DecryptWithPrivateKey(privateKey, randomPassBytes)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error decrypting random pass")
		return
	}

	// Comparar as senhas aleatórias
	if string(decryptedRandomPass) != params.RandomPass {
		RespondWithError(w, http.StatusBadRequest, "Random pass mismatch")
		return
	}

	params.UserID, err = utils.GetUserIDByEmail(dbAccessor, params.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn`t found that user")
		return
	}

	// Atualizar a senha no banco de dados
	err = recoveruserpass.UpdatePassword(dbAccessor, params.UserID, params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := "Password updated successfully, your new pass is: " + params.Password

	RespondWithJSON(w, http.StatusOK, response)
}

func SendRecoveryEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the email address
	var request struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if !utils.IsValidEmail(request.Email) {
		RespondWithError(w, http.StatusBadRequest, "invalid email")
		return
	}

	// Call SendRecoveryEmail function
	message, err := recoveruserpass.SendRecoveryEmail(request.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": message})
}
