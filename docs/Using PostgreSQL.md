# Using PostgreSQL

Default installation directory of PostgreSQL

    C:\Program Files\PostgreSQL\17

Add this directory to your PATH environment variable on Windows

    C:\Program Files\PostgreSQL\17\bin

Default database cluster for PostgreSQL

    C:\Program Files\PostgreSQL\17\data

Default logging directory for PostgreSQL

    C:\Program Files\PostgreSQL\17\data\log

The following files in the database cluster directory have special meanings.

    pg_hba.conf
    Client authentication configuration file. This file controls: which hosts are allowed to connect, how clients are authenticated, which PostgreSQL user names they can use, which databases they can access.

    pg_ident.conf
    Configure mapping of operating system user to PostgreSQL user names

    postgresql.conf
    Configure database server properties. IP address, port, maximum number of connections, etc.

Initialize a PostgreSQL database cluster

    pg_ctl -D /usr/local/pgsql/data initdb -o '-U bootstrap_username -W'

Start PostgreSQL database server

    pg_ctl start [-D datadir] [-l filename]

Stop PostgreSQL database server

    pg_ctl stop [-D datadir] 

Restart PostgreSQL database server

    pg_ctl restart

Connect to PostgreSQL database server

    psql -U username -d database_name 

Check if an application is already listening on a specific port

    netstat -aon | findstr 5432

## How to use PSQL

Show connection information

    \conninfo

list roles

    \dg
    \du

List schemas

    \dn

List tables in a schema

    \dt database_name.schema_name.*

Get all schema and table names

    select * from pg_tables;
