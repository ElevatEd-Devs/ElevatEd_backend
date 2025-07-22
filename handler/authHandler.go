package handler

import (
	"elevated_backend/functions"
	"elevated_backend/structs"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SignupHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var authFunc functions.AuthFunc
	var signInDetails structs.SignInDetails
	err := c.BodyParser(&signInDetails)
	if err != nil {
		c.Status(400)
		return c.JSON(authFunc.BuildErrorMResponse("malformed request"))
	}

	emailExists := authFunc.GetEmailExistence(c, conn, signInDetails.Email)

	if emailExists {
		c.Status(409)
		return c.JSON(authFunc.BuildErrorMResponse("email in use"))
	}

	passwordHash := base64.URLEncoding.EncodeToString(authFunc.GenerateHash(signInDetails.Password))

	var userDetails structs.UserDetails
	structs.ConvertSignInDetailsToUserDetails(signInDetails, &userDetails)

	createUserString := authFunc.BuildUserInsertString(userDetails, passwordHash)

	_, commandError := conn.Exec(c.Context(), createUserString)

	if commandError != nil {
		// fmt.Println(commandError)
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not create user"))
	}

	jwtToken, jwtGenErr := authFunc.GenerateJWT(c, userDetails)
	if jwtGenErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not create session"))
	}

	refreshToken := authFunc.GenerateRefreshToken()

	sessionErr := authFunc.CreateSession(c, conn, userDetails.Id, base64.StdEncoding.EncodeToString(authFunc.GenerateHash(refreshToken)), "")
	if sessionErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not create session"))
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "user created",
		"jwt":     jwtToken,
		"token":   refreshToken,
	})
}

func LoginHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var authFunc functions.AuthFunc
	var signInDetails structs.SignInDetails
	err := c.BodyParser(&signInDetails)
	if err != nil {
		c.Status(400)
		return c.JSON(authFunc.BuildErrorMResponse("malformed request"))
	}

	emailExists := authFunc.GetEmailExistence(c, conn, signInDetails.Email)

	if !emailExists {
		c.Status(404)
		return c.JSON(authFunc.BuildErrorMResponse("email does not exist"))
	}

	passwordHash := fmt.Sprintf("%v", base64.URLEncoding.EncodeToString(authFunc.GenerateHash(signInDetails.Password)))
	queryUserPassHashString := authFunc.BuildUserPassHashQueryString(signInDetails.Email)

	var storedHash string
	conn.QueryRow(c.Context(), queryUserPassHashString).Scan(&storedHash)

	if passwordHash != storedHash {
		c.Status(400)
		return c.JSON(authFunc.BuildErrorMResponse("incorrect password"))
	}

	queryUserString := authFunc.BuildUserQueryString(signInDetails.Email)

	var userDetailString string
	conn.QueryRow(c.Context(), queryUserString).Scan(&userDetailString)

	var userDetails structs.UserDetails
	jsonErr := json.Unmarshal([]byte(userDetailString), &userDetails)
	if jsonErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not get user"))
	}

	jwtToken, jwtGenErr := authFunc.GenerateJWT(c, userDetails)
	if jwtGenErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("JWT could not be generated"))
	}

	refreshToken := authFunc.GenerateRefreshToken()

	sessionErr := authFunc.CreateSession(c, conn, userDetails.Id, base64.StdEncoding.EncodeToString(authFunc.GenerateHash(refreshToken)), "")
	// fmt.Println(sessionErr)
	if sessionErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not create session"))
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "session created",
		"jwt":     jwtToken,
		"token":   refreshToken,
	})
}

func JWTHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var authFunc functions.AuthFunc
	var refreshToken functions.RefreshTokenReceptor
	err := c.BodyParser(&refreshToken)
	if err != nil {
		c.Status(400)
		return c.JSON(authFunc.BuildErrorMResponse("malformed request"))
	}

	refreshHash := base64.StdEncoding.EncodeToString(authFunc.GenerateHash(refreshToken.Token))
	validSession, userSession, sessionErr := authFunc.IsSessionValid(c, conn, refreshHash)

	if !validSession && sessionErr == nil {
		c.Status(401)
		return c.JSON(authFunc.BuildErrorMResponse("session expired"))
	}

	if sessionErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not verify session"))
	}

	endSessErr := authFunc.EndSession(c, conn, userSession, refreshHash)

	if endSessErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not end older session"))
	}

	queryUserString := authFunc.BuildUserQueryString(refreshToken.Email)

	var userDetailString string
	conn.QueryRow(c.Context(), queryUserString).Scan(&userDetailString)

	var userDetails structs.UserDetails
	jsonErr := json.Unmarshal([]byte(userDetailString), &userDetails)
	if jsonErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not refresh session"))
	}

	jwtToken, jwtGenErr := authFunc.GenerateJWT(c, userDetails)
	if jwtGenErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not refresh session"))
	}

	newRefreshToken := authFunc.GenerateRefreshToken()

	sessionErr = authFunc.CreateSession(c, conn, userDetails.Id, base64.StdEncoding.EncodeToString(authFunc.GenerateHash(newRefreshToken)), "")
	// fmt.Println(sessionErr)
	if sessionErr != nil {
		c.Status(500)
		return c.JSON(authFunc.BuildErrorMResponse("could not create session"))
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "session continued",
		"jwt":     jwtToken,
		"token":   newRefreshToken,
	})
}
