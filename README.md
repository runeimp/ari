Ari
===

This is a JSON log viewer written in Go. It is named in honor of [Ari Lehman][Ari Lehman - Wikipedia]. The first actor to portray [Jason Voorhees][Jason Voorhees - Wikipedia] in the [Friday the 13th][Friday the 13th (franchise) - Wikipedia] series.


Features
--------

* [ ] Read any size log file. Currently RAM limited.
* [ ] Be able to display the data of all fields
* [ ] Be able to filter out the fields you don't want displayed
* [ ] Able to handle integer fields as Unix timestamps and display them in [ISO 8601][ISO 8601 - Wikipedia] SQL style data-time format
* [ ] And last, but no lest, be fast!


Rational
--------

I've been looking for a good JSON log viewer for the command line for years. I tried [json-log-viewer][json-log-viewer - npm] early on. It made sense as JavaScript and JSON go hand-in-hand. But it was just very slow. And I didn't care for aspects of how it displayed the data. I've had similar results with all other JSON log viewers. They were often inflexible about the data they would display, seemed odd in how they displayed the data, and always too damn slow. And I'm not talking 1 GB log files with thousands of objects either. So Ari is my attempt to resolve these issues.






[Ari Lehman - Wikipedia]: https://en.wikipedia.org/wiki/Ari_Lehman
[Friday the 13th (franchise) - Wikipedia]: https://en.wikipedia.org/wiki/Friday_the_13th_(franchise)
[Jason Voorhees - Wikipedia]: https://en.wikipedia.org/wiki/Jason_Voorhees
[json-log-viewer - npm]: https://www.npmjs.com/package/json-log-viewer
[ISO 8601 - Wikipedia]: https://en.wikipedia.org/wiki/ISO_8601
