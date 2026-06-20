package token_test

import (
	"testing"
	"time"

	"proyecto-go/pkg/token"
)

const testSecret = "secreto-de-test-muy-largo"

func TestGenerate_Validate_Exitoso(t *testing.T) {
	tokenStr, err := token.Generate(42, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("error generando token: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("se esperaba un token no vacio")
	}

	claims, err := token.Validate(tokenStr, testSecret)
	if err != nil {
		t.Fatalf("error validando token: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("UserID esperado 42, obtenido %d", claims.UserID)
	}
}

func TestValidate_SecretoIncorrecto(t *testing.T) {
	tokenStr, _ := token.Generate(1, "secreto-correcto", time.Hour)

	_, err := token.Validate(tokenStr, "secreto-incorrecto")
	if err == nil {
		t.Error("se esperaba error con secreto incorrecto")
	}
}

func TestValidate_TokenExpirado(t *testing.T) {
	tokenStr, _ := token.Generate(1, testSecret, -time.Hour)

	_, err := token.Validate(tokenStr, testSecret)
	if err == nil {
		t.Error("se esperaba error con token ya expirado")
	}
}

func TestValidate_TokenMalformado(t *testing.T) {
	_, err := token.Validate("esto.no.es.un.jwt", testSecret)
	if err == nil {
		t.Error("se esperaba error con token malformado")
	}
}

func TestValidate_TokenVacio(t *testing.T) {
	_, err := token.Validate("", testSecret)
	if err == nil {
		t.Error("se esperaba error con token vacio")
	}
}
