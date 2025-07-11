package router

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/argon2"
)

const timing = 1
const memory = 64 * 1024
const threads = 4
const keylen = 32

var random_salt = os.Getenv("HASHING_SALT")

func SetAuthRouter(app *fiber.App, conn *pgx.Conn) {

	app.Post("/student/", func(c *fiber.Ctx) error {
		var signInDetails SignInDetails
		err := c.BodyParser(&signInDetails)
		signInDetails.Role = "student"
		if err != nil {
			c.Status(400)
			return c.JSON(buildErrorMResponse("malformed request"))
		}

		emailExists := getEmailExistence(c, conn, signInDetails.Email)

		if emailExists {
			c.Status(409)
			return c.JSON(buildErrorMResponse("email in use"))
		}

		passwordHash := base64.URLEncoding.EncodeToString(generateHash(signInDetails.Password))

		var userDetails UserDetails
		ConvertSignInDetailsToUserDetails(signInDetails, &userDetails)

		createUserString := buildUserInsertString(userDetails, passwordHash)

		_, commandError := conn.Exec(c.Context(), createUserString)

		if commandError != nil {
			fmt.Println(commandError)
			c.Status(500)
			return c.JSON(buildErrorMResponse("transaction failed"))
		}

		c.Status(200)
		return c.JSON(fiber.Map{
			"message":      "user created",
			"user_details": userDetails,
		})
	})

	app.Get("/student/", func(c *fiber.Ctx) error {
		var signInDetails SignInDetails
		err := c.BodyParser(&signInDetails)
		signInDetails.Role = "student"
		if err != nil {
			c.Status(400)
			return c.JSON(buildErrorMResponse("malformed request"))
		}

		emailExists := getEmailExistence(c, conn, signInDetails.Email)

		if !emailExists {
			c.Status(404)
			return c.JSON(buildErrorMResponse("email does not exist"))
		}

		passwordHash := fmt.Sprintf("%v", base64.URLEncoding.EncodeToString(generateHash(signInDetails.Password)))
		queryUserPassHashString := buildUserPassHashQueryString(signInDetails.Email)

		var storedHash string
		conn.QueryRow(c.Context(), queryUserPassHashString).Scan(&storedHash)

		if passwordHash != storedHash {
			c.Status(400)
			return c.JSON(buildErrorMResponse("incorrect password"))
		}

		queryUserString := buildUserQueryString(signInDetails.Email)

		var userDetailString string
		conn.QueryRow(c.Context(), queryUserString).Scan(&userDetailString)

		var userDetails UserDetails
		jsonErr := json.Unmarshal([]byte(userDetailString), &userDetails)
		if jsonErr != nil {
			c.Status(500)
			return c.JSON(buildErrorMResponse("server error"))
		}

		c.Status(200)
		return c.JSON(fiber.Map{
			"message": "user created",
			"token":   userDetails,
		})
	})

}

func getEmailExistence(c *fiber.Ctx, conn *pgx.Conn, email string) bool {
	emailQueryString := fmt.Sprintf(`SELECT email FROM USERS WHERE email = '%v'`, email)
	var foundEmail string
	conn.QueryRow(c.Context(), emailQueryString).Scan(&foundEmail)
	return email == foundEmail
}

func buildErrorMResponse(reason string) fiber.Map {
	return fiber.Map{
		"message": fmt.Sprintf("%v", reason),
	}
}

func buildUserInsertString(userDetails UserDetails, passwordHash string) string {
	id := userDetails.Id
	email := userDetails.Email
	role := userDetails.Role
	firstName := userDetails.First_name
	lastName := userDetails.Last_name
	avatarUrl := userDetails.Avatar_url
	phoneNumber := userDetails.Phone_number
	timzezone := userDetails.Timezone
	darkmode := userDetails.Dark_mode
	emailNotifications := userDetails.Email_notifications
	language := userDetails.Language
	isVerified := userDetails.Is_verified
	createUserString := fmt.Sprintf(`
		INSERT INTO USERS (id, email, password_hash, role, first_name, last_name, avatar_url, phone_number, timezone, dark_mode, email_notifications, language, is_verified)
		VALUES (%v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, %v, '%v', %v);`,
		id, email, passwordHash, role, firstName, lastName, avatarUrl, phoneNumber, timzezone, darkmode, emailNotifications, language, isVerified)
	return createUserString
}

func buildUserPassHashQueryString(email string) string {
	queryString := fmt.Sprintf(`
		SELECT password_hash FROM USERS WHERE email = '%v'`, email)
	return queryString
}

func buildUserQueryString(email string) string {
	queryString := fmt.Sprintf(`
		SELECT row_to_json(user_details) FROM (SELECT * FROM USERS WHERE email = '%v') AS user_details`, email)
	return queryString
}

func generateHash(password string) []byte {
	return argon2.IDKey([]byte(password), []byte(random_salt), timing, memory, threads, keylen)
}
