HTML Imports
============

Like many, I was saddened when the Polymer project decided to move to string
templates and away from using native HTML templates in web components.

I love the simplicity and speed of using the native browser APIs to create
custom elements. Why should HTML start its life as a JavaScript `String`, when
the browser has a native HTML parser?

The thing is, it's possible to use a template element when creating a custom
element or "web component". But without HTML imports, it would require writing
all of your component definitions _inside_ your index.html file.

This project is simply a static file server and tool that will take certain HTML
comments and tranform them into the contents of files. So you can write a web
component like this:

```html
<template id="x-foo-template">
  <h1>Hello</h1>
</template>

<script type="module">
  class XFoo extends HTMLElement {
    constructor() {
      super()
      this.template = document.getElementById('x-foo-template')
      this.attachShadow({ mode: 'open' }).appendChild(this.template.content)
    }
  }

  customElements.define('x-foo', XFoo)
</script>
```

And then include inside your `<head>` tag like this:

```html
<!-- import x-foo.html -->
```

And the html-imports server will serve up a version of your `index.html` with
that comment transformed into the content of the file. You can then use your
custom element at will.

Build with `go build` and run `html-imports`, then navigate to
`localhost:8000/example` to see the example component.

## Is this a good idea?

Probably not.

## What if I want to use Polymer?

Polymer does some magic to allow JavaScript module import syntax, which this
does not do. So it will not work.

## What about live reload?

Use Ctrl+R in the browser to activate live reload.

## Will this scale to building an entire web application?

This question is left as an exercise to the reader.
