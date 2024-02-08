## Docker run

create an `.env` file with the following content:
```
AWP_DB_HOST=localhost
AWP_DB_NAME=postgres
AWP_DB_USER=postgres
POSTGRES_PASSWORD=
```

run it in interactive mode:

```bash
  docker run -it --rm -p 9999:9999 --env-file .env --name awesome-server savarez/awesome-server
```

or run it like a daemon:

```bash
  docker run -d --rm -p 9999:9999 --env-file .env --name awesome-server savarez/awesome-server
```


## Docker-compose with dedicated Postgresql:

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