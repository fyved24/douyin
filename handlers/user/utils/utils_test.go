package utils

import (
	"testing"
)

func TestStringToUint(t *testing.T) {
	var cases = []struct {
		s   string
		ans uint
	}{
		{"123", uint(123)},
		{"1234", uint(1234)},
	}

	for _, testcase := range cases {
		if result := StringToUint(testcase.s); result != testcase.ans {
			t.Fatalf("string: %v to uint, except: %v, actual: %v", testcase.s, testcase.ans, result)
		} else {
			t.Logf("string: %v to uint, except: %v, pass", testcase.s, testcase.ans)
		}
	}
}

func TestGetIsLoginFromToken(t *testing.T) {
	var cases = []struct {
		token  string
		except bool
	}{
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imppbmh1aSIsInBhc3N3b3JkIjoiNDVjN2ZjZTNmM2Y4OTMxNzg1NGYwOWNmN2YyZjQ4YzYiLCJ1c2VyX2lkIjoiMjgiLCJpc19sb2dpbiI6dHJ1ZSwiZXhwIjoxNjgzODIwNTI3LCJpc3MiOiJkb3V5aW4iLCJuYmYiOjE2NzY2MjA0Njd9.QMquDFbh59YUZgC0l9Kq-Mmt8aJSMaI-oQzL2BTl3IY", true},
	}

	for _, testcase := range cases {
		if result := GetIsLoginFromToken(testcase.token); result != testcase.except {
			t.Fatalf("except: %v, actual: %v", testcase.except, result)
		} else {
			t.Logf("except: %v, pass", testcase.except)
		}
	}
}

func TestGetPasswordFromToken(t *testing.T) {
	var cases = []struct {
		token  string
		except string
	}{
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imppbmh1aSIsInBhc3N3b3JkIjoiNDVjN2ZjZTNmM2Y4OTMxNzg1NGYwOWNmN2YyZjQ4YzYiLCJ1c2VyX2lkIjoiMjgiLCJpc19sb2dpbiI6dHJ1ZSwiZXhwIjoxNjgzODIwNTI3LCJpc3MiOiJkb3V5aW4iLCJuYmYiOjE2NzY2MjA0Njd9.QMquDFbh59YUZgC0l9Kq-Mmt8aJSMaI-oQzL2BTl3IY", "45c7fce3f3f89317854f09cf7f2f48c6"},
	}

	for _, testcase := range cases {
		if result := GetPasswordFromToken(testcase.token); result != testcase.except {
			t.Fatalf("except: %v, actual: %v", testcase.except, result)
		} else {
			t.Logf("except: %v, pass", testcase.except)
		}
	}
}

func TestGetUsernameFromToken(t *testing.T) {
	var cases = []struct {
		token  string
		except string
	}{
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjEyMzQ1Njc4IiwicGFzc3dvcmQiOiIyNWQ1NWFkMjgzYWE0MDBhZjQ2NGM3NmQ3MTNjMDdhZCIsInVzZXJfaWQiOiIxOSIsImlzX2xvZ2luIjp0cnVlLCJleHAiOjE2ODMxMzY3NTIsImlzcyI6ImRvdXlpbiIsIm5iZiI6MTY3NTkzNjY5Mn0.--zA9-ggH66Rjz5N8xDAIzWPiXE1H0dYzOZN6l716fE", "12345678"},
	}

	for _, testcase := range cases {
		if result := GetUsernameFromToken(testcase.token); result != testcase.except {
			t.Fatalf("except: %v, actual: %v", testcase.except, result)
		} else {
			t.Logf("except: %v, pass", testcase.except)
		}
	}
}

func TestGetUserIDFromToken(t *testing.T) {
	var cases = []struct {
		token  string
		except uint
	}{
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjEyMzQ1Njc4IiwicGFzc3dvcmQiOiIyNWQ1NWFkMjgzYWE0MDBhZjQ2NGM3NmQ3MTNjMDdhZCIsInVzZXJfaWQiOiIxOSIsImlzX2xvZ2luIjp0cnVlLCJleHAiOjE2ODMxMzY3NTIsImlzcyI6ImRvdXlpbiIsIm5iZiI6MTY3NTkzNjY5Mn0.--zA9-ggH66Rjz5N8xDAIzWPiXE1H0dYzOZN6l716fE", uint(19)},
	}

	for _, testcase := range cases {
		if result := GetUserIDFromToken(testcase.token); result != testcase.except {
			t.Fatalf("except: %v, actual: %v", testcase.except, result)
		} else {
			t.Logf("except: %v, pass", testcase.except)
		}
	}
}

func TestGetUserToken(t *testing.T) {
	var cases = []struct {
		username       string
		password       string
		userID         uint
		islogin        bool
		exceptUsername string
		exceptPassword string
		exceptUserID   uint
		exceptIslogin  bool
	}{
		{username: "12344567890", password: "123456789", userID: uint(5), islogin: true,
			exceptUsername: "12344567890", exceptPassword: "123456789", exceptUserID: uint(5), exceptIslogin: true},
	}

	for _, testcase := range cases {
		token := GetUserToken(testcase.username, testcase.password, testcase.userID, testcase.islogin)
		actualUsername := GetUsernameFromToken(token)
		actualPassword := GetPasswordFromToken(token)
		actualUserID := GetUserIDFromToken(token)
		actualIsLogin := GetIsLoginFromToken(token)
		if actualIsLogin == testcase.exceptIslogin && actualUserID == testcase.exceptUserID &&
			actualPassword == testcase.exceptPassword && actualUsername == testcase.exceptUsername {
			t.Logf("pass")
		} else {
			t.Fatalf("fail")
		}
	}
}

func TestEncodePassword(t *testing.T) {
	var cases = []struct {
		password string
		except   string
	}{
		{"jinhui123", "45c7fce3f3f89317854f09cf7f2f48c6"},
		{"45c7fce3f3f89317854f09cf7f2f48c6", "0d21123d9a7c69bc1c91d7c708dca647"},
		{"0d21123d9a7c69bc1c91d7c708dca647", "20bc8c2d80f9a3d90189972b2eabdd4b"},
	}

	for _, testcase := range cases {
		if result := EncodePassword(testcase.password); result != testcase.except {
			t.Fatalf("except: %v, actual: %v", testcase.except, result)
		} else {
			t.Logf("except: %v, pass", testcase.except)
		}
	}
}
