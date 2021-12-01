Run 
```
sudo docker compose build
sudo docker compose up
```

## Notes
Only allow connection to the docker instance from 127.0.0.1 \
`iptables -I DOCKER-USER -i enp1s0 ! -s 127.0.0.1 -j DROP`

## Templating language


## TODO
* Sanitize html from the user.
* Use javascript as a templating language.
  * Look for script tags that have the template class, use these as code that generates html
  * Fill in any string literals with database data like "$userid"
  * run some function in the javascript code like generate()
  * replace the script tag with the generated html.

