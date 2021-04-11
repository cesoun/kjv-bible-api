const fs = require('fs');
const path = require('path');

function getRawText(fileName) {
	return fs.readFileSync(path.join(__dirname, fileName + '.txt'), 'utf8');
}

function getRawArray(rawData) {
	return rawData.split('\n\n\n');
}

function getRawBooks(rawArr) {
	let books = {};
	for (let i = 0; i < rawArr.length / 2; i += 2) {
		books[rawArr[i].replace('\n\n', '')] = rawArr[i + 1];
	}

	return books;
}

function parseChaptersAndVerses(rawBooks) {
	let output = {
		books: [],
	};

	// Create the book in each loop.
	for (const [k, v] of Object.entries(rawBooks)) {
		// Setup a starter object.
		let book = {
			title: k,
			// altTitle: '',
			chapters: {},
		};

		// Get the chapter contents & split it out.
		let contents = v.split('\n\n');

		// Loop through it.
		for (let index in contents) {
			// hold the content
			let c = contents[index];

			// Repalce all newlines.
			let exp = /\n/g;
			c = c.replace(exp, ' ');

			// Grab all potential sub-verses.
			exp = /((\d+):(\d+))/g;
			let found = c.match(exp);

			// Hanging lines, these are likely apart of the previous verse.
			if (!found) {
				// Append the content to the previous line.
				let prev = contents[index - 1];
				prev += ` ${c}`;

				contents[index - 1] = prev;

				// Remove it and continue.
				delete contents[index];
				continue;
			}

			// non-hanging, we need to insert them into the chapter:verse
			if (found.length > 1) {
				// move index ahead 1 & insert after current node.
				for (let i = 0; i < found.length; i++) {
					let [chap, verse] = found[i].split(':');
					let left = c.indexOf(found[i]);
					let v;

					// If we have another sub-verse, get substring. Otherwise to end of string.
					if (i + 1 < found.length) {
						let right = c.indexOf(found[i + 1]);

						v = c.substring(left + found[i].length, right);
					} else {
						v = c.substring(left + found[i].length);
					}

					// Append the verse to chapter.
					chapter = book.chapters[chap];

					if (!chapter) {
						chapter = {};
					}

					// Trim whitespace.
					chapter[verse] = v.trim();
					book.chapters[chap] = chapter;
				}
			} else {
				// Append the verse to the chapter.
				let [chap, verse] = found[0].split(':');

				chapter = book.chapters[chap];

				if (!chapter) {
					chapter = {};
				}

				// Remove the ##:## (chapter:verse) from the string and trim whitespace.
				chapter[verse] = c.replace(exp, '').trim();
				book.chapters[chap] = chapter;
			}
		}

		// Push the book in.
		output.books.push(book);
	}

	return output;
}

function writeJSON(obj, file, min = false) {
	let json = JSON.stringify(obj, null, min ? 0 : 4);
	let filename = min ? `${file}-min.json` : `${file}.json`;

	fs.writeFileSync(path.join(__dirname, '../data', `${filename}`), json);
	console.log(`${filename} written.`);
}

function parseBooks(files) {
	for (const file of files) {
		let raw = getRawText(file);
		let rawArr = getRawArray(raw);
		let rawBooks = getRawBooks(rawArr);
		let output = parseChaptersAndVerses(rawBooks);

		// Add ", true" to get a minified file.
		writeJSON(output, file);
	}

	console.log('Complete.');
}

const fileNames = ['kjv-old', 'kjv-new'];
parseBooks(fileNames);
