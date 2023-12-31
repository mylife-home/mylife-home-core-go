package internal

import (
	"encoding/json"
	"fmt"
	"mylife-home-common/components/metadata"
	"strings"
)

type Writer struct {
	builder *strings.Builder
}

func MakeWrite(packageName string) *Writer {
	writer := &Writer{
		builder: &strings.Builder{},
	}

	writer.appendLinef(`package %s`, packageName)
	writer.appendLine(``)
	writer.appendLine(`import (`)
	writer.appendLine(`	"mylife-home-common/components/metadata"`)
	writer.appendLine(`	"mylife-home-core-library/registry"`)
	writer.appendLine(`)`)
	writer.appendLine(``)

	return writer
}

func (writer *Writer) appendBlock(str string) {
	writer.builder.WriteString(str)
}

func (writer *Writer) appendLine(str string) {
	writer.appendBlock(str + "\n")
}

func (writer *Writer) appendLinef(format string, a ...any) {
	writer.appendLine(fmt.Sprintf(format, a...))
}

func (writer *Writer) BeginPlugin(pluginType string, module string, name string, description string, usage metadata.PluginUsage, version string) {
	writer.appendLine(`func init() {`)
	writer.appendLinef(`	builder := registry.MakePluginTypeBuilder[%s](%s, %s, %s, %s, %s)`,
		pluginType,
		renderStringLiteral(module),
		renderStringLiteral(name),
		renderStringLiteral(description),
		renderPluginUsage(usage),
		renderStringLiteral(version))
}

func (writer *Writer) AddState(fieldName string, name string, description string, valueType metadata.Type) {
	writer.appendLinef(`	builder.AddState(%s, %s, %s, %s)`,
		renderStringLiteral(fieldName),
		renderStringLiteral(name),
		renderStringLiteral(description),
		renderType(valueType))
}

func (writer *Writer) AddAction(methodName string, name string, description string, valueType metadata.Type) {
	writer.appendLinef(`	builder.AddAction(%s, %s, %s, %s)`,
		renderStringLiteral(methodName),
		renderStringLiteral(name),
		renderStringLiteral(description),
		renderType(valueType))

}

func (writer *Writer) AddConfig(fieldName string, name string, description string, valueType metadata.ConfigType) {
	writer.appendLinef(`	builder.AddConfig(%s, %s, %s, %s)`,
		renderStringLiteral(fieldName),
		renderStringLiteral(name),
		renderStringLiteral(description),
		renderConfigType(valueType))
}

func (writer *Writer) EndPlugin() {
	writer.appendLine(`	registry.RegisterPlugin(builder.Build())`)
	writer.appendLine(`}`)
	writer.appendLine(``)
}

func (writer *Writer) Content() []byte {
	return []byte(writer.builder.String())
}

func renderPluginUsage(usage metadata.PluginUsage) string {
	switch usage {
	case metadata.Sensor:
		return "metadata.Sensor"
	case metadata.Actuator:
		return "metadata.Actuator"
	case metadata.Logic:
		return "metadata.Logic"
	case metadata.Ui:
		return "metadata.Ui"
	default:
		return "???"
	}
}

func renderType(typ metadata.Type) string {
	switch typed := typ.(type) {

	case *metadata.RangeType:
		return fmt.Sprintf(`metadata.MakeTypeRange(%d, %d)`, typed.Min(), typed.Max())

	case *metadata.TextType:
		return `metadata.MakeTypeText()`

	case *metadata.FloatType:
		return `metadata.MakeTypeFloat()`

	case *metadata.BoolType:
		return `metadata.MakeTypeBool()`

	case *metadata.EnumType:
		builder := strings.Builder{}
		builder.WriteString(`metadata.MakeTypeEnum(`)
		for index := 0; index < typed.NumValues(); index += 1 {
			if index > 0 {
				builder.WriteString(`, `)
			}

			builder.WriteString(`"` + typed.Value(index) + `"`)
		}
		builder.WriteString(`)`)
		return builder.String()

	case *metadata.ComplexType:
		return `metadata.MakeTypeComplex()`

	default:
		return "???"
	}
}

func renderConfigType(configType metadata.ConfigType) string {
	switch configType {
	case metadata.String:
		return "metadata.String"
	case metadata.Bool:
		return "metadata.Bool"
	case metadata.Integer:
		return "metadata.Integer"
	case metadata.Float:
		return "metadata.Float"
	default:
		return "???"
	}
}

func renderStringLiteral(value string) string {
	// annotation lexer split the strings into []byte, then encode them back using 'str += string(b)
	// this crashed UTF8 encoding.
	// We fix that doing the reverse process: converting the runes back to bytes, and reassembling the string from []byte.
	// https://github.com/YReshetko/go-annotation/issues/29

	raw := make([]byte, 0)
	for _, rune := range value {
		raw = append(raw, byte(rune))
	}
	value = string(raw)

	// Consider for now json and golang string same
	raw, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	return string(raw)
}

/*
package plugin

import (
	"mylife-home-common/components/metadata"
	"mylife-home-core-library/registry"
)

func init() {
	builder := registry.MakePluginTypeBuilder[ValueBinary]("name", "desc", metadata.Logic)
	builder.AddState("Value", "value", "desc", metadata.MakeTypeBool())
	builder.AddAction("SetValue", "setValue", "desc", metadata.MakeTypeBool())
	builder.AddConfig("InitialValue", "initialValue", "desc", metadata.Bool)
	registry.RegisterPlugin(builder.Build())
}
*/
