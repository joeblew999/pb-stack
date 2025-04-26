# locale

Covers:

- Detection os the language of the User.

- Message Translation for many languages.

- Units translation such as Date, DataTime, Currency, Address, Scientific Units ( such as Temperature ).

## Detection

https://github.com/jeandeaual/go-locale does a great jop of this with no dependencies.

## Translation

We are using the ICU, which is an open Unicode standard. 

ICU, or International Components for Unicode, is a set of open-source libraries used for Unicode support and software internationalization

There are so many formats that its gotten ridiculous and ICU works well but had not been well adapter to golang yet.

https://github.com/romshark/toki
https://github.com/romshark/icumsg

## Compile versus Runtime

The module support CompileTime and RunTime.

CompileTime is when you want your Golang code tightly coupled to the ICU system.

RunTime is when you want your JS to be loosely couple to the ICU system.

## Admin

A Web GUI using HTMX so that AI and Humans can help translate the files.

The REST API using Ogen and SSE provides this for both Humans and AI.

As messages are fixed, they can be used immediately or pushed via git.

Host on Cloudflare or Local or any Cloud.

For Cloudflare using https://github.com/syumai/workers




