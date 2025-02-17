// Copyright 2022 The ILLA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package smtp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/golang-jwt/jwt/v4"
	email_cloud "github.com/illa-family/builder-backend/pkg/email-cloud"
)

var (
	USAGE_SIGNUP = "signup"
	USAGE_FORGET = "forgetpwd"
)

type Config struct {
	Username string `env:"ILLA_MAIL_USERNAME" envDefault:"m17602200056@163.com"`
	Password string `env:"ILLA_MAIL_PASSWORD" envDefault:"ESXNALKGBIAZCSYO"`
	Host     string `env:"ILLA_MAIL_HOST" envDefault:"smtp.163.com"`
	Port     string `env:"ILLA_MAIL_PORT" envDefault:"465"`
	Secret   string `env:"ILLA_SECRET_KEY" envDefault:"ausNV5NJfVCrz3tPXtW2ZGGCpUuWFVQbikZ6d7FyOfpw9RcyLiNpqx4pJ6fSX9JXhMfmIupKKjQElURR"`
}

type SMTPServer struct {
	From     string
	Password string
	Host     string
	Port     string
	Secret   string
}

type VCodeClaims struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	Usage string `json:"usage"`
	jwt.RegisteredClaims
}

func GetConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}

func NewSMTPServer(cfg *Config) SMTPServer {
	return SMTPServer{
		From:     cfg.Username,
		Password: cfg.Password,
		Host:     cfg.Host,
		Port:     cfg.Port,
		Secret:   cfg.Secret,
	}
}

func (s *SMTPServer) NewVerificationCode(email, usage string) (string, error) {
	// Authentication.
	//auth := smtp.PlainAuth("", s.From, s.Password, s.Host)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	//header := make(map[string]string)
	//header["From"] = "illa-builder" + "<" + s.From + ">"
	//header["To"] = email
	//header["Subject"] = "[ILLA]: Verification Code"
	//body := "Verification Code: " + vCode
	//message := ""
	//for k, v := range header {
	//	message += fmt.Sprintf("%s:%s\r\n", k, v)
	//}
	//message += "\r\n" + body
	//
	//// Sending email.
	//err := SendMailUsingTLS(s.Host+":"+s.Port, auth, s.From, email, []byte(message))
	//if err != nil {
	//	return "", err
	//}

	if err := email_cloud.SendVerificationEmail(email, vCode, usage); err != nil {
		return "", err
	}

	claims := &VCodeClaims{
		Email: email,
		Code:  vCode,
		Usage: usage,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "ILLA",
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute * 15),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	codeToken, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", err
	}

	return codeToken, nil
}

func (s *SMTPServer) SendSubscriptionEmail(email string) error {
	return email_cloud.SendSubscriptionEmail(email)
}

func (s *SMTPServer) ValidateVerificationCode(codeToken, vCode, email, usage string) (bool, error) {
	vCodeClaims := &VCodeClaims{}
	token, err := jwt.ParseWithClaims(codeToken, vCodeClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Secret), nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*VCodeClaims)
	if !(ok && claims.Usage == usage) {
		return false, errors.New("invalid verification token")
	}
	if !(claims.Code == vCode && claims.Email == email) {
		return false, errors.New("verification code wrong")
	}
	return true, nil
}

func SendMailUsingTLS(addr string, auth smtp.Auth, from string, to string, msg []byte) (err error) {
	c, err := Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	tos := strings.Split(to, ";")
	for _, addr := range tos {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
