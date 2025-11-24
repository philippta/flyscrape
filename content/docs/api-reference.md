---
title: 'API Reference'
weight: 4
next: '/docs/configuration'
---

## Query API

```javascript{filename="Reference"}
// <div class="element" foo="bar">Hey</div>
const el = doc.find(".element")
el.text()                                 // "Hey"
el.html()                                 // `<div class="element">Hey</div>`
el.attr("foo")                            // "bar"
el.hasAttr("foo")                         // true
el.hasClass("element")                    // true

// <ul>
//   <li class="a">Item 1</li>
//   <li>Item 2</li>
//   <li>Item 3</li>
// </ul>
const list = doc.find("ul")
list.children()                           // [<li class="a">Item 1</li>, <li>Item 2</li>, <li>Item 3</li>]

const items = list.find("li")
items.length()                            // 3
items.first()                             // <li>Item 1</li>
items.last()                              // <li>Item 3</li>
items.get(1)                              // <li>Item 2</li>
items.get(1).prev()                       // <li>Item 1</li>
items.get(1).next()                       // <li>Item 3</li>
items.get(1).parent()                     // <ul>...</ul>
items.get(1).siblings()                   // [<li class="a">Item 1</li>, <li>Item 2</li>, <li>Item 3</li>]
items.map(item => item.text())            // ["Item 1", "Item 2", "Item 3"]
items.filter(item => item.hasClass("a"))  // [<li class="a">Item 1</li>]
```

## Document Parsing

```javascript{filename="Reference"}
import { parse } from "flyscrape";

const doc = parse(`<div class="foo">bar</div>`);
const text = doc.find(".foo").text();
```

## File Downloads

```javascript{filename="Reference"}
import { download } from "flyscrape/http";

download("http://example.com/image.jpg")              // downloads as "image.jpg"
download("http://example.com/image.jpg", "other.jpg") // downloads as "other.jpg"
download("http://example.com/image.jpg", "dir/")      // downloads as "dir/image.jpg"

// If the server offers a filename via the Content-Disposition header and no
// destination filename is provided, Flyscrape will honor the suggested filename.
// E.g. `Content-Disposition: attachment; filename="archive.zip"`
download("http://example.com/generate_archive.php", "dir/") // downloads as "dir/archive.zip"
```
