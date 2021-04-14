const fs = require('fs');
const path = require('path');

function getRawText(fileName) {
	// File contents.
	let raw = fs.readFileSync(path.join(__dirname, fileName + '.txt'), 'utf8');

	raw = raw.replace(/(\r\n)/g, '\n');

	// Otherwise/Commonly Called: replacement.
	raw = raw.replace(/(\n.*Called:\n{2,}.*)/g, '');

	// Condensing lines downs nicely.
	raw = raw.replace(/\n{3,}/g, '\n');
	raw = raw.replace(/\n{2,}/g, '\n');

	// Replace double spaces. Idk why these are in here.
	raw = raw.replace(/ {2,}/g, ' ');

	// Titles for parsing.
	let titles = JSON.parse(
		fs.readFileSync(path.join(__dirname, fileName + '-titles.json'), 'utf8')
	);

	return [raw, titles];
}

function getRawArray(rawData) {
	return rawData.split('\n');
}

function collapseLines(rawArr, titles) {
	let books = {};
	let curBook = null;

	// Loop through our input.
	for (const line of rawArr) {
		// When we hit a title, set the curBook to this line & continue.
		for (const t of titles) {
			if (t.title.includes(line)) {
				curBook = line;
				continue;
			}
		}

		// If the object has no content, set it, otherwise append it.
		if (!books[curBook]) {
			books[curBook] = line;
		} else {
			books[curBook] += ` ${line}`;
		}
	}

	return books;
}

function parseChaptersAndVerses(rawBooks, titles) {
	let output = {
		books: [],
	};

	// Create the book in each loop.
	for (const [title, content] of Object.entries(rawBooks)) {
		let t = titles.find((t) => t.title === title);

		// Setup a starter object.
		let book = {
			...t,
			chapters: {},
		};

		// Grab all sub-verses.
		exp = /((\d+):(\d+))/g;
		let found = content.match(exp);

		for (let i = 0; i < found.length; i++) {
			let [chap, verse] = found[i].split(':');
			let left = content.indexOf(found[i]);
			let v;

			// If we have another sub-verse, get substring. Otherwise to end of string.
			if (i + 1 < found.length) {
				let right = content.indexOf(found[i + 1]);

				v = content.substring(left + found[i].length, right);
			} else {
				v = content.substring(left + found[i].length);
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

		// Push the book in.
		output.books.push(book);
	}

	// Simplify the formatting.
	for (const i in output.books) {
		const chapters = output.books[i].chapters;
		let formattedChapters = [];

		for (const key of Object.keys(chapters)) {
			const chapter = {
				chapter: key,
				verses: [],
			};

			for (const verse in chapters[key]) {
				chapter.verses.push(chapters[key][verse]);
			}

			formattedChapters.push(chapter);
		}

		output.books[i].chapters = formattedChapters;
	}

	return output;
}

function writeJSON(obj, titles, file, min = false) {
	let json = JSON.stringify({ titles, ...obj }, null, min ? 0 : 4);
	let filename = min ? `${file}-min.json` : `${file}.json`;
	let dataDir = path.join(__dirname, '../data');

	// Create directory if it doesn't exist.
	if (!fs.existsSync(dataDir)) {
		fs.mkdirSync(dataDir);
	}

	fs.writeFileSync(path.join(dataDir, `${filename}`), json);
	console.log(`${filename} written.`);
}

function parseBooks(files) {
	for (const file of files) {
		let [raw, titles] = getRawText(file);
		let rawArr = getRawArray(raw);
		let rawBooks = collapseLines(rawArr, titles);
		let output = parseChaptersAndVerses(rawBooks, titles);

		writeJSON(output, titles, file);
		writeJSON(output, titles, file, true);
	}

	console.log('Complete.');
}

const fileNames = ['kjv-old', 'kjv-new'];
parseBooks(fileNames);
