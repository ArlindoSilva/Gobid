package validator

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"
)

//2º
type Validator interface {
	Valid(ctx context.Context) Evaluator
}

type Evaluator map[string]string

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//4º
// AddFieldError adds an error for a specific field.
func (e *Evaluator) AddFieldError(key, message string) {
	if *e == nil {
		*e = make(map[string]string)
	}
	//Se já foi inicializado, se não é um ponteiro nulo, então adiciona o erro no mapa de erros.
	//sintaxe para acessar maps de um ponteiro, (*e)[key] acessa o valor associado à chave 'key' no mapa apontado por 'e'.
	if _, exists := (*e)[key]; !exists { 
		(*e)[key] = message
	}
}

//CheckField checa a condição e adiciona um erro para um campo específico se a condição for falsa. É apenas um wrapper para AddFieldError, que adiciona um erro de campo se a condição for falsa.
func (e *Evaluator) CheckField(ok bool, key, message string) {
	if !ok {
		e.AddFieldError(key, message)
	}
}

//NotBlank verifica se uma string não está vazia ou composta apenas por espaços em branco. 
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

//MaxChars verifica se uma string tem no máximo n caracteres(menor ou igual). A função utf8.RuneCountInString é usada para contar o número de runas (caracteres Unicode) na string, garantindo que a contagem seja precisa mesmo para caracteres multibyte. Pois o Length pode não trazer o tamanho correto.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

//MinChars verifica se uma string tem no mínimo n caracteres(maior ou igual)
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

//Matches verifica se uma string corresponde a uma expressão regular fornecida. A função rx.MatchString(value) retorna true se a string value corresponder ao padrão definido pela expressão regular rx, e false caso contrário.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}