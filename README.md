### install

install.sh:
```bash
#!/usr/bin/env bash
curl -L https://github.com/romanesko/awesomeProject/archive/refs/heads/master.zip -o m.zip
unzip m.zip -d .
rm m.zip
mv awesomeProject-master awesomeserver
cd awesomeserver
touch server.env
echo 'AWP_DB_HOST=172.17.0.1' >> server.env
echo 'AWP_DB_NAME=postgres' >> server.env
echo 'AWP_DB_USER=postgres' >> server.env
echo 'AWP_DB_PASSWORD=secret' >> server.env
```

### db function example
```sql
create or replace function server_time(params json, _token uuid) returns json language plpgsql as
$$
DECLARE
    _response JSON;
    _data TEXT;
BEGIN
    _data = params->>'data';
    SELECT to_json(a) INTO _response FROM (
    
        select now() as time,
        _data as data,
        _token as token
    
    ) as a;
    RETURN coalesce(_response,'{}');
END
$$;

-- alter function server_time(json, uuid) owner to postgres;

select server_time('{"data":"test"}','83797f92-2083-4a02-a983-c48f9cd5573a')
```


### call example

```bash
curl --location 'http://localhost:9999/public/server/time?data=foo' --header 'Authorization: Token 83797f92-2083-4a02-a983-c48f9cd5573a'
```