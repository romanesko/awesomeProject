## Direct run

prepare env variables:
```bash
echo 'AWP_DB_HOST=docker.for.mac.host.internal' >> server.env
echo 'AWP_DB_NAME=dbname' >> server.env
echo 'AWP_DB_USER=user' >> server.env
echo 'AWP_DB_PASSWORD=password' >> server.env
```

run it:
```bash
docker run -d --rm -p 9999:9999 --env-file server.env --name awesome-server savarez/awesome-server
```

check logs:
```bash
docker logs awesome-server -f
```

## Docker (compose) run:

create a `.env` file with the following content:
```
PROJECT_NAME=myproject
POSTGRES_DB=myproject
POSTGRES_USER=user
POSTGRES_PASSWORD=secret
```
run it in interactive mode:

```bash
make run
```

or run it like a daemon:

```bash
make start  # use `make log` to check the logs
```

## db function example
```sql
create or replace function public.server_time(params json, _token uuid) returns json language plpgsql as
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

-- alter function public.server_time(json, uuid) owner to postgres;

select public.server_time('{"data":"test"}','83797f92-2083-4a02-a983-c48f9cd5573a')
```


## call example

```bash
curl --location 'http://localhost:9999/public/server/time?data=foo' --header 'Authorization: Token 83797f92-2083-4a02-a983-c48f9cd5573a'
```