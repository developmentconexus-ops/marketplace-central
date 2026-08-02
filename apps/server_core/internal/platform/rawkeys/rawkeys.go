// Package rawkeys detects provider payload keys that no DTO declares.
//
// It exists because encoding/json silently DISCARDS undeclared keys: an
// adapter DTO that forgets a field the provider sends produces no error, no
// log and no symptom — the value simply never reaches the domain. That exact
// failure emptied price/currency/listing_type/published_quantity on every
// listing (ADR-C6).
//
// It is deliberately in platform, not in one adapter: ADR-C6 applies to every
// adapter (Mercado Livre, Sankhya, xlsx), and a per-adapter copy would rebuild
// the asymmetry that caused the defect.
package rawkeys

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Undeclared returns the TOP-LEVEL keys present in raw that dto does not
// declare via a json tag, sorted, minus anything in ignore.
//
// Top-level only: nested shapes belong to their own DTO and are checked by
// passing that nested value in its own call. Reporting nested keys here would
// mix two different DTOs' responsibilities in one list.
func Undeclared(raw json.RawMessage, dto any, ignore []string) ([]string, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, fmt.Errorf("rawkeys: payload is not a JSON object: %w", err)
	}

	declared := declaredKeys(reflect.TypeOf(dto))
	skip := make(map[string]struct{}, len(ignore))
	for _, key := range ignore {
		skip[key] = struct{}{}
	}

	var missing []string
	for key := range fields {
		if _, ok := declared[key]; ok {
			continue
		}
		if _, ok := skip[key]; ok {
			continue
		}
		missing = append(missing, key)
	}
	sort.Strings(missing)
	return missing, nil
}

func declaredKeys(t reflect.Type) map[string]struct{} {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	keys := map[string]struct{}{}
	if t == nil || t.Kind() != reflect.Struct {
		return keys
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		if name == "" {
			// No json tag: encoding/json matches the Go field name
			// case-insensitively. Declaring the field name keeps the
			// detector aligned with the decoder's real behavior.
			name = field.Name
		}
		keys[name] = struct{}{}
	}
	return keys
}
