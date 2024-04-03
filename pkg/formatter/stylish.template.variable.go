package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateTemplateVariableFromItem(writer *tabwriter.Writer, item *sdk.TemplateItem) {
	for index, variable := range item.GetVariablesSchema() {
		if index == 0 {
			fmt.Fprintf(writer, "\nTemplate Variables:\n")
			fmt.Fprintf(writer, "%v\t %v\t %v\t %v\n", "Name", "Default", "Type", "Description")
		}

		tabulateTemplateVariable(writer, variable)
	}
}

func tabulateTemplateVariable(writer *tabwriter.Writer, item sdk.TemplateItemVariablesSchemaInner) {
	switch {
	case item.BooleanTypeItem != nil:
		fmt.Fprintf(
			writer,
			"%v\t %v\t %v\t %v\n",
			item.BooleanTypeItem.GetName(),
			booleanTypeToString(item.BooleanTypeItem),
			"Bool",
			item.BooleanTypeItem.GetDescription(),
		)
	case item.IntegerTypeItem != nil:
		fmt.Fprintf(
			writer,
			"%v\t %v\t %v\t %v\n",
			item.IntegerTypeItem.GetName(),
			integerTypeToString(item.IntegerTypeItem),
			"Int",
			item.IntegerTypeItem.GetDescription(),
		)
	case item.FloatTypeItem != nil:
		fmt.Fprintf(
			writer,
			"%v\t %v\t %v\t %v\n",
			item.FloatTypeItem.GetName(),
			floatTypeToString(item.FloatTypeItem),
			"Float",
			item.FloatTypeItem.GetDescription(),
		)
	case item.StringTypeItem != nil:
		fmt.Fprintf(
			writer,
			"%v\t %v\t %v\t %v\n",
			item.StringTypeItem.GetName(),
			stringTypeToString(item.StringTypeItem),
			"String",
			item.StringTypeItem.GetDescription(),
		)
	case item.EnumTypeItem != nil:
		fmt.Fprintf(
			writer,
			"%v\t %v\t %v\t %v\n",
			item.EnumTypeItem.GetName(),
			enumTypeToString(item.EnumTypeItem),
			"Enum",
			item.EnumTypeItem.GetDescription()+" (choices: "+getEnumChoices(item.EnumTypeItem)+")",
		)
	default:
		fmt.Fprintf(writer, "Unknown variable type %v\n", item)
	}
}

func getEnumChoices(item *sdk.EnumTypeItem) string {
	result := ""

	for index, item := range item.GetValues() {
		if index > 0 {
			result += ", "
		}

		result += enumChoiceToString(item)
	}

	return result
}

func enumChoiceToString(item sdk.EnumTypeItemValuesInner) string {
	switch {
	case item.BooleanValueItem != nil:
		if value, ok := item.BooleanValueItem.GetValueOk(); ok {
			return fmt.Sprintf("%t", *value)
		}
	case item.IntegerValueItem != nil:
		if value, ok := item.IntegerValueItem.GetValueOk(); ok {
			return fmt.Sprintf("%d", *value)
		}
	case item.FloatValueItem != nil:
		if value, ok := item.FloatValueItem.GetValueOk(); ok {
			return fmt.Sprintf("%f", *value)
		}
	case item.StringValueItem != nil:
		if value, ok := item.StringValueItem.GetValueOk(); ok {
			return stringOrEmpty(*value)
		}
	}

	return ""
}

func enumTypeToString(item *sdk.EnumTypeItem) string {
	defaultValue, hasDefaultValue := item.GetDefaultValueOk()
	if !hasDefaultValue {
		return ""
	}

	switch {
	case defaultValue.BooleanValueItem != nil:
		value, ok := defaultValue.BooleanValueItem.GetValueOk()
		if ok {
			return fmt.Sprintf("%t", *value)
		}
	case defaultValue.IntegerValueItem != nil:
		value, ok := defaultValue.IntegerValueItem.GetValueOk()
		if ok {
			return fmt.Sprintf("%d", *value)
		}
	case defaultValue.FloatValueItem != nil:
		value, ok := defaultValue.FloatValueItem.GetValueOk()
		if ok {
			return fmt.Sprintf("%f", *value)
		}
	case defaultValue.StringValueItem != nil:
		value, ok := defaultValue.StringValueItem.GetValueOk()
		if ok {
			return stringOrEmpty(*value)
		}
	}

	return ""
}

func booleanTypeToString(item *sdk.BooleanTypeItem) string {
	if !item.HasDefaultValue() {
		return ""
	}

	defaultValue := item.GetDefaultValue()
	if !defaultValue.HasValue() {
		return ""
	}

	return fmt.Sprintf("%t", defaultValue.GetValue())
}

func integerTypeToString(item *sdk.IntegerTypeItem) string {
	if !item.HasDefaultValue() {
		return ""
	}

	defaultValue := item.GetDefaultValue()
	if !defaultValue.HasValue() {
		return ""
	}

	return fmt.Sprintf("%d", defaultValue.GetValue())
}

func floatTypeToString(item *sdk.FloatTypeItem) string {
	if !item.HasDefaultValue() {
		return ""
	}

	defaultValue := item.GetDefaultValue()
	if !defaultValue.HasValue() {
		return ""
	}

	return fmt.Sprintf("%f", defaultValue.GetValue())
}

func stringTypeToString(item *sdk.StringTypeItem) string {
	if !item.HasDefaultValue() {
		return ""
	}

	defaultValue := item.GetDefaultValue()
	if !defaultValue.HasValue() {
		return ""
	}

	return stringOrEmpty(defaultValue.GetValue())
}

func stringOrEmpty(data string) string {
	if data == "" {
		return "''"
	}

	return data
}
