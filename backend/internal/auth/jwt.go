package auth

import (
    "fmt"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var secret string = os.Getenv("SECRET_KEY")
var secretKey = []byte(secret) // Change this! Keep it secret.

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "my-backend",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

func VerifyToken(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Validate the alg is what you expect
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}

// func main() {
//     token, err := GenerateToken(123)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Println("Token:", token)

//     claims, err := VerifyToken(token)
//     if err != nil {
//         fmt.Println("Invalid:", err)
//     } else {
//         fmt.Printf("Valid! UserID: %d\n", claims.UserID)
//     }
// }