# kjv-bible-api

I'm fairly confident the text is parsed properly to json but if it isn't then just shoot a PR in or something. Hopefully nothing is missing though.

## Swagger / OpenAPI

Open `openapi.yaml` in your favorite tool(s) such as [Insomnia](https://insomnia.rest) or the [Swagger Online Editor](https://editor.swagger.io) to see and try out all the endpoints.

If the service is down _(because I stopped hosting it)_ change the servers section to your own servers or local instance. You might need to put some extra work in for https support in your local instance.

```yaml
servers:
    - url: 'http://localhost:8080/api'
      description: 'No SSL'
    # - url: 'https://kjb.heckin.dev/api'
    #   description: 'SSL'
```

## Data Format

```
{
	"titles": [
		{
			"title": string,
			"alt": string
		}
	],
	"books": [
		{
			"title": string,
			"alt": string,
			"chapters": [
				{
					"chapter": int
					"verses": [string]
				}
			]
		}
	]
}
```

## Legal Stuff?

The raw utf-8 content for the King James Bible (kjv/kjb) was acquired from [The Project Gutenberg eBook of The King James Bible](https://www.gutenberg.org/files/10/10-0.txt). This can be seen within the files named `kjv-new.[txt/json]` and `kjv-old.[txt/json]` which have been modified for parsing.

See their license at the bottom of the [file](https://www.gutenberg.org/files/10/10-0.txt) for information about their license if you plan to use/distribute/sell any part(s) of this codebase using their text.
