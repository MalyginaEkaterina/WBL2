package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetAnagrams(t *testing.T) {
	res := GetAnagrams([]string{"Тяпа", "Пятка", "ППП", "Пятак", "слиток", "ппп", "Листок", "Стол", "тяпка", "СТОЛИК", "Ппп"})
	wanted := map[string][]string{"пятка": {"пятак", "пятка", "тяпка"}, "слиток": {"листок", "слиток", "столик"}}
	areEqual := reflect.DeepEqual(res, wanted)
	assert.Equal(t, true, areEqual)
}
