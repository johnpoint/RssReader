# deoloy.md

## Docker

*TODO*

## Download && Compile

```bash
git clone https://github.com/johnpoint/RssReader
```

### Backend

```bash
cd RssReader
bash build.sh
```

### Web

```bash
cd web
yarn
yarn build
```

### Databse

#### Mysql

Example:
```sql
create database rssreader character set utf8mb4;
```

## Configuration

Edit `RssReader/dev/config.json`

Example:

```json
{
  "Port": "1323",
  "TLS": false,
  "CERTPath": "PATHtoCER",
  "KEYPath": "PATHtoKEY",
  "Salt": "ControlCenter",
  "Database": {
    "Address": "127.0.0.1:3306",
    "User": "root",
    "Password": "DNMP-lvcshu",
    "DBname": "rssreader",
    "Type": "sqlite"
  },
  "Debug": false,
  "UpdateTime": 5,
  "AllowOrigins": [
    "http://localhost:8080",
    "http://192.168.15.107:8080"
  ]
}
```

- Port: Backend listening port
- TLS: Backend TLS 
- CERTPath: Backend TLS certificate
- KEYPath: Backend TLS private key
- Salt: String used to encrypt user password **Note: Remember to record in another safe place**
- Database
    - Address: Database address(include port)
    - User: Database username(if need)
    - Password: Database password(if need)
    - DBname: Database name
    - Type: Database type("mysql"/"sqlite")
- Debug: (true/false)
- UpdateTime: Spider sleep time (min) // TODO
- AllowOrigins: CORS setting