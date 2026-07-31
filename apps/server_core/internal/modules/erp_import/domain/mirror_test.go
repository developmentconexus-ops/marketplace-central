package domain

import (
	"reflect"
	"testing"
)

func TestMirrorProductUnknownSellableFieldsRemainNil(t *testing.T) {
	row := MirrorProduct{CodigoProduto: "P-1"}
	value := reflect.ValueOf(row)
	for _, fieldName := range []string{"Usoprod", "ADEcommerce"} {
		field := value.FieldByName(fieldName)
		if !field.IsValid() {
			t.Fatalf("MirrorProduct.%s is missing", fieldName)
		}
		if field.Kind() != reflect.Ptr || !field.IsNil() {
			t.Fatalf("MirrorProduct.%s = %#v, want nil *string", fieldName, field.Interface())
		}
	}
}
