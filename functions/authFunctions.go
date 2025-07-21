package functions

import (
	"elevated_backend/structs"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/argon2"
)

type AuthFunc struct{}

type CustomClaimStruct struct {
	Details              structs.UserDetails `json:"details"`
	jwt.RegisteredClaims `json:"registered_claims"`
}

type RefreshTokenReceptor struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func (*AuthFunc) GetEmailExistence(c *fiber.Ctx, conn *pgx.Conn, email string) bool {
	emailQueryString := fmt.Sprintf(`SELECT email FROM USERS WHERE email = '%v'`, email)
	var foundEmail string
	conn.QueryRow(c.Context(), emailQueryString).Scan(&foundEmail)
	return email == foundEmail
}

func (*AuthFunc) BuildErrorMResponse(reason string) fiber.Map {
	return fiber.Map{
		"status":  "failed",
		"message": fmt.Sprintf("%v", reason),
	}
}

func (*AuthFunc) BuildUserInsertString(userDetails structs.UserDetails, passwordHash string) string {
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

func (*AuthFunc) BuildUserPassHashQueryString(email string) string {
	queryString := fmt.Sprintf(`
		SELECT password_hash FROM USERS WHERE email = '%v'`, email)
	return queryString
}

func (*AuthFunc) BuildUserQueryString(email string) string {
	queryString := fmt.Sprintf(`
		SELECT row_to_json(user_details) FROM (SELECT * FROM USERS WHERE email = '%v') AS user_details`, email)
	return queryString
}

func (*AuthFunc) GenerateHash(password string) []byte {
	const timing = 1
	const memory = 64 * 1024
	const threads = 4
	const keylen = 32
	var random_salt = os.Getenv("HASHING_SALT")

	return argon2.IDKey([]byte(password), []byte(random_salt), timing, memory, threads, keylen)
}

func (*AuthFunc) GenerateJWT(c *fiber.Ctx, userDetails structs.UserDetails) (string, error) {
	var claimStruct CustomClaimStruct
	claimStruct.Details = userDetails
	claimStruct.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 3))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimStruct)

	secretkey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretkey))

	if err != nil {
		return "", err
	}
	return tokenString, err
}

func (*AuthFunc) ExtractJWTFromHeader(c *fiber.Ctx) string {
	rawBearer := c.Get("Authorization")
	tokenString := ""
	bearerPattern := regexp.MustCompile(`^Bearer\s[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`)
	if bearerPattern.MatchString(rawBearer) {
		tokenString = rawBearer[7:]
	}
	return tokenString
}

func (*AuthFunc) VerifyJWT(tokenString string) (bool, CustomClaimStruct, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaimStruct{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return false, CustomClaimStruct{}, err
	}

	if claims, ok := token.Claims.(*CustomClaimStruct); ok && token.Valid {
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return false, CustomClaimStruct{}, nil
		}
		return true, *claims, nil
	} else {
		return false, CustomClaimStruct{}, nil
	}
}

func (*AuthFunc) GenerateRefreshToken() string {
	bucket := []byte("abcdefghijklmnopqrstuvABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	var rv [32]byte

	for i := 0; i < 32; i++ {
		rv[i] = bucket[rand.Intn(len(bucket))]
	}
	return string(rv[:])
}

func (*AuthFunc) CreateSession(c *fiber.Ctx, conn *pgx.Conn, id int, refreshHash string, ip string) error {
	sessionCreationStamp := time.Now()
	sessionExpireStamp := sessionCreationStamp.Add(time.Hour * 24) // refresh token lasts 24 hours

	writeString := fmt.Sprintf(`
		INSERT INTO USER_SESSIONS (id, user_id, session_token, ip_address, user_agent, expires_at, created_at)
		VALUES (%v, %v, '%v', '%v', '%v', '%v', '%v')`, rand.Intn(10000), id, refreshHash, ip, "agent", sessionExpireStamp.UTC().Format("2006-01-02 15:04:05.000000"), sessionCreationStamp.UTC().Format("2006-01-02 15:04:05.000000"))

	_, err := conn.Exec(c.Context(), writeString)

	if err != nil {
		return err
	}

	return nil
}

func (*AuthFunc) IsSessionValid(c *fiber.Ctx, conn *pgx.Conn, refreshHash string) (bool, *structs.UserSessionStruct, error) {
	writeString := fmt.Sprintf(`
		SELECT * FROM  USER_SESSIONS WHERE session_token = '%v'`, refreshHash)

	var storedHash string
	var userSession structs.UserSessionStruct

	err := conn.QueryRow(c.Context(), writeString).Scan(
		&userSession.Id,
		&userSession.UserId,
		&storedHash,
		&userSession.Ip,
		&userSession.UserAgent,
		&userSession.ExpiresAt,
		&userSession.CreatedAt,
	)

	if err != nil {
		return false, nil, err
	}

	if refreshHash != storedHash || userSession.ExpiresAt.Before(time.Now().UTC()) {
		return false, nil, nil
	}

	return true, &userSession, nil
}

func (*AuthFunc) EndSession(c *fiber.Ctx, conn *pgx.Conn, userSession *structs.UserSessionStruct, refreshHash string) error {
	userSession.ExpiresAt = time.Now().UTC()

	writeString := fmt.Sprintf(`UPDATE USER_SESSIONS SET expires_at='%v' WHERE session_token='%v';`,
		userSession.ExpiresAt.Format("2006-01-02 15:04:05.000000"), refreshHash)

	_, err := conn.Exec(c.Context(), writeString)

	if err != nil {
		return err
	}

	return nil
}
