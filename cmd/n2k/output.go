package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/open-ships/n2k/pgn"
	"github.com/spf13/cobra"
)

const (
	messageOutputJSON = "json"
	messageOutputText = "text"
)

type messageWriter struct {
	out     io.Writer
	format  string
	encoder *json.Encoder
}

func addMessageOutputFlag(command *cobra.Command, outputFormat *string) {
	command.Flags().StringVar(outputFormat, "output", messageOutputJSON, "message output: json or text")
	mustConfigure(command.RegisterFlagCompletionFunc("output", cobra.FixedCompletions(
		[]string{
			"json\tTyped message JSON lines with wire values",
			"text\tConcrete PGN types with physical values and units",
		},
		cobra.ShellCompDirectiveNoFileComp,
	)))
}

func newMessageWriter(out io.Writer, format string) (*messageWriter, error) {
	writer := &messageWriter{out: out, format: format}
	switch format {
	case messageOutputJSON:
		writer.encoder = json.NewEncoder(out)
	case messageOutputText:
	default:
		return nil, fmt.Errorf("unknown message output %q: use json or text", format)
	}
	return writer, nil
}

func (writer *messageWriter) Write(message pgn.Message) error {
	switch writer.format {
	case messageOutputJSON:
		if err := writer.encoder.Encode(message); err != nil {
			return fmt.Errorf("encoding message: %w", err)
		}
	case messageOutputText:
		if _, err := fmt.Fprintln(writer.out, formatTypedMessage(message)); err != nil {
			return fmt.Errorf("writing message: %w", err)
		}
	default:
		return errors.New("message writer has no output format")
	}
	return nil
}

func formatTypedMessage(message pgn.Message) string {
	if message == nil {
		return "<nil>"
	}
	messageValue := reflect.ValueOf(message)
	if messageValue.Kind() != reflect.Pointer || messageValue.IsNil() {
		return fmt.Sprintf("%T {pgn=%d}", message, message.PGNNumber())
	}
	structValue := messageValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return fmt.Sprintf("%T {pgn=%d}", message, message.PGNNumber())
	}

	typeName := structValue.Type().Name()
	var text strings.Builder
	_, _ = fmt.Fprintf(&text, "pgn.%s {pgn=%d", typeName, message.PGNNumber())
	if carrier, ok := message.(interface{ MessageInfo() pgn.MessageInfo }); ok {
		_, _ = fmt.Fprintf(&text, " src=%d", carrier.MessageInfo().SourceId)
	}

	descriptors := messageFieldDescriptors(message.PGNNumber(), typeName)
	for index := 0; index < structValue.NumField(); index++ {
		structField := structValue.Type().Field(index)
		if structField.Name == "Info" {
			continue
		}
		fieldValue := structValue.Field(index)
		if isNilValue(fieldValue) {
			continue
		}
		name := jsonFieldName(structField)
		if name == "" || name == "-" {
			continue
		}
		descriptor := descriptors[fieldOrder(structField)]
		_, _ = fmt.Fprintf(&text, " %s=%s", name, formatTypedField(messageValue, structField, fieldValue, descriptor))
	}
	text.WriteByte('}')
	return text.String()
}

func messageFieldDescriptors(number uint32, typeName string) map[int]*pgn.FieldDescriptor {
	for _, info := range pgn.PgnInfoLookup[number] {
		if info.Id == typeName {
			return info.Fields
		}
	}
	return nil
}

func fieldOrder(field reflect.StructField) int {
	order, err := strconv.Atoi(strings.Split(field.Tag.Get("n2k"), ",")[0])
	if err != nil {
		return 0
	}
	return order
}

func jsonFieldName(field reflect.StructField) string {
	name := strings.Split(field.Tag.Get("json"), ",")[0]
	if name == "" {
		return field.Name
	}
	return name
}

func isNilValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func formatTypedField(message reflect.Value, field reflect.StructField, value reflect.Value, descriptor *pgn.FieldDescriptor) string {
	if physical, ok := physicalFieldValue(message, field.Name); ok {
		formatted := formatPhysicalValue(physical, descriptor)
		if descriptor != nil && descriptor.Unit != "" {
			return formatted + " " + descriptor.Unit
		}
		return formatted
	}

	value = indirectValue(value)
	if !value.IsValid() {
		return "null"
	}
	if descriptor != nil {
		if lookup := lookupName(descriptor); lookup != "" {
			return lookup + "(" + fmt.Sprint(value.Interface()) + ")"
		}
	}
	switch value.Kind() {
	case reflect.String:
		return strconv.Quote(value.String())
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			return "0x" + hex.EncodeToString(value.Bytes())
		}
	}
	encoded, err := json.Marshal(value.Interface())
	if err != nil {
		return fmt.Sprint(value.Interface())
	}
	return string(encoded)
}

func formatPhysicalValue(value float64, descriptor *pgn.FieldDescriptor) string {
	if descriptor == nil || descriptor.Resolution == 0 {
		return strconv.FormatFloat(value, 'f', -1, 64)
	}
	resolution := math.Abs(float64(descriptor.Resolution))
	for precision := 0; precision <= 12; precision++ {
		scaled := resolution * math.Pow10(precision)
		nearest := math.Round(scaled)
		if nearest == 0 {
			continue
		}
		if math.Abs(scaled-nearest) > math.Max(1, math.Abs(nearest))*1e-6 {
			continue
		}
		formatted := strconv.FormatFloat(value, 'f', precision, 64)
		if strings.Contains(formatted, ".") {
			formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
		}
		if formatted == "-0" {
			return "0"
		}
		return formatted
	}
	return strconv.FormatFloat(value, 'g', -1, 64)
}

func physicalFieldValue(message reflect.Value, fieldName string) (float64, bool) {
	method := message.MethodByName(fieldName + "Value")
	if !method.IsValid() || method.Type().NumIn() != 0 || method.Type().NumOut() != 2 || method.Type().Out(0).Kind() != reflect.Float64 || method.Type().Out(1).Kind() != reflect.Bool {
		return 0, false
	}
	results := method.Call(nil)
	return results[0].Float(), results[1].Bool()
}

func indirectValue(value reflect.Value) reflect.Value {
	for value.IsValid() && (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	return value
}

func lookupName(descriptor *pgn.FieldDescriptor) string {
	for _, name := range []string{
		descriptor.LookupEnumeration,
		descriptor.LookupBitEnumeration,
		descriptor.LookupIndirectEnumeration,
		descriptor.LookupFieldTypeEnumeration,
	} {
		if name != "" {
			return name
		}
	}
	return ""
}
