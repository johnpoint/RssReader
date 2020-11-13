# deoloy.md

## Docker

*TODO*

## Download

```bash
git clone https://github.com/johnpoint/RssReader
```

## Configuration

### Backend

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

### Web

Edit `RssReader/web/src/config.js`

```javascript
export default {
  apiAddress: "http://127.0.0.1:1323"
};
```
- apiAddress: Backend url

## Compile

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

#### Nginx Url rewrite

```
location / {
  try_files $uri $uri/ /index.html;
}
```

### Databse

#### Mysql

Example:
```sql
create database rssreader character set utf8mb4;
```

## Run

```
cd RssReader/dev
./main start
```