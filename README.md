Run 
```
sudo docker compose build
sudo docker compose up
```

## Notes
Only allow connection to the docker instance from 127.0.0.1 \
`iptables -I DOCKER-USER -i enp1s0 ! -s 127.0.0.1 -j DROP`

## Documentation W.I.P
<attiribute> <query>

## TODO
* Sanitize html from the user.

