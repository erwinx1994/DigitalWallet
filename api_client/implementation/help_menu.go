package implementation

import "fmt"

func PrintHelpMenu() {

	fmt.Println("DIGITAL WALLET INC CLIENT")
	fmt.Println()

	fmt.Println("This command allows you to deposit an amount of money in the specified wallet in the specified currency.")
	fmt.Println()

	fmt.Println("\tapi_client deposit <wallet_id> <currency> <amount>")
	fmt.Println()

	fmt.Println("This command allows you to withdraw an amount of money from the specified wallet in the specified currency.")
	fmt.Println()

	fmt.Println("\tapi_client withdraw <wallet_id> <currency> <amount>")
	fmt.Println()

	fmt.Println("This command allows you to transfer an amount of money from source to destination wallets of the same currency.")
	fmt.Println()

	fmt.Println("\tapi_client transfer <source_wallet_id> <destination_wallet_id> <currency> <amount>")
	fmt.Println()

	fmt.Println("This command allows you to get the balance of the specified wallet.")
	fmt.Println()

	fmt.Println("\tapi_client get_balance <wallet_id>")
	fmt.Println()

	fmt.Println("This command allows you to get the transaction history of the specified wallet from a start date to end date, inclusive.")
	fmt.Println()

	fmt.Println("\tapi_client get_transaction_history <wallet_id> <start_date> <end_date>")
	fmt.Println()
}
