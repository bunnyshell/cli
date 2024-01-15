package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateSecretEncryptedItem(writer *tabwriter.Writer, item *sdk.SecretEncryptedItem) {
	fmt.Fprintf(writer, "%v\t %v\n", "Expression", item.GetExpression())
}

func tabulateSecretDecryptedItem(writer *tabwriter.Writer, item *sdk.SecretDecryptedItem) {
	fmt.Fprintf(writer, "%v\t %v\n", "Value", item.GetPlainText())
}
