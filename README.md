# kjv-bible-api

This is a project for CST-235. I really didn't want to do this in Java, so I didn't.

I'm fairly confident the text is parsed properly to json but if it isn't then just shoot a PR in or something. Hopefully nothing is missing though.

## Data Parsing Preface

All parsing was done on a Mac and I know that it doesn't parse properly on a Windows system. There seems to also be an issue with it finding the folder called `data` which you can create yourself before running `book-parser.js` in the root of the project. 

With that said, fixing this is not a priority as of right now & I'd recommend you parse the text through some unix system which should produce expected results until I get around to fixing it.

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
