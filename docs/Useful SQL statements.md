# SQL statements used by the backend

SQL statements for balance service 

    Get currency and ballance for a wallet
    select currency, balance from postgres.wallet.balances where wallet_id=$1

SQL Statements for deposit service 

    Check if wallet already exists
    select currency, balance from postgres.wallet.balances where wallet_id=$1

    Add new transaction
    insert into postgres.wallet.transactions(wallet_id, date_and_time, currency, amount) values ($1, $2, $3, $4)

    Update balance
    update postgres.wallet.balances set balance=$1 where wallet_id=$2

    Insert new balance (Create account)
    insert into postgres.wallet.balances(wallet_id, currency, balance) values ($1, $2, $3)

SQL statements for withdraw service

    Check if wallet already exists and determine its currency
    select wallet_id, balance, currency from postgres.wallet.balances where wallet_id=$1

    Add new transaction
    insert into postgres.wallet.transactions(wallet_id, date_and_time, currency, amount) values ($1, $2, $3, $4)

    Update balance
    update postgres.wallet.balances set balance=$1 where wallet_id=$2

SQL statements for transaction history service

    Check if wallet already exists
    select balance from postgres.wallet.balances where wallet_id=$1

    Get transaction history
    select wallet_id, date_and_time, currency, amount from postgres.wallet.transactions where wallet_id=$1 where date_and_time>=$2 and date_and_time<=$3

SQL statements for transfer service

    Check if wallet already exists
    select currency, balance from postgres.wallet.balances where wallet_id=$1

    Add new transaction
    insert into postgres.wallet.transactions(wallet_id, date_and_time, currency, amount) values ($1, $2, $3, $4)
